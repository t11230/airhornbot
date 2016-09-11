package modules

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/t11230/ramenbot/lib/modules/admin"
	"github.com/t11230/ramenbot/lib/modules/gambling"
	"github.com/t11230/ramenbot/lib/modules/greeter"
	"github.com/t11230/ramenbot/lib/modules/help"
	"github.com/t11230/ramenbot/lib/modules/rolemod"
	"github.com/t11230/ramenbot/lib/modules/modulebase"
	"github.com/t11230/ramenbot/lib/modules/soundboard"
	"github.com/t11230/ramenbot/lib/modules/voicebonus"
)

var (
	moduleSetupFunctions = map[string]modulebase.ModuleSetupFunc{
		admin.ConfigName:      admin.SetupFunc,
		gambling.ConfigName:   gambling.SetupFunc,
		greeter.ConfigName:    greeter.SetupFunc,
		rolemod.ConfigName:    rolemod.SetupFunc,
		soundboard.ConfigName: soundboard.SetupFunc,
		voicebonus.ConfigName: voicebonus.SetupFunc,
		help.ConfigName:       help.SetupFunc,
	}

	moduleHelpStrings = map[string]string{}

	commandMap = map[string]modulebase.CN{}

	eventList = *new([]interface{})

	dbStartFunctions = []modulebase.ModuleDBStartFunc{}
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

		if info.Events != nil {
			eventList = append(eventList, *info.Events...)
		}

		for _, l := range *info.Commands {
			commandMap[l.RootCommand] = linkModuleCommandTree(&l)
		}

		if info.DBStart != nil {
			dbStartFunctions = append(dbStartFunctions, info.DBStart)
		}

		if info.Help != "" {
			moduleHelpStrings[conf.Name] = info.Help
		}
		log.Debugf("Command trees: %v", commandMap)
	}
	modulebase.GetModuleHelp = getModuleHelpString
	log.Debug("Registered commands:")
	for c := range commandMap {
		log.Debugf("%v", c)
	}

	log.Debugf("Registered %v events", len(eventList))
	return nil
}

func getModuleHelpString() (map[string]string, error){
	return moduleHelpStrings, nil
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

	log.Debug("CHECKING LEN ARGS")
	// If there are no args, just call the root's function
	if len(cmd.Args) <= 0 {
		log.Debug("ROOT FUNCTION CALLED")
		message, err := node.Function(moduleCmd)
		if err == nil {
			cmd.Session.ChannelMessageSend(cmd.Message.ChannelID, message)
			return
		}
	}

	args := cmd.Args

	log.Debugf("Entering node %v", node)
	for {
		if len(args) == 0 {
			log.Debug("Longest prefix found")
			break
		}

		nextNode, ok := node.SubKeys[args[0]]

		if !ok {
			log.Debugf("Longest prefix found. %v args left, %v", len(args), args)
			break
		}
		args = args[1:]
		node = nextNode
		log.Debugf("Entering node %v, args is: %v", node, args)
	}

	var err error
	var message string

	if node.Function != nil {
		moduleCmd.Args = args
		message, err = node.Function(moduleCmd)
	} else {
		log.Error("Function was nil")
		err = errors.New("Nil Function")
	}

	if err == nil {
		cmd.Session.ChannelMessageSend(cmd.Message.ChannelID, message)
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

func DBHooks() error {
	for _, f := range dbStartFunctions {
		err := f()
		if err != nil {
			return err
		}
	}
	return nil
}
