package modulebase

import (
	"github.com/bwmarrin/discordgo"
	"github.com/t11230/ramenbot/lib/perms"
)

type ModuleConfig struct {
	Name    string
	Enable  bool
	Options interface{}
}
type ModuleHelpFunc func() (map[string]string, error)
type ModuleSetupFunc func(*ModuleConfig) (*ModuleSetupInfo, error)
type ModuleDBStartFunc func() error

type ModuleSetupInfo struct {
	Events   *[]interface{}
	Commands *[]ModuleCommandTree
	Help     string
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
	Permissions   []perms.Perm
}

type SK map[string]CN
type CN struct {
	Parent        *CN
	SubKeys       SK
	Function      ModuleCommandFunc
	ErrorFunction ModuleCommandErrorFunc
	Permissions   []perms.Perm
}

var GetModuleHelp ModuleHelpFunc
