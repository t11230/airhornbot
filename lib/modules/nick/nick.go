package nick

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/t11230/ramenbot/lib/modules/modulebase"
    "github.com/t11230/ramenbot/lib/bits"
	"github.com/t11230/ramenbot/lib/ramendb"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
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
	nickCollName = "nicktrack"

)


var commandTree = []modulebase.ModuleCommandTree{
	{
		RootCommand: "nick",
		Function: handleNickChange,
	},
}

type nickCollection struct {
	*mgo.Collection
}

type nicknameConfig struct {
	UserID         string
	Nickname  string `bson:",omitempty"`
}

// Called to initialize this module
func SetupFunc(config *modulebase.ModuleConfig) (*modulebase.ModuleSetupInfo, error) {
	events := []interface{}{
		nickChangeUpdateCallback,
	}
	return &modulebase.ModuleSetupInfo{
		Events:   &events,
		Commands: &commandTree,
		Help:     helpString,
	}, nil
}

func getNickName(msg string) string {
    msgArr := strings.Split(msg, " ")
    return strings.Join(msgArr[1:], " ")
}

func handleNickChange(cmd *modulebase.ModuleCommand) (string, error) {
	nicknames := nickCollection{ramendb.GetCollection(cmd.Guild.ID, nickCollName)}
    user := cmd.Message.Author
    guild := cmd.Guild
    s := cmd.Session
    if bits.GetBits(guild.ID, user.ID) < 200 {
		return "**FAILED TO CHANGE NICKNAME:** Insufficient bits.", nil
	}
    if (len(cmd.Args)<1)||(cmd.Args[0]=="help") {
        return nickHelpString, nil
    }
    nickname := getNickName(cmd.Message.Content)
	upsertdata := bson.M{"$set": nicknameConfig{
        UserID:    user.ID,
        Nickname: nickname,
    }}
	nicknames.Upsert(nicknameConfig{UserID: user.ID}, upsertdata)
    err := s.GuildMemberNickname(guild.ID, user.ID, nickname)
    if err != nil {
        log.Errorf("Failed to update user's nickname: %v", err)
        return "**Failed to update user's nickname**", nil
    }
    bits.RemoveBits(s, guild.ID, user.ID, 200, "Changed nickname to "+nickname)
    return "", nil
}

func nickChangeUpdateCallback(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {
	nicknames := nickCollection{ramendb.GetCollection(m.GuildID, nickCollName)}

	guild, _ := s.State.Guild(m.GuildID)
	if guild == nil {
		log.WithFields(log.Fields{
			"guild": m.GuildID,
		}).Warning("Failed to grab guild")
		return
	}

	member := m.Member
	user := member.User
	if member == nil {
		log.WithFields(log.Fields{
			"member": member,
		}).Warning("Failed to grab member")
		return
	}

	if member.User.Bot {
		return
	}

	data := nicknameConfig{}
	nicknames.Find(nicknameConfig{UserID: user.ID}).One(&data)
	trueNick := data.Nickname
	if trueNick == "" {
		log.Debug("No Nickname DB Info")
		return
	}
	if member.Nick != trueNick {
		err := s.GuildMemberNickname(guild.ID, user.ID, trueNick)
	    if err != nil {
	        log.Errorf("Failed to correct user's nickname: %v", err)
	        return
	    }
		return
	}
}
