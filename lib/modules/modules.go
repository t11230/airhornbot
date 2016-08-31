package modules

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/t11230/ramenbot/lib/modules/greeter"
	"github.com/t11230/ramenbot/lib/modules/modulebase"
	"github.com/t11230/ramenbot/lib/modules/soundboard"
)

var (
	moduleSetupFunctions = map[string]modulebase.ModuleSetupFunc{
		greeter.ConfigName:    greeter.SetupFunc,
		soundboard.ConfigName: soundboard.SetupFunc,
	}

	commandMap = map[string]modulebase.CN{}

	eventList = *new([]interface{})
)

type Command struct {
	Session *discordgo.Session
	Guild   *discordgo.Guild
	Message *discordgo.Message
	Args    []string
	Command string
}

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
			commandMap[l.RootCommand] = linkModuleCommandTree(&l)
		}

		log.Debugf("Command trees: %v", commandMap)
	}

	log.Debug("Registered commands:")
	for c := range commandMap {
		log.Debugf("%v", c)
	}

	log.Debugf("Registered %v events", len(eventList))
	return nil
}

func linkModuleCommandTree(tree *modulebase.ModuleCommandTree) modulebase.CN {
	root := modulebase.CN{
		Parent:        nil,
		SubKeys:       modulebase.SK{},
		Function:      tree.Function,
		ErrorFunction: tree.ErrorFunction,
	}

	for c := range tree.SubKeys {
		node := tree.SubKeys[c]
		node.Parent = &root
		linkModuleCommandNode(&node)
		root.SubKeys[c] = node
	}

	return root
}

func linkModuleCommandNode(parent *modulebase.CN) {
	for c := range parent.SubKeys {
		node := parent.SubKeys[c]
		node.Parent = parent
		linkModuleCommandNode(&node)
	}
}

func HandleCommand(cmd *Command) {
	// Try to do a longest prefix match

	node, ok := commandMap[cmd.Command]
	if !ok {
		log.Debugf("Invalid command %v", cmd.Command)
		return
		// log.Errorf("Error processing command %v: %v", cmd.Command, err)
	}

	moduleCmd := &modulebase.ModuleCommand{
		Session: cmd.Session,
		Guild:   cmd.Guild,
		Message: cmd.Message,
	}

	// If there are no args, just call the root's function
	if len(cmd.Args) <= 0 {
		node.Function(moduleCmd)
		return
	}

	args := cmd.Args

	log.Debugf("Entering node %v", node)
	for {
		if len(args) == 0 {
			log.Debug("Longest prefix found")
			break
		}

		nextNode, ok := node.SubKeys[args[0]]
		args = cmd.Args[1:]

		if !ok {
			log.Debugf("Longest prefix found. %v args left", len(args))
			break
		}
		node = nextNode
		log.Debugf("Entering node %v, args is: %v", node, args)
	}

	var err error

	if node.Function != nil {
		moduleCmd.Args = args
		err = node.Function(moduleCmd)
	} else {
		log.Error("Function was nil")
		err = errors.New("Nil Function")
	}

	if err == nil {
		return
	}

	log.Errorf("Error in function: %v", err)
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
