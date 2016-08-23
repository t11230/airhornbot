package modules

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/t11230/ramenbot/lib/modules/greeter"
	"github.com/t11230/ramenbot/lib/modules/modulebase"
)

var (
	moduleSetupFunctions = map[string]modulebase.ModuleSetupFunc{
		greeter.ConfigName: greeter.SetupFunc,
	}

	commandMap = map[string]modulebase.ModuleCommandCB{}

	eventList = *new([]interface{})
)

func LoadModules(configs []modulebase.ModuleConfig) error {
	for _, conf := range configs {
		if !conf.Enable {
			log.Debugf("Skipping module %v", conf.Name)
			continue
		}

		log.Debugf("Loading module %v", conf.Name)
		info, err := moduleSetupFunctions[conf.Name](&conf)
		if err != nil {
			log.Errorf("Error loading module %v: %v", conf.Name, err)
			continue
		}

		if info == nil {
			log.Errorf("Error loading module %v: Callback was nil", conf.Name)
			continue
		}

		eventList = append(eventList, *info.Events...)

		for _, l := range *info.Commands {
			commandMap[l.Command] = l.Callback
		}
	}

	log.Debug("Registered commands:")
	for c := range commandMap {
		log.Debugf("%v", c)
	}

	log.Debugf("Registered %v events", len(eventList))
	return nil
}

func HandleCommand(s *discordgo.Session, cmd *modulebase.ModuleCommand) {
	err := commandMap[cmd.Command](s, cmd)
	if err != nil {
		log.Errorf("Error processing command %v: %v", cmd.Command, err)
	}
}

func InitEvents(s *discordgo.Session) {
	log.Debug("Adding event handlers")
	for _, e := range eventList {
		s.AddHandler(e)
	}

	log.Debug("Event handlers added")
}

// Callbacks to add:
//
// Ready
// MessageCreate
// PresenceUpdate
// VoiceStateUpdate
// GuildCreate
// GuildMemberAdd
// GuildMemberRemove
// GuildRoleCreate
// GuildRoleUpdate
// GuildRoleDelete
// GuildUpdate
// GuildDelete
