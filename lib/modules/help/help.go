package help

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"github.com/t11230/ramenbot/lib/modules/modulebase"
)

const (
	ConfigName = "help"
)

var commandTree = []modulebase.ModuleCommandTree{
	{
		RootCommand: "help",
		Function:    handleHelp,
	},
}

var helpString string

// Called to initialize this module
func SetupFunc(config *modulebase.ModuleConfig) (*modulebase.ModuleSetupInfo, error) {
	return &modulebase.ModuleSetupInfo{
		Commands: &commandTree,
	}, nil
}

func handleHelp(cmd *modulebase.ModuleCommand) (string, error) {
	log.Debug("HELP FUNCTION CALLED")
	if modulebase.GetModuleHelp == nil {
		return "Our clever solution has failed", nil
	}
	helpStrings, _ := modulebase.GetModuleHelp()
	var buffer bytes.Buffer
	buffer.WriteString("**MODULES**\n\n")
	for _, help := range helpStrings {
		buffer.WriteString(help)
	}
	buffer.WriteString("\nFor more info on using any of these modules, type the module followed by **help**\n")
	helpString = buffer.String()
	return helpString, nil
}
