package soundboard

import (
	"errors"
	"net/http"
	"io"
	"os"
	"os/exec"
	"time"
	log "github.com/Sirupsen/logrus"
	"github.com/t11230/ramenbot/lib/modules/modulebase"
	"github.com/t11230/ramenbot/lib/sound"
	"github.com/t11230/ramenbot/lib/utils"
	"github.com/t11230/ramenbot/lib/bits"
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"image/png"
	"strconv"
	"sync"
	"github.com/layeh/gopus"
)

// Module name used in the config file
const (
	ConfigName = "soundboard"
	helpString = "**!!s** : This module allows the user to play sounds from a dank soundboard.\n"
	sHelpStringHead = `**S**
This module allows the user to play sounds from a dank soundboard.

**usage:** !!s *collection* *sound*
	Plays the sound *sound* from the collection of sounds *collection*
	For a list of sounds in a collection, type !!s *collection* help
	Collections listed below:
`
	sHelpStringTail = `
**EXAMPLE:** !!s airhorn default

For the command to upload sounds to the soundboard, type **!!s upload help**

For the command to disable sounds for a period of time, type **!!s silence help**
`

	uploadHelpString = `**UPLOAD**
This module allows the user to upload sounds to the bot's soundboard.

**usage:** put !!s upload *collection* *soundname* in the comments of an audio file attachment
Processes the attached soundfile and adds it to the soundboard as *soundname* in the collection *collection*
**WARNING** Uploading a sound to the soundboard costs **300 bits**
`

	silenceHelpString = `**SILENCE**
This module allows the user to upload sounds to the bot's soundboard.

**usage:** !!s silence *duration*
Prevents any sound clips from being played for *duration* minutes
**WARNING** Silencing the soundboard costs **100 bits** per minute
`

	// The current version of the DCA format
	FormatVersion int8 = 1

	// The current version of the DCA program
	ProgramVersion string = "0.0.1"

	GitHubRepositoryURL string = "https://github.com/bwmarrin/dca"

	Volume = 256
	Channels = 2
	FrameRate = 48000
	FrameSize = 960
	Bitrate = 64
	CoverFormat = "jpeg"
	Application = "audio"
)

var (
	SoundCommandsEnabled = true
	// Buffer for some commands
	CmdBuf bytes.Buffer
	PngBuf bytes.Buffer

	CoverImage string

	// Metadata structures
	Metadata    MetadataStruct
	FFprobeData FFprobeMetadata

	// Magic bytes to write at the start of a DCA file
	MagicBytes string = fmt.Sprintf("DCA%d", FormatVersion)

	MaxBytes  int // max size of opus data


	OpusEncoder *gopus.Encoder

	InFile  string
	InFD    *os.File

	OutFile string
	OutFD    *os.File
	OutBuf  []byte

	EncodeChan chan []int16
	OutputChan chan []byte

	err error

	wg sync.WaitGroup
)


// List of commands that this module accepts
var commandTree = []modulebase.ModuleCommandTree{
	{
		RootCommand: "s",
		SubKeys:     modulebase.SK{
			"upload": modulebase.CN{
				Function: uploadSoundFile,
			},
			"silence": modulebase.CN{
				Function: silenceSounboard,
			},
		},
		Function:    handleSoundCommand,
	},
}

// Called to initialize this module
func SetupFunc(config *modulebase.ModuleConfig) (*modulebase.ModuleSetupInfo, error) {
	return &modulebase.ModuleSetupInfo{
		Events:   nil,
		Commands: &commandTree,
		Help:     helpString,
	}, nil
}

func uploadSoundFile(cmd *modulebase.ModuleCommand) (string, error) {
	if len(cmd.Args)!= 3 {
		return uploadHelpString, nil
	}
	user := cmd.Message.Author
	if bits.GetBits(cmd.Guild.ID, user.ID) < 300 {
		return "**FAILED TO ADD SOUND:** Insufficient bits.", nil
	}
	prefix := cmd.Args[0]
	name := cmd.Args[1]
	filename := prefix+"_"+name
	InFile = "audio/"+filename
	OutFile = "audio/"+filename+".dca"
	InFD, _ = os.Create(InFile)
	OutFD, _ = os.Create(OutFile)
	defer InFD.Close()
	defer InFD.Close()
	defer os.Remove(InFile)
	resp, _ := http.Get(cmd.Args[2])
	defer resp.Body.Close()
	io.Copy(InFD, resp.Body)
	result := processUploadFile()
	if result < 0 {
		if result == -2 {
			return "**FAILED TO ADD SOUND:** Clip was longer than 15 seconds.", nil
		}
		return "**FAILED TO ADD SOUND:** FFProbe Error", nil
	}
	log.Debug("Starting to add sound to soundboard")
	collections := sound.GetCollections()
	log.Debug("Got Collections")
	for _, collection := range collections {
		if collection.Prefix == prefix {
			log.Debug("Existing Collection")
			for _, sound := range collection.Sounds {
				if sound.Name == name{
					return "**FAILED TO ADD SOUND:** Sound with that name already exists in collection", nil
				}
			}
			newSound := sound.CreateSound(name, 50, 0)
			collection.Sounds = append(collection.Sounds, newSound)
			newSound.Load(collection)
			log.Debug("Added sound")
			bits.RemoveBits(cmd.Session, cmd.Guild.ID, user.ID, 300, "Added sound "+prefix+"_"+name)
			return "**"+utils.GetPreferredName(cmd.Guild, user.ID)+"** added sound **"+name+"** to collection **"+prefix+"**", nil
		}
	}
	log.Debug("New Collection")
	var NEW *sound.SoundCollection = &sound.SoundCollection{
		Prefix: prefix,
		Commands: []string{
			prefix,
		},
		Sounds: []*sound.Sound{
			sound.CreateSound(name, 50, 0),
		},
	}
	sound.AddCollection(NEW)
	NEW.Sounds[0].Load(NEW)
	log.Debug("Added sound")
	bits.RemoveBits(cmd.Session, cmd.Guild.ID, user.ID, 300, "Added sound "+prefix+"_"+name)
	return "**"+utils.GetPreferredName(cmd.Guild, user.ID)+"** added sound **"+name+"** to collection **"+prefix+"**", nil
}
//ALL THIS CODE FROM https://github.com/bwmarrin/dca/cmd/dca/main.go

func processUploadFile() int {
	MaxBytes = (FrameSize * Channels) * 2
	OpusEncoder, err = gopus.NewEncoder(FrameRate, Channels, gopus.Audio)
	if err != nil {
		fmt.Println("NewEncoder Error:", err)
		return -1
	}

	OpusEncoder.SetBitrate(Bitrate * 1000)

	OpusEncoder.SetApplication(gopus.Audio)

	OutputChan = make(chan []byte, 10)
	EncodeChan = make(chan []int16, 10)

	Metadata = MetadataStruct{
		Dca: &DCAMetadata{
			Version: FormatVersion,
			Tool: &DCAToolMetadata{
				Name:    "dca",
				Version: ProgramVersion,
				Url:     GitHubRepositoryURL,
				Author:  "bwmarrin",
			},
		},
		SongInfo: &SongMetadata{},
		Origin:   &OriginMetadata{},
		Opus: &OpusMetadata{
			Bitrate:     Bitrate * 1000,
			SampleRate:  FrameRate,
			Application: Application,
			FrameSize:   FrameSize,
			Channels:    Channels,
		},
		Extra: &ExtraMetadata{},
	}

	_ = Metadata

	// get ffprobe data
	ffprobe := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", InFile)
	ffprobe.Stdout = &CmdBuf


	err = ffprobe.Start()
	if err != nil {
		fmt.Println("RunStart Error:", err)
		return -1
	}

	err = ffprobe.Wait()
	if err != nil {
		fmt.Println("FFprobe Error:", err)
		return -1
	}

	err = json.Unmarshal(CmdBuf.Bytes(), &FFprobeData)
	if err != nil {
		fmt.Println("Error unmarshaling the FFprobe JSON:", err)
		return -1
	}

	bitrateInt, err := strconv.Atoi(FFprobeData.Format.Bitrate)
	if err != nil {
		fmt.Println("Could not convert bitrate to int:", err)
		return -1
	}
	duration, _ := strconv.ParseFloat(FFprobeData.Format.Duration, 32)
	if duration > 15.0 {
		return -2
	}
	log.Debug("Finished ffprobe")
	log.Debugf("%v", FFprobeData)
	Metadata.SongInfo = &SongMetadata{
		Title:    FFprobeData.Format.Tags.Title,
		Artist:   FFprobeData.Format.Tags.Artist,
		Album:    FFprobeData.Format.Tags.Album,
		Genre:    FFprobeData.Format.Tags.Genre,
		Comments: "", // change later?
	}

	Metadata.Origin = &OriginMetadata{
		Source:   "file",
		Bitrate:  bitrateInt,
		Channels: Channels,
		Encoding: FFprobeData.Format.FormatLongName,
	}

	CmdBuf.Reset()

	// get cover art
	cover := exec.Command("ffmpeg", "-loglevel", "0", "-i", InFile, "-f", "singlejpeg", "pipe:1")
	cover.Stdout = &CmdBuf

	err = cover.Start()
	if err != nil {
		fmt.Println("RunStart Error:", err)
		return -1
	}

	err = cover.Wait()
	if err == nil {
		buf := bytes.NewBufferString(CmdBuf.String())

		if CoverFormat == "png" {
			img, err := jpeg.Decode(buf)
			if err == nil { // silently drop it, no image
				err = png.Encode(&PngBuf, img)
				if err == nil {
					CoverImage = base64.StdEncoding.EncodeToString(PngBuf.Bytes())
				}
			}
		} else {
			CoverImage = base64.StdEncoding.EncodeToString(CmdBuf.Bytes())
		}

		Metadata.SongInfo.Cover = &CoverImage
	}

	CmdBuf.Reset()
	PngBuf.Reset()

	//////////////////////////////////////////////////////////////////////////
	// BLOCK : Start reader and writer workers
	//////////////////////////////////////////////////////////////////////////

	wg.Add(1)
	go reader()

	wg.Add(1)
	go encoder()

	wg.Add(1)
	go writer()

	// wait for above goroutines to finish, then exit.
	wg.Wait()
	log.Debug("Finished processing file")
	return 0
}

// reader reads from the input
func reader() {
	log.Debug("Started reading")
	defer func() {
		close(EncodeChan)
		wg.Done()
	}()

	// read from file

	// Create a shell command "object" to run.
	ffmpeg := exec.Command("ffmpeg", "-i", InFile, "-vol", strconv.Itoa(Volume), "-f", "s16le", "-ar", strconv.Itoa(FrameRate), "-ac", strconv.Itoa(Channels), "pipe:1")
	stdout, err := ffmpeg.StdoutPipe()
	if err != nil {
		fmt.Println("StdoutPipe Error:", err)
		return
	}

	// Starts the ffmpeg command
	err = ffmpeg.Start()
	if err != nil {
		fmt.Println("RunStart Error:", err)
		return
	}

	for {

		// read data from ffmpeg stdout
		InBuf := make([]int16, FrameSize*Channels)
		err = binary.Read(stdout, binary.LittleEndian, &InBuf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			log.Debug("Finished reading")
			return
		}
		if err != nil {
			fmt.Println("error reading from ffmpeg stdout :", err)
			return
		}

		// write pcm data to the EncodeChan
		EncodeChan <- InBuf
	}
}

// encoder listens on the EncodeChan and encodes provided PCM16 data
// to opus, then sends the encoded data to the OutputChan
func encoder() {
	log.Debug("Started encoding")
	defer func() {
		close(OutputChan)
		wg.Done()
	}()

	for {
		pcm, ok := <-EncodeChan
		if !ok {
			// if chan closed, exit
			log.Debug("Finished encoding")
			return
		}

		// try encoding pcm frame with Opus
		opus, err := OpusEncoder.Encode(pcm, FrameSize, MaxBytes)
		if err != nil {
			log.Error("Encoding Error:", err)
			return
		}

		// write opus data to OutputChan
		OutputChan <- opus
	}
}

// writer listens on the OutputChan and writes the output to stdout pipe
// TODO: Add support for writing directly to a file
func writer() {
	log.Debug("Started writing")
	defer wg.Done()

	var opuslen int16
	// var jsonlen int32

	// 16KB output buffer
	wbuf := bufio.NewWriterSize(OutFD, 16384)
	defer wbuf.Flush()
	for {
		opus, ok := <-OutputChan
		if !ok {
			// if chan closed, exit
			log.Debug("Finished writing")
			return
		}

		// write header
		opuslen = int16(len(opus))
		err = binary.Write(wbuf, binary.LittleEndian, &opuslen)
		if err != nil {
			fmt.Println("error writing output: ", err)
			return
		}

		// write opus data to stdout
		err = binary.Write(wbuf, binary.LittleEndian, &opus)
		if err != nil {
			fmt.Println("error writing output: ", err)
			return
		}
	}
}

//END OF CODE FROM https://github.com/bwmarrin/dca/cmd/dca/main.go

func handleSoundCommand(cmd *modulebase.ModuleCommand) (string, error) {
	log.Debugf("Sound :%v", cmd.Args)
	if len(cmd.Args) == 0 || cmd.Args[0] =="help" {
		availableCollections()
		return sHelpStringHead+sound.PrintCollections()+sHelpStringTail, nil
	}
	if !SoundCommandsEnabled {
		return "**SOUND COMMANDS ARE CURRENTLY DISABLED**", nil
	}
	for _, coll := range sound.GetCollections() {
		if utils.Scontains(cmd.Args[0], coll.Commands...) {

			// If they passed a specific sound effect, find and select that (otherwise play nothing)
			var snd *sound.Sound
			if len(cmd.Args) > 1 {
				if cmd.Args[1]=="help" {
					return sound.PrintCollection(coll), nil
				}
				for _, s := range coll.Sounds {
					if cmd.Args[1] == s.Name {
						snd = s
					}
				}

				if snd == nil {
					return "", errors.New("Sound was nil")
				}
			} else {
				return sound.PrintCollection(coll), nil
			}

			go sound.EnqueuePlay(cmd.Session, cmd.Message.Author, cmd.Guild, coll, snd)
			return "", nil
		}
	}

	return "Unable to find sound", nil
}

func silenceSounboard(cmd *modulebase.ModuleCommand) (string, error) {
	if (len(cmd.Args)!=1)||(cmd.Args[0]=="help"){
		return silenceHelpString, nil
	}
	duration, err := strconv.Atoi(cmd.Args[0])
	if err != nil {
		return "**Invalid silence duration**", nil
	}
	user := cmd.Message.Author
	if bits.GetBits(cmd.Guild.ID, user.ID) < (duration * 100) {
		return "**FAILED TO SILENCE SOUNDBOARD:** Insufficient bits.", nil
	}
	SoundCommandsEnabled = false
	time.Sleep(time.Duration(duration)*time.Minute)
	SoundCommandsEnabled = true
	return "**Soundboard is no longer silenced**", nil
}

func availableCollections() []string {
	colls := []string{}
	for _, c := range sound.GetCollections() {
		colls = append(colls, c.Commands[0])
	}
	return colls
}
