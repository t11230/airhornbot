package modulebase

import (
	"github.com/bwmarrin/discordgo"
)

type ModuleConfig struct {
	Name    string
	Enable  bool
	Options interface{}
}

type ModuleSetupFunc func(*ModuleConfig) (*ModuleSetupInfo, error)

type ModuleSetupInfo struct {
	Events   *[]interface{}
	Commands *[]ModuleCommandListener
}

type ModuleCommandCB func(*discordgo.Session, *ModuleCommand) error

type ModuleCommandListener struct {
	Command  string
	Callback ModuleCommandCB
}

type ModuleCommand struct {
	Guild   *discordgo.Guild
	Message *discordgo.Message
	Command string
	Args    []string
}
