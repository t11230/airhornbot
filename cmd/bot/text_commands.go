package main

import (
    "github.com/bwmarrin/discordgo"
)

type TextFunction func(*discordgo.Guild, *discordgo.User, []string) string

type TextCollection struct {
    Commands    []string
    Function    TextFunction
}

var GITHUB *TextCollection = &TextCollection{
    Commands: []string{
        "github",
        "git",
    },
    Function: func (guild *discordgo.Guild,
                    user *discordgo.User,
                    args []string) string {
        return "https://github.com/t11230/airhornbot"
    },
}

var SOUNDCOMMANDS *TextCollection = &TextCollection{
    Commands: []string{
        "imanoob",
        "l2p",
    },
    Function: func (*discordgo.Guild, *discordgo.User, []string) string {
        return sndGetSoundCommands()
    },
}

var HILLARY *TextCollection = &TextCollection{
    Commands: []string{
        "hillary",
    },
    Function: func (*discordgo.Guild, *discordgo.User, []string) string {
        return "https://i.imgur.com/1PFAZsV.jpg"
    },
}

var MARKOV *TextCollection = &TextCollection{
    Commands: []string{
        "text",
    },
    Function: func (guild *discordgo.Guild,
                    user *discordgo.User,
                    args []string) string {
        return mkGetMessage(guild, user)
    },
}

var STATS *TextCollection = &TextCollection{
    Commands: []string{
        "stats",
    },
    Function: gpPrintStats,
}

var TEXTCMDS []*TextCollection = []*TextCollection{
    SOUNDCOMMANDS, 
    GITHUB,
    HILLARY,
    MARKOV,
    STATS,
}