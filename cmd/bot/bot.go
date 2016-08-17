package main

import (
    "flag"
    "math/rand"
    "os"
    "os/signal"
    "strconv"
    "strings"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/bwmarrin/discordgo"
    redis "gopkg.in/redis.v3"
)

var (
    // discordgo session
    discord *discordgo.Session

    // Redis client connection (used for stats)
    rcli *redis.Client

    // Sound encoding settings
    BITRATE        = 128
    MAX_QUEUE_SIZE = 6

    // Prefix for chat commands
    PREFIX = "!!"

    // Owner
    OWNER string

    // Temporary bool for enabling welcome
    WelcomeEnabled bool

)

func init() {
    // Seed the random number generator.
    rand.Seed(time.Now().UnixNano())
}

func onReady(s *discordgo.Session, event *discordgo.Ready) {
    log.Info("Recieved READY payload")
    s.UpdateStatus(0, "Dank memes")
}

func processGameplayLoop(ticker *time.Ticker) {
    for {
        select {
        case <- ticker.C:
            var processedUsers []string
            for _, g := range discord.State.Guilds {
                for _, p := range g.Presences {
                    if p.Game != nil && len(p.Game.Name) > 0 &&
                            !utilStringInSlice(p.User.ID, processedUsers) {

                        processedUsers = append(processedUsers, p.User.ID)

                        db := dbGetSession(g.ID)
                        db.GameTrackIncGameEntry(p.User.ID, p.Game.Name, 60)
                    }
                }
            }
        }
    }
}

func onGuildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
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

func airhornBomb(cid string, guild *discordgo.Guild, user *discordgo.User, cs string) {
    count, _ := strconv.Atoi(cs)
    discord.ChannelMessageSend(cid, ":ok_hand:"+strings.Repeat(":trumpet:", count))

    // Cap it at something
    if count > 100 {
        return
    }

    play := sndCreatePlay(user, guild, AIRHORN, nil)
    vc, err := discord.ChannelVoiceJoin(play.GuildID, play.ChannelID, true, true)
    if err != nil {
        return
    }

    for i := 0; i < count; i++ {
        AIRHORN.Random().Play(vc)
    }

    vc.Disconnect()
}

// Handles bot operator messages, should be refactored (lmao)
func handleBotControlMessages(s *discordgo.Session, m *discordgo.MessageCreate, parts []string, g *discordgo.Guild) {
    c,_ := s.UserChannelCreate(m.Author.ID)

    if utilScontains(parts[1], "status") {
        rdDisplayBotStats(c.ID)

    } else if utilScontains(parts[1], "stats") {
        if len(m.Mentions) >= 2 {
            rdDisplayUserStats(c.ID, utilGetMentioned(s, m).ID)
        } else if len(parts) >= 3 {
            rdDisplayUserStats(c.ID, parts[2])
        } else {
            rdDisplayServerStats(c.ID, g.ID)
        }

    } else if utilScontains(parts[1], "aps") {
        s.ChannelMessageSend(c.ID, ":ok_hand: give me a sec m8")
        go rdCalculateAirhornsPerSecond(c.ID)

    } else if utilScontains(parts[1], "toggle_welcome") {
        s.ChannelMessageSend(c.ID, ":ok_hand: give me a sec m8")
        WelcomeEnabled = !WelcomeEnabled
    }
}

func onVoiceStateUpdate(s *discordgo.Session, m *discordgo.VoiceStateUpdate) {
    if m.ChannelID == "" {
        return
    }

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
            "member":   member,
        }).Warning("Failed to grab member")
        return
    }

    if member.User.Bot {
        return
    }

    startTime := time.Date(2016, time.August, 16, 23, 0, 0, 0, time.UTC)
    endTime := time.Date(2016, time.August, 17, 5, 0, 0, 0, time.UTC)

    if utilInTimeSpan(startTime, endTime, time.Now().UTC()) {
        db := dbGetSession(guild.ID)
        e := db.GetVoiceJoinEntry(member.User.ID)

        if e != nil {
            for _, t := range e.Dates {
                p := time.Unix(t, 0).UTC()
                // log.Info(p)
                if utilInTimeSpan(startTime, endTime, p) {
                    return
                }
            }
        }

        db.UpsertVoiceJoinEntry(member.User.ID)

        // Give weekly bit bonus
        message:= giveWeeklyBitBonus(guild, member.User.ID)
        c,_ := s.UserChannelCreate(member.User.ID)
        s.ChannelMessageSend(c.ID, message)

        // Welcome them to the family
        if WelcomeEnabled {
            var sound *Sound
            for _, s := range MEMES.Sounds {
                if "welcomebdc" == s.Name {
                    sound = s
                }
            }
            go sndEnqueuePlay(member.User, guild, MEMES, sound)
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

    // If we have a message not starting with "!", then handle markov stuff
    if (!strings.HasPrefix(m.Content, "!") && len(m.Mentions) < 1) {
        mkWriteMessage(guild, m.Content)
        rando := rand.Intn(100)
        if rando < 10 {
            log.Printf("Sending markov message")
            go func() {
                s.ChannelMessageSend(m.ChannelID,
                                        mkGetMessage(guild, m.Author))
            }()
        }
        return
    }

    msg := strings.Replace(m.ContentWithMentionsReplaced(), s.State.Ready.User.Username, "username", 1)
    parts := strings.Split(strings.ToLower(msg), " ")

    // If this is a mention, it should come from the owner (otherwise we don't care)
    if len(m.Mentions) > 0 && m.Author.ID == OWNER && len(parts) > 0 {
        mentioned := false
        for _, mention := range m.Mentions {
            mentioned = (mention.ID == s.State.Ready.User.ID)
            if mentioned {
                break
            }
        }

        if mentioned {
            handleBotControlMessages(s, m, parts, guild)
        }
        return
    }

    // Filter out commands for airhornbot
    if (!strings.HasPrefix(m.Content, PREFIX)) {
        log.Printf("Filtering out airhornbot command")
        return
    }

    baseCommand := strings.Replace(parts[0], PREFIX, "", 1)

    // Process text based commands
    for _, tcoll := range TEXTCMDS {
        if utilScontains(baseCommand, tcoll.Commands...) {
            s.ChannelMessageSend(m.ChannelID,
                                    tcoll.Function(guild, m.Message, parts))
        }
    }

    // Process sound commands
    for _, coll := range COLLECTIONS {
        if utilScontains(baseCommand, coll.Commands...) {

            // If they passed a specific sound effect, find and select that (otherwise play nothing)
            var sound *Sound
            if len(parts) > 1 {
                for _, s := range coll.Sounds {
                    if parts[1] == s.Name {
                        sound = s
                    }
                }

                if sound == nil {
                    return
                }
            }

            go sndEnqueuePlay(m.Author, guild, coll, sound)
            return
        }
    }
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
        Token      = flag.String("t", "", "Discord Authentication Token")
        Redis      = flag.String("r", "", "Redis Connection String")
        Owner      = flag.String("o", "", "Owner ID")
        err        error
    )
    flag.Parse()

    if *Owner != "" {
        OWNER = *Owner
    }

    // Preload all the sounds
    log.Info("Preloading sounds...")
    for _, coll := range COLLECTIONS {
        coll.Load()
    }

    // If we got passed a redis server, try to connect
    if *Redis != "" {
        log.Info("Connecting to redis...")
        rcli = redis.NewClient(&redis.Options{Addr: *Redis, DB: 0})
        _, err = rcli.Ping().Result()

        if err != nil {
            log.WithFields(log.Fields{
                "error": err,
            }).Fatal("Failed to connect to redis")
            return
        }
    }

    // Open new database
    log.Info("Opening MongoDB")
    dbMongoOpen("localhost")

    // log.Info("Testing bits")
    db := dbGetSession("1")
    // db.SetBitStats("1", "2", 15)
    // bits := db.GetBitStats("1", "2")
    // log.Info(bits)

    // db.IncBitStats("1", "2", 20)
    // bits = db.GetBitStats("1", "2")
    // log.Info(bits)

    // db.DecBitStats("1", "2", 5)
    // bits = db.GetBitStats("1", "2")
    // log.Info(bits)

    // err = db.DecCheckBitStats("1", "2", 20)
    // bits = db.GetBitStats("1", "2")
    // log.Info(bits)

    // if err != nil {
    //     log.Error("NEB")
    // }

    // err = db.DecCheckBitStats("1", "2", 20)
    // bits = db.GetBitStats("1", "2")
    // log.Info(bits)

    // if err != nil {
    //     log.Error("NEB")
    // }

    db.GameTrackIncGameEntry("2", "Game1", 5)
    db.GameTrackIncGameEntry("2", "Game2", 6)
    db.GameTrackIncGameEntry("2", "Game3", 7)
    db.GameTrackIncGameEntry("2", "Game4", 8)


    // Create a discord session
    log.Info("Starting discord session...")
    discord, err = discordgo.New(*Token)
    if err != nil {
        log.WithFields(log.Fields{
            "error": err,
        }).Fatal("Failed to create discord session")
        return
    }

    WelcomeEnabled = false

    discord.AddHandler(onReady)
    discord.AddHandler(onGuildCreate)
    discord.AddHandler(onMessageCreate)
    discord.AddHandler(onPresenceUpdate)
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
