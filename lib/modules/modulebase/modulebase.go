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
type ModuleDBStartFunc func() error

type ModuleSetupInfo struct {
	Events   *[]interface{}
	Commands *[]ModuleCommandTree
	DBStart  ModuleDBStartFunc
}

type ModuleCommand struct {
	Session *discordgo.Session
	Guild   *discordgo.Guild
	Message *discordgo.Message
	Args    []string
}

type ModuleCommandFunc func(*ModuleCommand) (string, error)
type ModuleCommandErrorFunc func(*ModuleCommand, error)

type ModuleCommandTree struct {
	RootCommand   string
	SubKeys       SK
	Function      ModuleCommandFunc
	ErrorFunction ModuleCommandErrorFunc
}

type SK map[string]CN
type CN struct {
	Parent        *CN
	SubKeys       SK
	Function      ModuleCommandFunc
	ErrorFunction ModuleCommandErrorFunc
}
