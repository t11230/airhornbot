package rolemod

// import (
//     log "github.com/Sirupsen/logrus"

//     "github.com/bwmarrin/discordgo"

//     "github.com/t11230/ramenbot/lib/utils"
// )

// var (
//     red string
//     orange string
//     yellow string
//     green string
//     blue string
//     purple string
// )
// func createColors(s *discordgo.Session, guild *discordgo.Guild, message *discordgo.Message, args []string) string {
//     red_exists := false
//     orange_exists := false
//     yellow_exists := false
//     green_exists := false
//     blue_exists := false
//     purple_exists := false
//     perms := 0x00000001 | 0x00000400 | 0x00000800 | 0x00001000 | 0x00004000 | 0x00008000 | 0x00010000 | 0x00020000 | 0x00100000 | 0x00200000
//     for _, role := range(guild.Roles){
//         if role.Name == "red" {
//             red_exists = true
//             red = role.ID
//         } else if role.Name == "orange" {
//             orange_exists = true
//             orange = role.ID
//         } else if role.Name == "yellow" {
//             yellow_exists = true
//             yellow = role.ID
//         }
//         if role.Name == "green" {
//             green_exists = true
//             green = role.ID
//         }
//         if role.Name == "blue" {
//             blue_exists = true
//             blue = role.ID
//         }
//         if role.Name == "purple" {
//             purple_exists = true
//             purple = role.ID
//         }
//     }
//     if !red_exists {
//         redrole, err := s.GuildRoleCreate(guild.ID)
//         if err != nil {
//             log.Error("Failed to create red role")
//             return ""
//         }
//         redrole, err = s.GuildRoleEdit(guild.ID, redrole.ID, "red", 0xe74c3c, false, perms)
//         if err != nil {
//             log.Error("Failed to add red role")
//             return ""
//         }
//         red = redrole.ID
//     }

//     if !orange_exists {
//         orangerole, err := s.GuildRoleCreate(guild.ID)
//         if err != nil {
//             log.Error("Failed to create orange role")
//             return ""
//         }
//         orangerole, err = s.GuildRoleEdit(guild.ID, orangerole.ID, "orange", 0xe67e22, false, perms)
//         if err != nil {
//             log.Error("Failed to add orange role")
//             return ""
//         }
//         orange = orangerole.ID
//     }

//     if !yellow_exists {
//         yellowrole, err := s.GuildRoleCreate(guild.ID)
//         if err != nil {
//             log.Error("Failed to create yellow role")
//             return ""
//         }
//         yellowrole, err = s.GuildRoleEdit(guild.ID, yellowrole.ID, "yellow", 0xf1c40f, false, perms)
//         if err != nil {
//             log.Error("Failed to add yellow role")
//             return ""
//         }
//         yellow = yellowrole.ID
//     }

//     if !green_exists {
//         greenrole, err := s.GuildRoleCreate(guild.ID)
//         if err != nil {
//             log.Error("Failed to create green role")
//             return ""
//         }
//         greenrole, err = s.GuildRoleEdit(guild.ID, greenrole.ID, "green", 0x2ecc71, false, perms)
//         if err != nil {
//             log.Error("Failed to add green role")
//             return ""
//         }
//         green = greenrole.ID
//     }

//     if !blue_exists {
//         bluerole, err := s.GuildRoleCreate(guild.ID)
//         if err != nil {
//             log.Error("Failed to create blue role")
//             return ""
//         }
//         bluerole, err = s.GuildRoleEdit(guild.ID, bluerole.ID, "blue", 0x3498db, false, perms)
//         if err != nil {
//             log.Error("Failed to add blue role")
//             return ""
//         }
//         blue = bluerole.ID
//     }

//     if !purple_exists {
//         purplerole, err := s.GuildRoleCreate(guild.ID)
//         if err != nil {
//             log.Error("Failed to create purple role")
//             return ""
//         }
//         purplerole, err = s.GuildRoleEdit(guild.ID, purplerole.ID, "purple", 0x9b59b6, false, perms)
//         if err != nil {
//             log.Error("Failed to add purple role")
//             return ""
//         }
//         purple = purplerole.ID
//     }
//     return ""
// }

// func changeColor(s *discordgo.Session, guild *discordgo.Guild, message *discordgo.Message, args []string) string {
//     colors := []string{red, orange, yellow, green, blue, purple}
//     var err error
//     var color string
//     user := message.Author
//     member := utils.GetMember(guild, user.ID)
//     disco := false
//     log.Info("Case: "+args[0])
//     switch args[0] {
//     case "!!red":
//         color = red
//     case "!!orange":
//         color = orange
//     case "!!yellow":
//         color = yellow
//     case "!!green":
//         color = green
//     case "!!blue":
//         color = blue
//     case "!!purple":
//         color = purple
//     case "!!disco":
//          color = red
//          disco = true
//     case "!!clear":
//         for i, role := range(member.Roles){
//             if utils.Scontains(role, colors...) {
//                 member.Roles = append(member.Roles[:i], member.Roles[i+1:]...)
//             }
//         }
//         err = s.GuildMemberEdit(guild.ID, user.ID, member.Roles)
//         if err != nil {
//             log.Error("Failed to update user's role")
//             return ""
//         }
//         return ""
//     }
//     role, err := s.State.Role(guild.ID, color)
//     if err != nil {
//         log.Error("Failed to change disco role")
//         return ""
//     }
//     log.Info(role.Color)
//     for i, role := range(member.Roles){
//         if utils.Scontains(role, colors...) {
//             member.Roles = append(member.Roles[:i], member.Roles[i+1:]...)
//         }
//     }
//     member.Roles = append(member.Roles, role.ID)
//     err = s.GuildMemberEdit(guild.ID, user.ID, member.Roles)
//     if err != nil {
//         log.Error("Failed to update user's role")
//         return ""
//     }
//     if disco {
//         go discoParty(s, guild, message)
//     }
//     return ""
// }

// func discoParty(s *discordgo.Session, guild *discordgo.Guild, message *discordgo.Message) string {
//     red := []string{"!!red"}
//     orange := []string{"!!orange"}
//     yellow := []string{"!!yellow"}
//     green := []string{"!!green"}
//     blue := []string{"!!blue"}
//     purple := []string{"!!purple"}
//     clear := []string{"!!clear"}
//     for i := 0; i < 30; i++ {
//         changeColor(s, guild, message, red)
//         changeColor(s, guild, message, orange)
//         changeColor(s, guild, message, yellow)
//         changeColor(s, guild, message, green)
//         changeColor(s, guild, message, blue)
//         changeColor(s, guild, message, purple)
//     }
//     changeColor(s, guild, message, clear)
//     return ""
// }
