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
