package main

import (
	"flag"
	"github.com/t11230/ramenbot/lib/sound"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"time"
	"syscall"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"

	"github.com/t11230/ramenbot/lib/config"
	"github.com/t11230/ramenbot/lib/modules"
	"github.com/t11230/ramenbot/lib/ramendb"
)

var (
	// Discordgo session
	discord *discordgo.Session

	// Prefix for chat commands
	PREFIX = "!!"
	buildstamp = ""
	githash = ""
)

const keikaku = "**TRANSLATOR'S NOTE:** Keikaku means plan"
func init() {
	// Seed the random number generator.
	rand.Seed(time.Now().UnixNano())
}

func onReady(s *discordgo.Session, event *discordgo.Ready) {
	log.Info("Recieved READY payload")
	s.UpdateStatus(0, "Dank memes")

	modules.InitEvents(s)
}

func handlePatchUpdate() {
	var err error
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGUSR1, syscall.SIGUSR2)
    sig := <-sigs
    if sig == syscall.SIGUSR1 {
		err = syscall.Exec("./patch", []string{}, nil)
	}
	if sig == syscall.SIGUSR2 {
		err = syscall.Exec("./patch_verbose", []string{}, nil)
	} else {
		log.Error("Unrecognized Signal")
	}
	if err!=nil{
		log.Error("Exec failed in handlePatchUpdate")
	}
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

	// Handle version check
	if strings.Compare(m.Content, "--version")==0 {
		gitstring := fmt.Sprintf("**Git Commit Hash:** %s\n", githash)
		buildstring := fmt.Sprintf("**UTC Build Time:** %s\n", buildstamp)
		s.ChannelMessageSend(channel.ID, gitstring)
		s.ChannelMessageSend(channel.ID, buildstring)
		return
	}

	// Handle massages with Keikaku in them
	if strings.Contains(m.Content, "keikaku") {
		s.ChannelMessageSend(channel.ID, keikaku)
		return
	}
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
}

func main() {
	var (
		version = flag.Bool("version", false, "Version")
		verbose = flag.Bool("v", false, "Verbose")
		patch = flag.Bool("p", false, "Patching")
		err     error
	)
	flag.Parse()

	if *version {
		fmt.Printf("Git Commit Hash: %s\n", githash)
		fmt.Printf("UTC Build Time: %s\n", buildstamp)
		return
	}
	if *verbose {
		log.SetLevel(log.DebugLevel)
	}
	if *patch {
		if *verbose {
			err = syscall.Exec("./patch_verbose",[]string{},nil)
		} else {
			err = syscall.Exec("./patch", []string{}, nil)
		}
		if err != nil {
			log.Errorf("Patching failed: %v", err)
		}
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

	log.Info("Processing Module DB start functions")
	err = modules.DBHooks()
	if err != nil {
		log.Errorf("Error on modules DB start: ", err)
		return
	}

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

	err = discord.Open()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to create discord websocket connection")
		return
	}

	// We're running!
	log.Info("RamenBot is running.")

	//
	go handlePatchUpdate()
	// Wait for a signal to quit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}
