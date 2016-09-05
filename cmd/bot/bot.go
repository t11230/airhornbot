package main

import (
	"flag"
	"github.com/t11230/ramenbot/lib/sound"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"

	// "github.com/t11230/ramenbot/lib/utils"
	"github.com/t11230/ramenbot/lib/config"
	"github.com/t11230/ramenbot/lib/modules"
	"github.com/t11230/ramenbot/lib/ramendb"
)

var (
	// Discordgo session
	discord *discordgo.Session

	// Prefix for chat commands
	PREFIX = "!!"
)

func init() {
	// Seed the random number generator.
	rand.Seed(time.Now().UnixNano())
}

func onReady(s *discordgo.Session, event *discordgo.Ready) {
	log.Info("Recieved READY payload")
	s.UpdateStatus(0, "Dank memes")

	modules.InitEvents(s)
}

func onGuildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	// This filters out guilds that we aren't joining for the first time
	if event.Guild.Unavailable != nil {
		return
	}

	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			s.ChannelMessageSend(channel.ID, "**RAMENBOT READY**")
			return
		}
	}
}

func onVoiceStateUpdate(s *discordgo.Session, m *discordgo.VoiceStateUpdate) {
	guild, _ := discord.State.Guild(m.GuildID)
	if guild == nil {
		log.WithFields(log.Fields{
			"guild":   m.GuildID,
			"channel": m,
		}).Warning("Failed to grab guild")
		return
	}

	member, _ := discord.State.Member(m.GuildID, m.UserID)
	if member == nil {
		log.WithFields(log.Fields{
			"member": member,
		}).Warning("Failed to grab member")
		return
	}

	if member.User.Bot {
		return
	}
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(m.Content) <= 0 {
		return
	}

	channel, _ := discord.State.Channel(m.ChannelID)
	if channel == nil {
		log.WithFields(log.Fields{
			"channel": m.ChannelID,
			"message": m.ID,
		}).Warning("Failed to grab channel")
		return
	}

	guild, _ := discord.State.Guild(channel.GuildID)
	if guild == nil {
		log.WithFields(log.Fields{
			"guild":   channel.GuildID,
			"channel": channel,
			"message": m.ID,
		}).Warning("Failed to grab guild")
		return
	}

	// // If we have a message not starting with "!", then handle markov stuff
	// if (!strings.HasPrefix(m.Content, "!") && len(m.Mentions) < 1) {
	//     mkWriteMessage(guild, m.Content)
	//     rando := rand.Intn(100)
	//     if rando < 10 {
	//         log.Printf("Sending markov message")
	//         go func() {
	//             s.ChannelMessageSend(m.ChannelID,
	//                                     mkGetMessage(guild, m.Author))
	//         }()
	//     }
	//     return
	// }

	// Filter out normal messages
	if !strings.HasPrefix(m.Content, PREFIX) {
		log.Debug("Filtering non-command")
		return
	}

	msg := strings.Replace(m.ContentWithMentionsReplaced(), s.State.Ready.User.Username, "username", 1)
	parts := strings.Split(strings.ToLower(msg), " ")
	baseCommand := strings.Replace(parts[0], PREFIX, "", 1)

	cmd := modules.Command{
		Session: s,
		Guild:   guild,
		Message: m.Message,
		Command: baseCommand,
		Args:    parts[1:],
	}

	modules.HandleCommand(&cmd)

	// // Process text based commands
	// for _, tcoll := range TEXTCMDS {
	//     if utils.Scontains(baseCommand, tcoll.Commands...) {
	//         s.ChannelMessageSend(m.ChannelID,
	//                                 tcoll.Function(guild, m.Message, parts))
	//     }
	// }

	// // Process role based commands
	// for _, rcoll := range ROLECMDS {
	//     if utils.Scontains(baseCommand, rcoll.Commands...) {
	//         s.ChannelMessageSend(m.ChannelID,
	//                                 rcoll.Function(s, guild, m.Message, parts))
	//     }
	// }

	// // Process sound commands
	// for _, coll := range COLLECTIONS {
	//     if utils.Scontains(baseCommand, coll.Commands...) {

	//         // If they passed a specific sound effect, find and select that (otherwise play nothing)
	//         var sound *Sound
	//         if len(parts) > 1 {
	//             for _, s := range coll.Sounds {
	//                 if parts[1] == s.Name {
	//                     sound = s
	//                 }
	//             }

	//             if sound == nil {
	//                 return
	//             }
	//         }

	//         go sndEnqueuePlay(m.Author, guild, coll, sound)
	//         return
	//     }
	// }
}

// Handle updating of presences in the current session, because the API doesnt...
func onPresenceUpdate(s *discordgo.Session, u *discordgo.PresenceUpdate) {
	if s == nil {
		return
	}

	guild, err := s.Guild(u.GuildID)
	if err != nil {
		return
	}

	s.Lock()
	defer s.Unlock()

	for i, p := range guild.Presences {
		if p.User.ID == u.User.ID {
			guild.Presences[i].Status = u.Status
			guild.Presences[i].Game = u.Game
			return
		}
	}

	return
}

func main() {
	var (
		verbose = flag.Bool("v", false, "Verbose")
		err     error
	)
	flag.Parse()

	if *verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.Info("Loading config file")
	conf, err := config.LoadConfig("config.json")
	if err != nil {
		log.Errorf("Invalid config file format: %v", err)
		return
	}

	log.Info("Loading Modules")
	err = modules.LoadModules(conf.Modules)
	if err != nil {
		log.Errorf("Error loading modules: ", err)
		return
	}

	// Preload all the sounds
	log.Info("Preloading sounds...")
	sound.LoadSounds()

	// Open database
	log.Info("Opening MongoDB")
	ramendb.MongoOpen(conf.MongoDB)

	// Create a discord session
	log.Info("Starting discord session...")
	discord, err = discordgo.New(conf.Token)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to create discord session")
		return
	}

	discord.AddHandler(onReady)
	discord.AddHandler(onGuildCreate)
	discord.AddHandler(onMessageCreate)
	// discord.AddHandler(onPresenceUpdate)
	discord.AddHandler(onVoiceStateUpdate)

	err = discord.Open()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to create discord websocket connection")
		return
	}

	// We're running!
	log.Info("RamenBot is ready to horn it up.")

	// log.Info("Setting up Game watch tick")
	// ticker := time.NewTicker(time.Second * 60)
	// go processGameplayLoop(ticker)

	// Wait for a signal to quit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}
