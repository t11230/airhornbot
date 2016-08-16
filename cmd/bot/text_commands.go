package main

import (
    "github.com/bwmarrin/discordgo"
)

type TextFunction func(*discordgo.Guild, *discordgo.Message, []string) string

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
                    user *discordgo.Message,
                    args []string) string {
        return "https://github.com/t11230/airhornbot"
    },
}

var SOUNDCOMMANDS *TextCollection = &TextCollection{
    Commands: []string{
        "imanoob",
        "l2p",
    },
    Function: func (*discordgo.Guild, *discordgo.Message, []string) string {
        return sndGetSoundCommands()
    },
}

var HILLARY *TextCollection = &TextCollection{
    Commands: []string{
        "hillary",
    },
    Function: func (*discordgo.Guild, *discordgo.Message, []string) string {
        return "https://i.imgur.com/1PFAZsV.jpg"
    },
}

var MARKOV *TextCollection = &TextCollection{
    Commands: []string{
        "text",
    },
    Function: func (guild *discordgo.Guild,
                    message *discordgo.Message,
                    args []string) string {
        return mkGetMessage(guild, message.Author)
    },
}

var STATS *TextCollection = &TextCollection{
    Commands: []string{
        "stats",
    },
    Function: gpHandleStatsCommand,
}

var BITS *TextCollection = &TextCollection{
    Commands: []string{
        "bits",
    },
    Function: bitsPrintStats,
}

var DICEROLL *TextCollection = &TextCollection{
    Commands: []string{
        "roll",
    },
    Function: rollDice,
}

var BETROLL *TextCollection = &TextCollection{
    Commands: []string{
        "betroll",
    },
    Function: betRoll,
}

var BID *TextCollection = &TextCollection{
    Commands: []string{
        "bid",
    },
    Function: bid,
}


var TEXTCMDS []*TextCollection = []*TextCollection{
    SOUNDCOMMANDS,
    GITHUB,
    HILLARY,
    MARKOV,
    STATS,
    BITS,
    DICEROLL,
    BETROLL,
    BID,
}
