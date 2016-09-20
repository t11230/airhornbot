package greeter

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/t11230/ramenbot/lib/modules/modulebase"
	"github.com/t11230/ramenbot/lib/perms"
	"github.com/t11230/ramenbot/lib/ramendb"
	"github.com/t11230/ramenbot/lib/sound"
	"github.com/t11230/ramenbot/lib/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

// Module name used in the config file
const (
	ConfigName = "greeter"
	helpString = "**!!greet** : This module allows the user to set text or voice greetings for joining a call.\n"

	greetHelpString = `**GREET**
This module allows the user to set text or voice greetings for joining a call.

**usage:** !!greet *function* *args...*

**permissions required:** greet-control

**functions:**
    *pm:* This function allows the user to set a welcome private message to be sent by the bot.
	*voice:* This function allows the user to set a welcome sound to be played by the bot.

For more info on using any of these functions, type **!!greet [function name] help**`

    pmHelpString = `**PM**

**usage:** !!greet pm *action* *<message>*
    Handles the welcome message sent by the bot to users who join the call.

**action names:** enable, disable, set
    **enable** turns on the greeting (ignores *<message>*)
	**disable** turns off the greeting (ignores *<message>*)
    **set** sets the greeting to *<message>*`

	voiceHelpString = `**VOICE**

**usage:** !!greet voice *action* *<collection>* *<sound>*
    Handles the welcome sound played by the bot when users join the call.

**action names:** enable, disable, set
    **enable** turns on the greeting (ignores *<collection>* *<sound>*)
	**disable** turns off the greeting (ignores *<collection>* *<sound>*)
    **set** sets the greeting sound to *<collection>* *<sound>* (see !!s help)`
)

// List of commands that this module accepts
var commandTree = []modulebase.ModuleCommandTree{
	{
		RootCommand: "greet",
		SubKeys: modulebase.SK{
			"pm": modulebase.CN{
				Function:    handleGreetPm,
				Permissions: []perms.Perm{greetControlPerm},
			},
			"voice": modulebase.CN{
				Function:    handleGreetVoice,
				Permissions: []perms.Perm{greetControlPerm},
			},
		},
		Function: handleGreet,
	},
}

var greetControlPerm = perms.Perm{"greet-control"}

// Called to initialize this module
func SetupFunc(config *modulebase.ModuleConfig) (*modulebase.ModuleSetupInfo, error) {
	events := []interface{}{
		voiceStateUpdateCallback,
		guildCreateCallback,
	}

	return &modulebase.ModuleSetupInfo{
		Events:   &events,
		Commands: &commandTree,
		DBStart:  handleDbStart,
		Help:     helpString,
	}, nil
}

func handleDbStart() error {
	err := perms.CreatePerm(greetControlPerm.Name)
	if err != nil {
		log.Errorf("Error creating perm: %v", err)
		return err
	}
	return nil
}

func handleGreet(cmd *modulebase.ModuleCommand) (string, error) {
	log.Debug("Called greet")
	return greetHelpString, nil
}

func handleGreetPm(cmd *modulebase.ModuleCommand) (string, error) {
	if len(cmd.Args) == 0 || cmd.Args[0] == "help" {
		return pmHelpString, nil
	}

	c := greeterCollection{ramendb.GetCollection(cmd.Guild.ID, ConfigName)}
	if cmd.Args[0] == "enable" {
		c.PMGreetEnable(cmd.Guild.ID, true)
	} else if cmd.Args[0] == "disable" {
		c.PMGreetEnable(cmd.Guild.ID, false)
	} else if cmd.Args[0] == "set" {
		c.SetPMGreetMessage(cmd.Guild.ID, strings.Join(cmd.Args[1:], " "))
	} else {
		return pmHelpString, nil
	}
	return "Updated greet pm config", nil
}

func handleGreetPmError(cmd *modulebase.ModuleCommand, e error) {
	cmd.Session.ChannelMessageSend(cmd.Message.ChannelID, fmt.Sprintf("Err: %v", e))
}

func handleGreetVoice(cmd *modulebase.ModuleCommand) (string, error) {
	if len(cmd.Args) == 0 || cmd.Args[0] == "help" {
		return voiceHelpString, nil
	}

	c := greeterCollection{ramendb.GetCollection(cmd.Guild.ID, ConfigName)}
	if cmd.Args[0] == "enable" {
		c.VoiceGreetEnable(cmd.Guild.ID, true)
	} else if cmd.Args[0] == "disable" {
		c.VoiceGreetEnable(cmd.Guild.ID, false)
	} else if cmd.Args[0] == "set" {
		if len(cmd.Args) != 3 {
			return voiceHelpString, nil
		}
		if sound.FindSoundByName(cmd.Args[1], cmd.Args[2]) == nil {
			errString := fmt.Sprintf("Invalid Sound effect: %v", cmd.Args[1:3])
			return errString, nil
		}
		c.SetVoiceGreetSound(cmd.Guild.ID, strings.Join(cmd.Args[1:3], " "))
	} else {
		return voiceHelpString, nil
	}
	return "Updated greet voice config", nil
}

// Called in response to a VoiceStateUpdate event
func voiceStateUpdateCallback(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	log.Debugf("Greeter On voice state update: %v", v.VoiceState)

	// Check if it was a part
	if v.ChannelID == "" {
		return
	}

	guild, _ := s.State.Guild(v.GuildID)
	if guild == nil {
		log.WithFields(log.Fields{
			"guild": v.GuildID,
		}).Warning("Failed to grab guild")
		return
	}

	member, _ := s.State.Member(v.GuildID, v.UserID)
	if member == nil {
		log.WithFields(log.Fields{
			"member": member,
		}).Warning("Failed to grab member")
		return
	}

	if member.User.Bot {
		return
	}
	if v.VoiceState.SelfMute == true ||  v.VoiceState.SelfDeaf == true {
		return
	}
	c := greeterCollection{ramendb.GetCollection(v.GuildID, ConfigName)}
	voiceGreet, pmGreet := c.GreetEnabled(v.GuildID)

	// Handle PM greets
	if pmGreet {
		message := fmt.Sprintf(c.PMGreetMessage(v.GuildID),
			utils.GetPreferredName(guild, v.UserID))

		channel, _ := s.UserChannelCreate(v.UserID)
		s.ChannelMessageSend(channel.ID, message)
	}
	// Handle Voice greets
	if voiceGreet {
		log.Debugf("Greeting :%v", member.User)
		name := strings.Split(c.VoiceGreetSound(v.GuildID), " ")
		if len(name) != 2 {
			return
		}
		snd := sound.FindSoundByName(name[0], name[1])
		go sound.EnqueuePlay(s, member.User, guild, nil, snd)
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
	GuildID         string
	VoiceGreet      *bool  `bson:",omitempty"`
	VoiceGreetSound string `bson:",omitempty"`
	PMGreet         *bool  `bson:",omitempty"`
	PMGreetMessage  string `bson:",omitempty"`
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
		GuildID:         guildId,
		VoiceGreet:      &[]bool{false}[0],
		VoiceGreetSound: "meme welcomebdc",
		PMGreet:         &[]bool{false}[0],
		PMGreetMessage:  "Welcome to the server %s!",
	}
	c.Insert(defaultConfig)
}

func (c *greeterCollection) GreetEnabled(guildId string) (bool, bool) {
	config := greeterConfig{}
	c.Find(greeterConfig{GuildID: guildId}).One(&config)

	voiceGreet := false
	pmGreet := false
	if config.VoiceGreet != nil {
		voiceGreet = *config.VoiceGreet
	}
	if config.PMGreet != nil {
		pmGreet = *config.PMGreet
	}
	return voiceGreet, pmGreet
}

func (c *greeterCollection) VoiceGreetEnable(guildId string, enable bool) {
	upsertdata := bson.M{"$set": greeterConfig{
		GuildID:    guildId,
		VoiceGreet: &enable,
	}}

	c.Update(greeterConfig{GuildID: guildId}, upsertdata)
}

func (c *greeterCollection) PMGreetEnable(guildId string, enable bool) {
	upsertdata := bson.M{"$set": greeterConfig{
		GuildID: guildId,
		PMGreet: &enable,
	}}

	c.Update(greeterConfig{GuildID: guildId}, upsertdata)
}

func (c *greeterCollection) PMGreetMessage(guildId string) string {
	config := greeterConfig{}
	c.Find(greeterConfig{GuildID: guildId}).One(&config)

	return config.PMGreetMessage
}

func (c *greeterCollection) SetPMGreetMessage(guildId string, message string) {
	data := bson.M{"$set": greeterConfig{
		GuildID:        guildId,
		PMGreetMessage: message,
	}}

	c.Update(greeterConfig{GuildID: guildId}, data)
}

func (c *greeterCollection) VoiceGreetSound(guildId string) string {
	config := greeterConfig{}
	c.Find(greeterConfig{GuildID: guildId}).One(&config)

	return config.VoiceGreetSound
}

func (c *greeterCollection) SetVoiceGreetSound(guildId string, name string) {
	data := bson.M{"$set": greeterConfig{
		GuildID:         guildId,
		VoiceGreetSound: name,
	}}

	c.Update(greeterConfig{GuildID: guildId}, data)
}
