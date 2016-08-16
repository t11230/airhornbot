package main

import (
    log "github.com/Sirupsen/logrus"

    "github.com/bwmarrin/discordgo"
)


func changeColor(s *discordgo.Session, guild *discordgo.Guild, message *discordgo.Message, args []string) string {
    red :=  "214978951339048961"
    orange := "215238415866658816"
    yellow := "215238471558627329"
    green := "215238500235214848"
    blue := "215238540118851585"
    purple := "215238568820473856"
    colors := []string{red, orange, yellow, green, blue, purple}
    var err error
    var color string
    user := message.Author
    member := utilGetMember(guild, user.ID)
    disco := false
    log.Info("Case: "+args[0])
    switch args[0] {
    case "!!red":
        color = red
    case "!!orange":
        color = orange
    case "!!yellow":
        color = yellow
    case "!!green":
        color = green
    case "!!blue":
        color = blue
    case "!!purple":
        color = purple
    case "!!disco":
         color = red
         disco = true
    case "!!clear":
        for i, role := range(member.Roles){
            if utilStringInSlice(role, colors) {
                member.Roles = append(member.Roles[:i], member.Roles[i+1:]...)
            }
        }
        err = s.GuildMemberEdit(guild.ID, user.ID, member.Roles)
        if err != nil {
            log.Error("Failed to update user's role")
            return ""
        }
        return ""
    }
    role, err := s.State.Role(guild.ID, color)
    if err != nil {
        log.Error("Failed to change disco role")
        return ""
    }
    log.Info(role.Color)
    for i, role := range(member.Roles){
        if utilStringInSlice(role, colors) {
            member.Roles = append(member.Roles[:i], member.Roles[i+1:]...)
        }
    }
    member.Roles = append(member.Roles, role.ID)
    err = s.GuildMemberEdit(guild.ID, user.ID, member.Roles)
    if err != nil {
        log.Error("Failed to update user's role")
        return ""
    }
    if disco {
        go discoParty(s, guild, message)
    }
    return ""
}

func discoParty(s *discordgo.Session, guild *discordgo.Guild, message *discordgo.Message) string {
    red := []string{"!!red"}
    orange := []string{"!!orange"}
    yellow := []string{"!!yellow"}
    green := []string{"!!green"}
    blue := []string{"!!blue"}
    purple := []string{"!!purple"}
    clear := []string{"!!clear"}
    for i := 0; i < 30; i++ {
        changeColor(s, guild, message, red)
        changeColor(s, guild, message, orange)
        changeColor(s, guild, message, yellow)
        changeColor(s, guild, message, green)
        changeColor(s, guild, message, blue)
        changeColor(s, guild, message, purple)
    }
    changeColor(s, guild, message, clear)
    return ""
}
