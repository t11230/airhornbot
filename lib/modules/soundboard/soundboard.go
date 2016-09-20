package soundboard

import (
	"errors"
	"net/http"
	"io"
	"os"
	"os/exec"
	log "github.com/Sirupsen/logrus"
	"github.com/t11230/ramenbot/lib/modules/modulebase"
	"github.com/t11230/ramenbot/lib/sound"
	"github.com/t11230/ramenbot/lib/utils"
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
	sHelpString = `**S**
This module allows the user to play sounds from a dank soundboard.

**usage:** !!s *collection* *sound*
	Plays the sound *sound* from the collection of sounds *collection*
	Collections and the sounds they contain listed below:

**airhorn:** default, reverb, spam, tripletap, fourtap, distant, echo, clownfull, clownshort, clownspam, highfartlong, highfartshort, midshort, truck
**anotha:** one, one_echo, one_classic, dialup
**johncena:** airhorn, echo, full, jc, nameis, spam, collect
**ethan:** areyou_classic, areyou_condensed, areyou_crazy, areyou_ethan, classic, echo, high, slowandlow, cuts, beat, sodiepop, vape
**stan:** herd, moo, x3
**trump:** 10ft, wall, mess, bing, getitout, tractor, worstpres, china, mexico, special
**music:** serbian, techno
**meme:** headshot, wombo, triple, camera, gandalf, mad, ateam, bennyhill, tuba, donethis, leeroy, slam, nerd, kappa, digitalsports, csi, nogod, welcomebdc
**birthday:** horn, horn3, sadhorn, weakhorn
**owult:** dva_enemy, genji_enemy, genji_friendly, hanzo_enemy, hanzo_friendly, junkrat_enemy, junkrat_friendly, lucio_friendly, lucio_enemy, mccree_enemy, mccree_friendly, mei_friendly, mei_enemy, pharah_enemy, reaper_friendly, 76_enemy, symmetra_friendly, torbjorn_enemy, tracer_enemy, tracer_friendly, widow_enemy, widow_friendly, zarya_enemy, zarya_friendly, zenyatta_enemy, dva_;), anyong
**dota:** waow, balance, rekt, stick, mana, disaster, liquid, history, smut, team, aegis
**overwatch:** payload, whoa, woah, winky, turd, ryuugawagatekiwokurau, cyka, noon, somewhere, lift, russia
**wc3:** work, awake
**sp:** screw, authority
**sv:** piss, fucks, shittalk, attractive, win
**archer:** dangerzone, klog

**EXAMPLE:** !!s airhorn default

For the command to upload sounds to the soundboard, type **!!s upload help**
`

	uploadHelpString = `**UPLOAD**
This module allows the user to upload sounds to the bot's soundboard.

**usage:** put !!s upload *collection* *soundname* in the comments of an audio file attachment
Processes the attached soundfile and adds it to the soundboard as *soundname* in the collection *collection*
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

// Base metadata struct
//
// https://github.com/bwmarrin/dca/issues/5#issuecomment-189713886
type MetadataStruct struct {
    Dca             *DCAMetadata    `json:"dca"`
    SongInfo        *SongMetadata   `json:"info"`
    Origin          *OriginMetadata `json:"origin"`
    Opus            *OpusMetadata   `json:"opus"`
    Extra           *ExtraMetadata  `json:"extra"`
}

// DCA metadata struct
//
// Contains the DCA version.
type DCAMetadata struct {
    Version int8                `json:"version"`
    Tool    *DCAToolMetadata    `json:"tool"`
}

// DCA tool metadata struct
//
// Contains the Git revisions, commit author etc.
type DCAToolMetadata struct {
    Name        string  `json:"name"`
    Version     string  `json:"version"`
    Url         string  `json:"url"`
    Author      string  `json:"author"`
}

// Song Information metadata struct
//
// Contains information about the song that was encoded.
type SongMetadata struct {
    Title       string  `json:"title"`
    Artist      string  `json:"artist"`
    Album       string  `json:"album"`
    Genre       string  `json:"genre"`
    Comments    string  `json:"comments"`
    Cover       *string `json:"cover"`
}

// Origin information metadata struct
//
// Contains information about where the song came from,
// audio bitrate, channels and original encoding.
type OriginMetadata struct {
    Source      string  `json:"source"`
    Bitrate     int     `json:"abr"`
    Channels    int     `json:"channels"`
    Encoding    string  `json:"encoding"`
    Url         string  `json:"url"`
}

// Opus metadata struct
//
// Contains information about how the file was encoded
// with Opus.
type OpusMetadata struct {
    Bitrate     int     `json:"abr"`
    SampleRate  int     `json:"sample_rate"`
    Application string  `json:"mode"`
    FrameSize   int     `json:"frame_size"`
    Channels    int     `json:"channels"`
}

// Extra metadata struct
type ExtraMetadata struct {}

////////////////////////////////////////////////////////
/// FFprobe Structures
////////////////////////////////////////////////////////

type FFprobeMetadata struct {
    Format  *FFprobeFormat  `json:"format"`
}

type FFprobeFormat struct {
    FileName        string          `json:"filename"`
    NumStreams      int             `json:"nb_streams"`
    NumPrograms     int             `json:"nb_programs"`
    FormatName      string          `json:"format_name"`
    FormatLongName  string          `json:"format_long_name"`
    StartTime       string          `json:"start_time"`
    Duration        string          `json:"duration"`
    Size            string          `json:"size"`
    Bitrate         string          `json:"bit_rate"`
    ProbeScore      int             `json:"probe_score"`

    Tags            *FFprobeTags    `json:"tags"`
}

type FFprobeTags struct {
    Date        string  `json:"date"`
    Track       string  `json:"track"`
    Artist      string  `json:"artist"`
    Genre       string  `json:"genre"`
    Title       string  `json:"title"`
    Album       string  `json:"album"`
    Compilation string  `json:"compilation"`
}

// List of commands that this module accepts
var commandTree = []modulebase.ModuleCommandTree{
	{
		RootCommand: "s",
		SubKeys:     modulebase.SK{
			"upload": modulebase.CN{
				Function: uploadSoundFile,
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
	prefix := cmd.Args[0]
	name := cmd.Args[1]
	filename := prefix+"_"+name
	InFile = "audio/"+filename
	OutFile = "audio/"+filename+".dca"
	InFD, _ = os.Create(InFile)
	OutFD, _ = os.Create(OutFile)
	defer InFD.Close()
	defer InFD.Close()
	resp, _ := http.Get(cmd.Args[2])
	defer resp.Body.Close()
	io.Copy(InFD, resp.Body)
	processUploadFile()
	log.Debug("Starting to add sound to soundboard")
	collections := sound.GetCollections()
	log.Debug("Got Collections")
	for _, collection := range collections {
		if collection.Prefix == prefix {
			log.Debug("Existing Collection")
			newSound := sound.CreateSound(name, 50, 0)
			collection.Sounds = append(collection.Sounds, newSound)
			newSound.Load(collection)
			log.Debug("Added sound")
		}
	}
	for _, collection := range collections {
		for _, sound := range collection.Sounds {
			log.Debug("%v", sound.Name)
		}
	}
	// log.Debug("New Collection")
	// var NEW *sound.SoundCollection = &sound.SoundCollection{
	// 	Prefix: prefix,
	// 	Commands: []string{
	// 		prefix,
	// 	},
	// 	Sounds: []*sound.Sound{
	// 		sound.CreateSound(name, 50, 0),
	// 	},
	// }
	// sound.AddCollection(NEW)
	// NEW.Sounds[0].Load(NEW)
	// log.Debug("Added sound")
	return "", nil
}

func processUploadFile() {
	MaxBytes = (FrameSize * Channels) * 2
	OpusEncoder, err = gopus.NewEncoder(FrameRate, Channels, gopus.Audio)
	if err != nil {
		fmt.Println("NewEncoder Error:", err)
		return
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
		return
	}

	err = ffprobe.Wait()
	if err != nil {
		fmt.Println("FFprobe Error:", err)
		return
	}

	err = json.Unmarshal(CmdBuf.Bytes(), &FFprobeData)
	if err != nil {
		fmt.Println("Error unmarshaling the FFprobe JSON:", err)
		return
	}

	bitrateInt, err := strconv.Atoi(FFprobeData.Format.Bitrate)
	if err != nil {
		fmt.Println("Could not convert bitrate to int:", err)
		return
	}

	log.Debug("Finished ffprobe")

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
		return
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
	return
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

	// write the magic bytes
	// wbuf.Write([]byte(MagicBytes))
	//
	// // encode and write json length
	// json, err := json.Marshal(Metadata)
	// if err != nil {
	// 	fmt.Println("Failed to encode the Metadata JSON:", err)
	// 	return
	// }
	//
	// jsonlen = int32(len(json))
	// err = binary.Write(wbuf, binary.LittleEndian, &jsonlen)
	// if err != nil {
	// 	fmt.Println("error writing output: ", err)
	// 	return
	// }
	//
	// // write the actual json
	// wbuf.Write(json)

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

func handleSoundCommand(cmd *modulebase.ModuleCommand) (string, error) {
	log.Debugf("Sound :%v", cmd.Args)
	if len(cmd.Args) == 0 || cmd.Args[0] =="help" {
		availableCollections()
		return sHelpString, nil
	}

	for _, coll := range sound.GetCollections() {
		if utils.Scontains(cmd.Args[0], coll.Commands...) {

			// If they passed a specific sound effect, find and select that (otherwise play nothing)
			var snd *sound.Sound
			if len(cmd.Args) > 1 {
				for _, s := range coll.Sounds {
					if cmd.Args[1] == s.Name {
						snd = s
					}
				}

				if snd == nil {
					return "", errors.New("Sound was nil")
				}
			}

			go sound.EnqueuePlay(cmd.Session, cmd.Message.Author, cmd.Guild, coll, snd)
			return "", nil
		}
	}

	return "Unable to find sound", nil
}

func availableCollections() []string {
	colls := []string{}
	for _, c := range sound.GetCollections() {
		colls = append(colls, c.Commands[0])
	}
	return colls
}
