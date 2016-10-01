package sound

import (
	// "bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"
	"strings"
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"

	"github.com/t11230/ramenbot/lib/utils"
)

const (
	// Sound encoding settings
	BITRATE        = 128
	MAX_QUEUE_SIZE = 6
)

var (
	// Map of Guild id's to *Play channels, used for queuing and rate-limiting guilds
	queues map[string]chan *Play = make(map[string]chan *Play)
	collections []*SoundCollection = []*SoundCollection{}
)

// Play represents an individual use of the !airhorn command
type Play struct {
	GuildID   string
	ChannelID string
	UserID    string
	Sound     *Sound

	// The next play to occur after this, only used for chaining sounds like anotha
	Next *Play

	// If true, this was a forced play using a specific airhorn sound name
	Forced bool
}

type SoundCollection struct {
	Prefix    string
	Commands  []string
	Sounds    []*Sound
	ChainWith *SoundCollection

	soundRange int
}

// Sound represents a sound clip
type Sound struct {
	Name string

	// Weight adjust how likely it is this song will play, higher = more likely
	Weight int

	// Delay (in milliseconds) for the bot to wait before sending the disconnect request
	PartDelay int

	// Buffer to store encoded PCM packets
	buffer [][]byte
}

// Create a Sound struct
func CreateSound(Name string, Weight int, PartDelay int) *Sound {
	return &Sound{
		Name:      Name,
		Weight:    Weight,
		PartDelay: PartDelay,
		buffer:    make([][]byte, 0),
	}
}

func (sc *SoundCollection) Load() {
	for _, sound := range sc.Sounds {
		sc.soundRange += sound.Weight
		sound.Load(sc)
	}
}

func (s *SoundCollection) Random() *Sound {
	var (
		i      int
		number int = utils.RandomRange(0, s.soundRange)
	)

	for _, sound := range s.Sounds {
		i += sound.Weight

		if number < i {
			return sound
		}
	}
	return nil
}

// Load attempts to load an encoded sound file from disk
// DCA files are pre-computed sound files that are easy to send to Discord.
// If you would like to create your own DCA files, please use:
// https://github.com/nstafie/dca-rs
// eg: dca-rs --raw -i <input wav file> > <output file>
func (s *Sound) Load(c *SoundCollection) error {
	path := fmt.Sprintf("audio/%v_%v.dca", c.Prefix, s.Name)

	file, err := os.Open(path)

	if err != nil {
		fmt.Println("error opening dca file :", err)
		return err
	}

	var opuslen int16

	for {
		// read opus frame length from dca file
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil
		}

		if err != nil {
			fmt.Println("error reading frame length from dca file :", err)
			return err
		}

		// read encoded pcm from dca file
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)
		// Should not be any end of file errors
		if err != nil {
			fmt.Println("error reading encoded pcm from dca file :", err)
			return err
		}

		// append encoded pcm data to the buffer
		s.buffer = append(s.buffer, InBuf)
	}
}

// Plays this sound over the specified VoiceConnection
func (s *Sound) Play(vc *discordgo.VoiceConnection) {
	vc.Speaking(true)
	defer vc.Speaking(false)

	for _, buff := range s.buffer {
		vc.OpusSend <- buff
	}
}

// Prepares a play
func CreatePlay(s *discordgo.Session, user *discordgo.User, guild *discordgo.Guild, coll *SoundCollection, sound *Sound) *Play {
	// Grab the users voice channel
	channel := utils.GetCurrentVoiceChannel(s, user, guild)
	if channel == nil {
		log.WithFields(log.Fields{
			"user":  user.ID,
			"guild": guild.ID,
		}).Warning("Failed to find channel to play sound in")
		return nil
	}

	// Create the play
	play := &Play{
		GuildID:   guild.ID,
		ChannelID: channel.ID,
		UserID:    user.ID,
		Sound:     sound,
		Forced:    true,
	}

	if coll == nil {
		if sound == nil {
			return nil
		}
		return play
	}

	// If we didn't get passed a manual sound, generate a random one
	if play.Sound == nil {
		play.Sound = coll.Random()
		play.Forced = false
	}

	// If the collection is a chained one, set the next sound
	if coll.ChainWith != nil {
		play.Next = &Play{
			GuildID:   play.GuildID,
			ChannelID: play.ChannelID,
			UserID:    play.UserID,
			Sound:     coll.ChainWith.Random(),
			Forced:    play.Forced,
		}
	}

	return play
}

// Prepares and enqueues a play into the ratelimit/buffer guild queue
func EnqueuePlay(s *discordgo.Session, user *discordgo.User, guild *discordgo.Guild, coll *SoundCollection, sound *Sound) {
	play := CreatePlay(s, user, guild, coll, sound)
	if play == nil {
		return
	}

	// Check if we already have a connection to this guild
	//   yes, this isn't threadsafe, but its "OK" 99% of the time
	_, exists := queues[guild.ID]

	if exists {
		if len(queues[guild.ID]) < MAX_QUEUE_SIZE {
			queues[guild.ID] <- play
		}
	} else {
		queues[guild.ID] = make(chan *Play, MAX_QUEUE_SIZE)
		PlaySound(s, play, nil)
	}
}

// Play a sound
func PlaySound(s *discordgo.Session, play *Play, vc *discordgo.VoiceConnection) (err error) {
	log.WithFields(log.Fields{
		"play": play,
	}).Info("Playing sound")

	if vc == nil {
		vc, err = s.ChannelVoiceJoin(play.GuildID, play.ChannelID, false, false)
		// vc.Receive = false
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Failed to play sound")
			delete(queues, play.GuildID)
			return err
		}
	}

	// If we need to change channels, do that now
	if vc.ChannelID != play.ChannelID {
		vc.ChangeChannel(play.ChannelID, false, false)
		time.Sleep(time.Millisecond * 125)
	}

	// // Track stats for this play in redis
	// go rdTrackSoundStats(play)

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(time.Millisecond * 32)

	// Play the sound
	play.Sound.Play(vc)

	// If this is chained, play the chained sound
	if play.Next != nil {
		PlaySound(s, play.Next, vc)
	}

	// If there is another song in the queue, recurse and play that
	if len(queues[play.GuildID]) > 0 {
		play := <-queues[play.GuildID]
		PlaySound(s, play, vc)
		return nil
	}

	// If the queue is empty, delete it
	time.Sleep(time.Millisecond * time.Duration(play.Sound.PartDelay))
	delete(queues, play.GuildID)
	vc.Disconnect()
	return nil
}

// func GetSoundCommands() string {
//     buffer := bytes.NewBufferString("")
//     for _, coll := range collections {
//         buffer.WriteString("**")
//         buffer.WriteString(coll.Commands[0])
//         buffer.WriteString(":** ")
//         for idx, snd := range coll.Sounds {
//             buffer.WriteString(snd.Name)
//             if(idx != len(coll.Sounds)-1) {
//                 buffer.WriteString(", ")
//             }
//         }
//         buffer.WriteString("\n")
//     }
//     return buffer.String();
// }

func FindSoundByName(base string, name string) *Sound {
	for _, c := range collections {
		if utils.Scontains(base, c.Commands...) {
			for _, s := range c.Sounds {
				if name == s.Name {
					return s
				}
			}
		}
	}
	return nil
}

func LoadSounds() {
	audio, _ := os.Open("audio")
	files, _ := audio.Readdirnames(0)
	for _, file := range files {
		parts := strings.Split(file, "_")
		prefix := parts[0]
		name := parts[1]
		for i := 2; i < len(parts); i++ {
			name = name+"_"+parts[i]
		}
		name = (strings.Split(name, "."))[0]
		coll := getCollection(prefix)
		if(coll==nil){
			log.Debug("Creating collection "+prefix+" with sound "+name)
			var NEW *SoundCollection = &SoundCollection{
				Prefix: prefix,
				Commands: []string{
					prefix,
				},
				Sounds: []*Sound{
					CreateSound(name, 50, 0),
				},
			}
			log.Debug("Created collection "+prefix)
			AddCollection(NEW)
			log.Debug("Added collection "+prefix)
			NEW.Sounds[0].Load(NEW)
			log.Debug("Loaded sound "+name)
		} else{
			log.Debug("Adding sound "+name+" to collection "+prefix)
			newSound := CreateSound(name, 50, 0)
			coll.Sounds = append(coll.Sounds, newSound)
			newSound.Load(coll)
		}
		//coll.Load()
	}
	// collections := GetCollections()
	// for _, coll := range collections {
	// 	coll.Load()
	// }
}

func GetCollections() []*SoundCollection {
	return collections
}

func getCollection(prefix string) *SoundCollection{
	for _, coll := range collections {
		if(coll.Prefix == prefix){
			return coll
		}
	}
	return nil
}


func AddCollection(newCollection *SoundCollection) error {
	collections = append(collections, newCollection)
	return nil
}

func PrintCollections() string {
	result := ""
	for _, coll := range collections {
		result += "**"+coll.Prefix+":** "
		for i, sound := range coll.Sounds {
			if i < (len(coll.Sounds) - 1) {
				result+= sound.Name+", "
			} else {
				result+= sound.Name+"\n"
			}
		}
	}
	return result
}
