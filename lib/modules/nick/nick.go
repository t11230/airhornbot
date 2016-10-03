package nick

import (
	log "github.com/Sirupsen/logrus"
	"github.com/t11230/ramenbot/lib/modules/modulebase"
    "github.com/t11230/ramenbot/lib/bits"
	"strings"
)

const (
	ConfigName = "nick"
	helpString = "**!!nick** : This module allows the user to change their nickname.\n"

	nickHelpString = `**NICK**
This module allows the user to change their nickname.

**usage:** !!nick *nickname*
    Changes user's nickname to *nickname*
    **WARNING** Changing your nickname costs **200 bits**
`
)

var commandTree = []modulebase.ModuleCommandTree{
	{
		RootCommand: "nick",
		Function: handleNickChange,
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

func getNickName(msg string) string {
    msgArr := strings.Split(msg, " ")
    return strings.Join(msgArr[1:], " ")
}

func handleNickChange(cmd *modulebase.ModuleCommand) (string, error) {
    user := cmd.Message.Author
    guild := cmd.Guild
    s := cmd.Session
    if bits.GetBits(guild.ID, user.ID) < 200 {
		return "**FAILED TO ADD ROLE:** Insufficient bits.", nil
	}
    if (len(cmd.Args)<1)||(cmd.Args[0]=="help") {
        return nickHelpString, nil
    }
    nickname := getNickName(cmd.Message.Content)
    err := s.GuildMemberNickname(guild.ID, user.ID, nickname)
    if err != nil {
        log.Errorf("Failed to update user's nickname: %v", err)
        return "**Failed to update user's nickname**", nil
    }
    bits.RemoveBits(s, guild.ID, user.ID, 200, "Changed nickname to "+nickname)
    return "", nil
}
