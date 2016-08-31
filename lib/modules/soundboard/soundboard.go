package soundboard

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/t11230/ramenbot/lib/modules/modulebase"
	"github.com/t11230/ramenbot/lib/sound"
	"github.com/t11230/ramenbot/lib/utils"
)

// Module name used in the config file
const ConfigName = "soundboard"

// List of commands that this module accepts
var commandTree = []modulebase.ModuleCommandTree{
	{
		RootCommand: "sound",
		SubKeys:     modulebase.SK{},
		Function:    handleSoundCommand,
	},
}

// Called to initialize this module
func SetupFunc(config *modulebase.ModuleConfig) (*modulebase.ModuleSetupInfo, error) {
	return &modulebase.ModuleSetupInfo{
		Events:   nil,
		Commands: &commandTree,
	}, nil
}

func handleSoundCommand(cmd *modulebase.ModuleCommand) error {
	log.Debugf("Sound :%v", cmd.Args)
	if len(cmd.Args) == 0 {
		availableCollections()
		return nil
	}

	for _, coll := range sound.GetCollections() {
		if utils.Scontains(cmd.Args[0], coll.Commands...) {

			// If they passed a specific sound effect, find and select that (otherwise play nothing)
			var snd *sound.Sound
			if len(cmd.Args) > 1 {
				for _, s := range coll.Sounds {
					if cmd.Args[1] == s.Name {
						snd = s
					}
				}

				if snd == nil {
					return errors.New("Sound was nil")
				}
			}

			go sound.EnqueuePlay(cmd.Session, cmd.Message.Author, cmd.Guild, coll, snd)
			return nil
		}
	}

	return nil
}

func availableCollections() []string {
	colls := []string{}
	for _, c := range sound.GetCollections() {
		colls = append(colls, c.Commands[0])
	}
	return colls
}

// import (
//     "github.com/bwmarrin/discordgo"
// )

// type TextFunction func(*discordgo.Guild, *discordgo.Message, []string) string

// type TextCollection struct {
//     Commands    []string
//     Function    TextFunction
// }

// var GITHUB *TextCollection = &TextCollection{
//     Commands: []string{
//         "github",
//         "git",
//     },
//     Function: func (guild *discordgo.Guild,
//                     user *discordgo.Message,
//                     args []string) string {
//         return "https://github.com/t11230/airhornbot"
//     },
// }

// var SOUNDCOMMANDS *TextCollection = &TextCollection{
//     Commands: []string{
//         "imanoob",
//         "l2p",
//     },
//     Function: func (*discordgo.Guild, *discordgo.Message, []string) string {
//         return sndGetSoundCommands()
//     },
// }

// var HILLARY *TextCollection = &TextCollection{
//     Commands: []string{
//         "hillary",
//     },
//     Function: func (*discordgo.Guild, *discordgo.Message, []string) string {
//         return "https://i.imgur.com/1PFAZsV.jpg"
//     },
// }

// var MARKOV *TextCollection = &TextCollection{
//     Commands: []string{
//         "text",
//     },
//     Function: func (guild *discordgo.Guild,
//                     message *discordgo.Message,
//                     args []string) string {
//         return mkGetMessage(guild, message.Author)
//     },
// }

// var STATS *TextCollection = &TextCollection{
//     Commands: []string{
//         "stats",
//     },
//     Function: gpHandleStatsCommand,
// }

// var BITS *TextCollection = &TextCollection{
//     Commands: []string{
//         "bits",
//     },
//     Function: bitsPrintStats,
// }

// var DICEROLL *TextCollection = &TextCollection{
//     Commands: []string{
//         "roll",
//     },
//     Function: rollDice,
// }

// var BETROLL *TextCollection = &TextCollection{
//     Commands: []string{
//         "betroll",
//     },
//     Function: betRoll,
// }

// var BID *TextCollection = &TextCollection{
//     Commands: []string{
//         "bid",
//     },
//     Function: bid,
// }

// var TEXTCMDS []*TextCollection = []*TextCollection{
//     SOUNDCOMMANDS,
//     GITHUB,
//     HILLARY,
//     MARKOV,
//     STATS,
//     BITS,
//     DICEROLL,
//     BETROLL,
//     BID,
// }
