package admin

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/t11230/ramenbot/lib/modules/modulebase"
	"github.com/t11230/ramenbot/lib/perms"
	"github.com/t11230/ramenbot/lib/utils"
	"strings"
	"text/tabwriter"
)

const ConfigName = "admin"

var commandTree = []modulebase.ModuleCommandTree{
	{
		RootCommand: "adm",
		SubKeys: modulebase.SK{
			"addperm": modulebase.CN{
				Function: handleAddPerm,
			},
			"delperm": modulebase.CN{
				Function: handleDelPerm,
			},
			"showperms": modulebase.CN{
				Function: handleShowPerms,
			},
		},
	},
}

// Called to initialize this module
func SetupFunc(config *modulebase.ModuleConfig) (*modulebase.ModuleSetupInfo, error) {
	return &modulebase.ModuleSetupInfo{
		Events:   nil,
		Commands: &commandTree,
	}, nil
}

func handleAddPerm(cmd *modulebase.ModuleCommand) (string, error) {
	log.Debug("Handling addperm command")

	if !canManageRoles(cmd.Guild, cmd.Message.Author.ID) {
		return "Insufficient permissions", nil
	}

	if len(cmd.Args) != 2 {
		return "Invalid or missing arguments", nil
	}

	h := perms.GetPermsHandle(cmd.Guild.ID, ConfigName)
	permExist, err := perms.PermExists(cmd.Args[0])
	if err != nil {
		return "Error checking perm", nil
	}
	if !permExist {
		return "Invalid perm", nil
	}

	userName := strings.Join(cmd.Args[1:], " ")
	user, err := utils.FindUser(cmd.Guild, userName)
	if err != nil {
		return "Unable to find user", nil
	}

	if h.CheckPerm(user.ID, cmd.Args[0]) {
		return "User already has that perm", nil
	}

	err = h.AddPerm(user.ID, cmd.Args[0])
	if err != nil {
		return "Error adding perm to user", nil
	}

	return "Permission added", nil
}

func handleDelPerm(cmd *modulebase.ModuleCommand) (string, error) {
	log.Debug("Handling delperm command")

	if !canManageRoles(cmd.Guild, cmd.Message.Author.ID) {
		return "Insufficient permissions", nil
	}

	if len(cmd.Args) != 2 {
		return "Invalid or missing arguments", nil
	}

	h := perms.GetPermsHandle(cmd.Guild.ID, ConfigName)
	permExist, err := perms.PermExists(cmd.Args[0])
	if err != nil {
		return "Error checking perm", nil
	}
	if !permExist {
		return "Invalid perm", nil
	}

	userName := strings.Join(cmd.Args[1:], " ")
	user, err := utils.FindUser(cmd.Guild, userName)
	if err != nil {
		return "Unable to find user", nil
	}

	if !h.CheckPerm(user.ID, cmd.Args[0]) {
		return "User does not have that permission", nil
	}

	err = h.RemovePerm(user.ID, cmd.Args[0])
	if err != nil {
		log.Errorf("Error: %v", err)
		return "Error removing perm from user", nil
	}

	return "Permission removed", nil
}

func handleShowPerms(cmd *modulebase.ModuleCommand) (string, error) {
	if !canManageRoles(cmd.Guild, cmd.Message.Author.ID) {
		return "Insufficient permissions", nil
	}

	if len(cmd.Args) != 0 {
		return "This command takes no arguments", nil
	}

	permList, err := perms.PermsList()
	if err != nil {
		log.Errorf("Error getting perms %v", err)
		return "Error getting perms", nil
	}

	w := &tabwriter.Writer{}
	buf := &bytes.Buffer{}
	w.Init(buf, 0, 3, 0, ' ', 0)
	fmt.Fprint(w, "```\n")

	for _, perm := range permList {
		fmt.Fprintf(w, "%s\n", perm.Name)
	}

	fmt.Fprint(w, "```\n")
	w.Flush()
	return buf.String(), nil
}

func canManageRoles(guild *discordgo.Guild, userId string) bool {
	member := utils.GetMember(guild, userId)
	for _, role := range member.Roles {
		for _, gRole := range guild.Roles {
			if gRole.ID == role {
				if gRole.Permissions&0x10000000 > 0 {
					return true
				}
			}
		}
	}
	return false
}
