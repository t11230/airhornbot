package rolemod

import (
    log "github.com/Sirupsen/logrus"
    "github.com/t11230/ramenbot/lib/modules/modulebase"
    "github.com/bwmarrin/discordgo"
    "github.com/t11230/ramenbot/lib/perms"
    "github.com/t11230/ramenbot/lib/utils"
)

var (
    red string
    orange string
    yellow string
    green string
    blue string
    purple string
    disco string
    is_4ever bool
    is_5ever bool
)

// Module name used in the config file

const (
    ConfigName = "rolemod"

    helpString = "**!!role** : This module allows the user to choose from a set of roles.\n"

	roleHelpString = `**ROLE**
This module allows the user to choose from a set of roles.

**usage:** !!role *function* *args...*

**permissions required:** none

**functions:**
    *color:* This function allows the user to change the color of their name.

For more info on using any of these functions, type **!!role [function name] help**`

    colorHelpString = `**COLOR**

**usage:** !!role color *colorname*
    Changes *user's* name color to *colorname*.

**color names:** red, orange, yellow, green, blue, purple, disco, clear
    Use color *clear* to reset to black.
    Use color *disco* for disco party.`

    role_perms = 0x00000001 | 0x00000400 | 0x00000800 | 0x00001000 | 0x00004000 | 0x00008000 | 0x00010000 | 0x00020000 | 0x00100000 | 0x00200000
)

// List of commands that this module accepts
var commandTree = []modulebase.ModuleCommandTree{
	{
		RootCommand: "role",
		SubKeys: modulebase.SK{
			"color": modulebase.CN{
				Function:      handleChangeColor,
                Permissions:   []perms.Perm{roleControlPerm},
			},
            "help": modulebase.CN{
				Function:   handleRoleHelp,
			},
		},
	},
}

var roleControlPerm = perms.Perm{"role-control"}

// Called to initialize this module
func SetupFunc(config *modulebase.ModuleConfig) (*modulebase.ModuleSetupInfo, error) {
	return &modulebase.ModuleSetupInfo{
		Commands: &commandTree,
        Help:     helpString,
        DBStart:  handleDbStart,
	}, nil
}

func handleDbStart() error {
	perms.CreatePerm(roleControlPerm.Name)
	return nil
}

func handleRoleHelp(cmd *modulebase.ModuleCommand) (string, error) {
    return roleHelpString, nil
}

func createColors(s *discordgo.Session, guild *discordgo.Guild, m *discordgo.Message) string {
    red_exists := false
    orange_exists := false
    yellow_exists := false
    green_exists := false
    blue_exists := false
    purple_exists := false
    disco_exists := false
    for _, role := range(guild.Roles){
        if role.Name == "red" {
            red_exists = true
            red = role.ID
        }
        if role.Name == "orange" {
            orange_exists = true
            orange = role.ID
        }
        if role.Name == "yellow" {
            yellow_exists = true
            yellow = role.ID
        }
        if role.Name == "green" {
            green_exists = true
            green = role.ID
        }
        if role.Name == "blue" {
            blue_exists = true
            blue = role.ID
        }
        if role.Name == "purple" {
            purple_exists = true
            purple = role.ID
        }
        if role.Name == "disco" {
            disco_exists = true
            disco = role.ID
        }
    }
    if !red_exists {
        redrole, err := s.GuildRoleCreate(guild.ID)
        if err != nil {
            log.Error("Failed to create red role")
            return ""
        }
        redrole, err = s.GuildRoleEdit(guild.ID, redrole.ID, "red", 0xe74c3c, false, role_perms)
        if err != nil {
            log.Error("Failed to add red role")
            return ""
        }
        red = redrole.ID
    }

    if !orange_exists {
        orangerole, err := s.GuildRoleCreate(guild.ID)
        if err != nil {
            log.Error("Failed to create orange role")
            return ""
        }
        orangerole, err = s.GuildRoleEdit(guild.ID, orangerole.ID, "orange", 0xe67e22, false, role_perms)
        if err != nil {
            log.Error("Failed to add orange role")
            return ""
        }
        orange = orangerole.ID
    }

    if !yellow_exists {
        yellowrole, err := s.GuildRoleCreate(guild.ID)
        if err != nil {
            log.Error("Failed to create yellow role")
            return ""
        }
        yellowrole, err = s.GuildRoleEdit(guild.ID, yellowrole.ID, "yellow", 0xf1c40f, false, role_perms)
        if err != nil {
            log.Error("Failed to add yellow role")
            return ""
        }
        yellow = yellowrole.ID
    }

    if !green_exists {
        greenrole, err := s.GuildRoleCreate(guild.ID)
        if err != nil {
            log.Error("Failed to create green role")
            return ""
        }
        greenrole, err = s.GuildRoleEdit(guild.ID, greenrole.ID, "green", 0x2ecc71, false, role_perms)
        if err != nil {
            log.Error("Failed to add green role")
            return ""
        }
        green = greenrole.ID
    }

    if !blue_exists {
        bluerole, err := s.GuildRoleCreate(guild.ID)
        if err != nil {
            log.Error("Failed to create blue role")
            return ""
        }
        bluerole, err = s.GuildRoleEdit(guild.ID, bluerole.ID, "blue", 0x3498db, false, role_perms)
        if err != nil {
            log.Error("Failed to add blue role")
            return ""
        }
        blue = bluerole.ID
    }

    if !purple_exists {
        purplerole, err := s.GuildRoleCreate(guild.ID)
        if err != nil {
            log.Error("Failed to create purple role")
            return ""
        }
        purplerole, err = s.GuildRoleEdit(guild.ID, purplerole.ID, "purple", 0x9b59b6, false, role_perms)
        if err != nil {
            log.Error("Failed to add purple role")
            return ""
        }
        purple = purplerole.ID
    }

    if !disco_exists {
        discorole, err := s.GuildRoleCreate(guild.ID)
        if err != nil {
            log.Error("Failed to create red role")
            return ""
        }
        discorole, err = s.GuildRoleEdit(guild.ID, discorole.ID, "disco", 0xe74c3c, false, role_perms)
        if err != nil {
            log.Error("Failed to add red role")
            return ""
        }
        disco = discorole.ID
    }
    return ""
}

func handleChangeColor(cmd *modulebase.ModuleCommand) (string, error) {
    createColors(cmd.Session, cmd.Guild, cmd.Message)
    // permHandle := perms.GetPermsHandle(cmd.Guild.ID, ConfigName)
    user := cmd.Message.Author
    if len(cmd.Args) < 1 || len(cmd.Args) > 2{
        return colorHelpString, nil
    }
    log.Info("Case: "+cmd.Args[0])
    color := cmd.Args[0]
    if color == "help" {
        return colorHelpString, nil
    }
    if len(cmd.Args) == 2 {
        time := cmd.Args[1]
        ever:=0
        if time == "4ever" {
            ever = 4
        }
        changeColor(cmd.Session, cmd.Guild, user, color, ever)
        return "", nil
    }
    changeColor(cmd.Session, cmd.Guild, user, color, 0)
    return "", nil
}

func changeColor(s *discordgo.Session, guild *discordgo.Guild, user *discordgo.User, c string, ever int) string {
    var err error
    var color string
    colors := []string{red, orange, yellow, green, blue, purple, disco}
    member := utils.GetMember(guild, user.ID)
    switch c {
    case "red":
        color = red
    case "orange":
        color = orange
    case "yellow":
        color = yellow
    case "green":
        color = green
    case "blue":
        color = blue
    case "purple":
        color = purple
    case "disco":
        color = disco
        if ever==4 {
            go discoParty4ever(s, guild, user)
        } else{
            go discoParty(s, guild, user)
        }
    case "clear":
        for i, role := range(member.Roles){
            if utils.Scontains(role, colors...) {
                member.Roles = append(member.Roles[:i], member.Roles[i+1:]...)
            }
        }
        err = s.GuildMemberEdit(guild.ID, user.ID, member.Roles)
        if err != nil {
            log.Error("Failed to update user's role")
            return ""
        }
        is_4ever = false
        return ""
    }
    role, err := s.State.Role(guild.ID, color)
    if err != nil {
        log.Error("Failed to change disco role")
        return ""
    }

    log.Info(role.Color)
    for i, role := range(member.Roles){
        if utils.Scontains(role, colors...) {
            member.Roles = append(member.Roles[:i], member.Roles[i+1:]...)
        }
    }
    member.Roles = append(member.Roles, role.ID)
    err = s.GuildMemberEdit(guild.ID, user.ID, member.Roles)
    if err != nil {
        log.Error("Failed to update user's role")
        return ""
    }
    return ""
}

func discoParty(s *discordgo.Session, guild *discordgo.Guild, user *discordgo.User) string {
    for i := 0; i < 30; i++ {
        s.GuildRoleEdit(guild.ID, disco, "disco", 0xe74c3c, false, role_perms)
        s.GuildRoleEdit(guild.ID, disco, "disco", 0xe67e22, false, role_perms)
        s.GuildRoleEdit(guild.ID, disco, "disco", 0xf1c40f, false, role_perms)
        s.GuildRoleEdit(guild.ID, disco, "disco", 0x2ecc71, false, role_perms)
        s.GuildRoleEdit(guild.ID, disco, "disco", 0x3498db, false, role_perms)
        s.GuildRoleEdit(guild.ID, disco, "disco", 0x9b59b6, false, role_perms)
    }
    changeColor(s, guild, user, "clear", 0)
    return ""
}

func discoParty4ever(s *discordgo.Session, guild *discordgo.Guild, user *discordgo.User) string {
    is_4ever = true
    for is_4ever{
        s.GuildRoleEdit(guild.ID, disco, "disco", 0xe74c3c, false, role_perms)
        s.GuildRoleEdit(guild.ID, disco, "disco", 0xe67e22, false, role_perms)
        s.GuildRoleEdit(guild.ID, disco, "disco", 0xf1c40f, false, role_perms)
        s.GuildRoleEdit(guild.ID, disco, "disco", 0x2ecc71, false, role_perms)
        s.GuildRoleEdit(guild.ID, disco, "disco", 0x3498db, false, role_perms)
        s.GuildRoleEdit(guild.ID, disco, "disco", 0x9b59b6, false, role_perms)
    }
    return ""
}
