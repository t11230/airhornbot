package main

type TextCollection struct {
    Commands    []string
    Text        string
}

var SOUNDCOMMANDS *TextCollection = &TextCollection{
    Commands: []string{
        "imanoob",
        "l2p",
    },
    Text: sndGetSoundCommands(),
}

var GITHUB *TextCollection = &TextCollection{
    Commands: []string{
        "github",
        "git",
    },
    Text: "https://github.com/t11230/airhornbot",
}

var HILLARY *TextCollection = &TextCollection{
    Commands: []string{
        "hillary",
    },
    Text: "https://i.imgur.com/1PFAZsV.jpg",
}

var TEXTCMDS []*TextCollection = []*TextCollection{
    SOUNDCOMMANDS, 
    GITHUB,
    HILLARY,
}