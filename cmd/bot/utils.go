package main

import (
    "regexp"
    "strings"
    "unicode"
    "unicode/utf8"
    "math/rand"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/bwmarrin/discordgo"
)

var (
    soundRegex *regexp.Regexp
    emoteRegex *regexp.Regexp
    urlRegex *regexp.Regexp
    soundBlacklist = []string {
        "!airhorn",
    }
)

func init() {
    soundRegex = regexp.MustCompile(`!(?P<group>[^\s]+)\s*(?P<effect>[^\s]*)`)
    emoteRegex = regexp.MustCompile(`<([^>]+)>`) 
    urlRegex = regexp.MustCompile(`http.*`)
}

func utilParseText(text string) string {
    soundNames := soundRegex.SubexpNames()
    soundMatches := soundRegex.FindAllStringSubmatch(text, -1)
    if len(soundMatches) > 0 {
        return ""
        soundMatch := soundMatches[0]
        soundNameMap := map[string]string{}

        for i, n := range soundMatch {
            soundNameMap[soundNames[i]] = n
            log.Printf("%s", n)
        }

        for _, b := range soundBlacklist {
            if b == soundNameMap["group"] {
                log.Printf("Excluding: %s", b)
                return ""
            }
        }
    }

    if utilScontains(text, "http"){
        log.Printf("Excluding: %s", text)
        return ""
    }

    text = emoteRegex.ReplaceAllString(text, "")
    text = urlRegex.ReplaceAllString(text, "")

    text = strings.TrimSpace(text)

    text = strings.Replace(text, "&lt;", "<", -1)
    text = strings.Replace(text, "&gt;", ">", -1)
    text = strings.Replace(text, "&amp;", "&", -1)

    return text
}

func utilLowerFirst(s string) string {
    if s == "" {
        return ""
    }
    r, n := utf8.DecodeRuneInString(s)
    return string(unicode.ToLower(r)) + s[n:]
}

func utilStringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func utilScontains(key string, options ...string) bool {
    for _, item := range options {
        if item == key {
            return true
        }
    }
    return false
}

// Attempts to find the current users voice channel inside a given guild
func utilGetCurrentVoiceChannel(user *discordgo.User, guild *discordgo.Guild) *discordgo.Channel {
    for _, vs := range guild.VoiceStates {
        if vs.UserID == user.ID {
            channel, _ := discord.State.Channel(vs.ChannelID)
            return channel
        }
    }
    return nil
}

// Returns a random integer between min and max
func utilRandomRange(min, max int) int {
    rand.Seed(time.Now().UTC().UnixNano())
    return rand.Intn(max-min) + min
}

func utilGetMentioned(s *discordgo.Session, m *discordgo.MessageCreate) *discordgo.User {
    for _, mention := range m.Mentions {
        if mention.ID != s.State.Ready.User.ID {
            return mention
        }
    }
    return nil
}