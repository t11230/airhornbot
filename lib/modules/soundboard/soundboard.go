package soundboard

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/t11230/ramenbot/lib/modules/modulebase"
	"github.com/t11230/ramenbot/lib/sound"
	"github.com/t11230/ramenbot/lib/utils"
)

// Module name used in the config file
const (
	ConfigName = "soundboard"
	helpString = "**!!s** : This module allows the user to play sounds from a dank soundboard.\n"
	sHelpString = `**S**
This module allows the user to play sounds from a dank soundboard.

**usage:** !!s *collection* *sound*
	Plays the sound *sound* from the collection of sounds *collection*
	Collections and the sounds they contain listed below:

**airhorn:** default, reverb, spam, tripletap, fourtap, distant, echo, clownfull, clownshort, clownspam, highfartlong, highfartshort, midshort, truck
**anotha:** one, one_echo, one_classic, dialup
**johncena:** airhorn, echo, full, jc, nameis, spam, collect
**ethan:** areyou_classic, areyou_condensed, areyou_crazy, areyou_ethan, classic, echo, high, slowandlow, cuts, beat, sodiepop, vape
**stan:** herd, moo, x3
**trump:** 10ft, wall, mess, bing, getitout, tractor, worstpres, china, mexico, special
**music:** serbian, techno
**meme:** headshot, wombo, triple, camera, gandalf, mad, ateam, bennyhill, tuba, donethis, leeroy, slam, nerd, kappa, digitalsports, csi, nogod, welcomebdc
**birthday:** horn, horn3, sadhorn, weakhorn
**owult:** dva_enemy, genji_enemy, genji_friendly, hanzo_enemy, hanzo_friendly, junkrat_enemy, junkrat_friendly, lucio_friendly, lucio_enemy, mccree_enemy, mccree_friendly, mei_friendly, mei_enemy, pharah_enemy, reaper_friendly, 76_enemy, symmetra_friendly, torbjorn_enemy, tracer_enemy, tracer_friendly, widow_enemy, widow_friendly, zarya_enemy, zarya_friendly, zenyatta_enemy, dva_;), anyong
**dota:** waow, balance, rekt, stick, mana, disaster, liquid, history, smut, team, aegis
**overwatch:** payload, whoa, woah, winky, turd, ryuugawagatekiwokurau, cyka, noon, somewhere, lift, russia
**wc3:** work, awake
**sp:** screw, authority
**sv:** piss, fucks, shittalk, attractive, win
**archer:** dangerzone, klog

**EXAMPLE:** !!s airhorn default

`
)

// List of commands that this module accepts
var commandTree = []modulebase.ModuleCommandTree{
	{
		RootCommand: "s",
		SubKeys:     modulebase.SK{},
		Function:    handleSoundCommand,
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

func handleSoundCommand(cmd *modulebase.ModuleCommand) (string, error) {
	log.Debugf("Sound :%v", cmd.Args)
	if len(cmd.Args) == 0 || cmd.Args[0] =="help" {
		availableCollections()
		return sHelpString, nil
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
					return "", errors.New("Sound was nil")
				}
			}

			go sound.EnqueuePlay(cmd.Session, cmd.Message.Author, cmd.Guild, coll, snd)
			return "", nil
		}
	}

	return "Unable to find sound", nil
}

func availableCollections() []string {
	colls := []string{}
	for _, c := range sound.GetCollections() {
		colls = append(colls, c.Commands[0])
	}
	return colls
}
