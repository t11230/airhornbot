package greeter

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/t11230/ramenbot/lib/modules/modulebase"
	"github.com/t11230/ramenbot/lib/ramendb"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Module name used in the config file
const ConfigName = "greeter"

// List of commands that this module accepts
var commandList = []modulebase.ModuleCommandListener{
	{Command: "greet", Callback: greetCallback},
}

// Called to initialize this module
func SetupFunc(config *modulebase.ModuleConfig) (*modulebase.ModuleSetupInfo, error) {
	// Subscribe to VoiceStateUpdate events
	events := []interface{}{
		voiceStateUpdateCallback,
	}

	return &modulebase.ModuleSetupInfo{
		Events:   &events,
		Commands: &commandList,
	}, nil
}

// Called when the greet command is seen by the bot
func greetCallback(s *discordgo.Session, cmd *modulebase.ModuleCommand) error {
	log.Debugf("Greeter command: %v", cmd.Args)
	switch cmd.Args[0] {
	case "voice":
		c := greeterCollection{ramendb.GetCollection(cmd.Guild.ID, ConfigName)}
		if cmd.Args[1] == "enable" {
			c.VoiceGreetEnable(cmd.Guild.ID, true)
		} else {
			c.VoiceGreetEnable(cmd.Guild.ID, false)
		}
	case "pm":
		c := greeterCollection{ramendb.GetCollection(cmd.Guild.ID, ConfigName)}
		if cmd.Args[1] == "enable" {
			c.PMGreetEnable(cmd.Guild.ID, true)
		} else {
			c.PMGreetEnable(cmd.Guild.ID, false)
		}
	}
	return nil
}

// Called in response to a VoiceStateUpdate event
func voiceStateUpdateCallback(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	log.Debugf("Greeter On voice state update: %v", v.VoiceState)
	// Check if it was a part
	if v.ChannelID == "" {
		return
	}

	c := greeterCollection{ramendb.GetCollection(v.GuildID, ConfigName)}
	voiceGreet, _ := c.GreetEnabled(v.GuildID)
	if voiceGreet {
		log.Info("Greet enabled")
	} else {
		log.Info("Greet disabled")
	}
}

func guildCreateCallback(s *discordgo.Session, g *discordgo.GuildCreate) {
	if g.Unavailable != nil && *g.Unavailable == true {
		return
	}

	c := greeterCollection{ramendb.GetCollection(g.ID, ConfigName)}
	c.CreateConfig(g.ID)
}

// Database functionality
type greeterConfig struct {
	GuildID    string
	VoiceGreet *bool `bson:",omitempty"`
	PMGreet    *bool `bson:",omitempty"`
}

type greeterCollection struct {
	*mgo.Collection
}

func (c *greeterCollection) CreateConfig(guildId string) {
	count, _ := c.Find(greeterConfig{GuildID: guildId}).Count()
	if count > 0 {
		return
	}

	log.Debug("Creating new config")

	// Setup default values
	defaultConfig := greeterConfig{
		GuildID:    guildId,
		VoiceGreet: &[]bool{false}[0],
		PMGreet:    &[]bool{false}[0],
	}
	c.Insert(defaultConfig)
}

func (c *greeterCollection) GreetEnabled(guildId string) (bool, bool) {
	config := greeterConfig{}
	c.Find(greeterConfig{GuildID: guildId}).One(&config)

	return *config.VoiceGreet, false
}

func (c *greeterCollection) VoiceGreetEnable(guildId string, enable bool) {
	upsertdata := bson.M{"$set": greeterConfig{
		GuildID:    guildId,
		VoiceGreet: &enable,
	}}

	c.Upsert(greeterConfig{GuildID: guildId}, upsertdata)
}

func (c *greeterCollection) PMGreetEnable(guildId string, enable bool) {
	upsertdata := bson.M{"$set": greeterConfig{
		GuildID: guildId,
		PMGreet: &enable,
	}}

	c.Upsert(greeterConfig{GuildID: guildId}, upsertdata)
}

// import (
//     "encoding/gob"
//     "fmt"
//     "math/rand"
//     "os"
//     "strings"
//     "sync"
// )

// // Welcome them to the family
// if WelcomeEnabled {
//     var sound *Sound
//     for _, s := range MEMES.Sounds {
//         if "welcomebdc" == s.Name {
//             sound = s
//         }
//     }
//     go sndEnqueuePlay(member.User, guild, MEMES, sound)
// }
