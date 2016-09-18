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

const (
	ConfigName = "admin"
	helpString = "**!!adm** : This module allows the user to control permissions for using other modules.\n"

	admHelpString = `**ADM**
This module allows the user to control permissions for using other modules.

**usage:** !!adm *function* *args...*

**discord permissions required:** Manage Roles

**functions:**
	*showperms:* This function displays the names of all possible permissions.
    *addperm:* This function allows the user to add a set of permissions to a user.
	*delperm:* This function allows the user to remove a set of permissions from a user.

For more info on using any of these functions, type **!!adm [function name] help**`

	showHelpString = `**SHOWPERMS**

**usage:** !!adm showperms
	Displays the names of all possible permissions.`

	addHelpString = `**ADDPERM**

**usage:** !!adm addperm *permissions* *username*
	Adds the permissions specified by *permisions* to the user specified by *username*`

	delHelpString = `**DELPERM**

**usage:** !!adm delperm *permissions* *username*
	Removes the permissions specified by *permisions* from the user specified by *username*`
)

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
		Function: handleAdmHelp,
	},
}

// Called to initialize this module
func SetupFunc(config *modulebase.ModuleConfig) (*modulebase.ModuleSetupInfo, error) {
	return &modulebase.ModuleSetupInfo{
		Events:   nil,
		Commands: &commandTree,
		Help:     helpString,
	}, nil
}

func handleAdmHelp(cmd *modulebase.ModuleCommand) (string, error) {
	return admHelpString, nil
}

func handleAddPerm(cmd *modulebase.ModuleCommand) (string, error) {
	log.Debug("Handling addperm command")

	if !canManageRoles(cmd.Guild, cmd.Message.Author.ID) {
		return "Insufficient permissions", nil
	}

	if len(cmd.Args) < 2 {
		return addHelpString, nil
	}

	h := perms.GetPermsHandle(cmd.Guild.ID)
	permExist, err := perms.PermExists(cmd.Args[0])
	if err != nil {
		return "Error checking perm", nil
	}
	if permExist == nil {
		return "Invalid perm", nil
	}

	userName := strings.Join(cmd.Args[1:], " ")
	user, err := utils.FindUser(cmd.Guild, userName)
	if err != nil {
		return "Unable to find user", nil
	}

	if h.CheckPerm(user.ID, permExist) {
		return "User already has that perm", nil
	}

	err = h.AddPerm(user.ID, permExist)
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

	if len(cmd.Args) < 2 {
		return delHelpString, nil
	}

	h := perms.GetPermsHandle(cmd.Guild.ID)
	permExist, err := perms.PermExists(cmd.Args[0])
	if err != nil {
		return "Error checking perm", nil
	}
	if permExist == nil {
		return "Invalid perm", nil
	}

	userName := strings.Join(cmd.Args[1:], " ")
	user, err := utils.FindUser(cmd.Guild, userName)
	if err != nil {
		return "Unable to find user", nil
	}

	if !h.CheckPerm(user.ID, permExist) {
		return "User does not have that permission", nil
	}

	err = h.RemovePerm(user.ID, permExist)
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
		return showHelpString, nil
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
