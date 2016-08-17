package main

import (
    "github.com/bwmarrin/discordgo"
)

type RoleFunction func(*discordgo.Session, *discordgo.Guild, *discordgo.Message, []string) string

type RoleCollection struct {
    Commands    []string
    Function    RoleFunction
}

var COLOR *RoleCollection = &RoleCollection{
    Commands: []string{
        "red",
        "orange",
        "yellow",
        "green",
        "blue",
        "purple",
        "disco",
        "clear",
    },
    Function: changeColor,
}

var COLORCREATE *RoleCollection = &RoleCollection{
    Commands: []string{
        "create_colors",
    },
    Function: createColors,
}

var ROLECMDS []*RoleCollection = []*RoleCollection{
    COLOR,
    COLORCREATE,
}
