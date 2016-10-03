package store

import (
	"github.com/t11230/ramenbot/lib/modules/modulebase"
    "github.com/t11230/ramenbot/lib/bits"
    "fmt"
)

const (
	ConfigName = "store"
	helpString = "**!!store** : Displays what features you can buy with your bits.\n"
)

var commandTree = []modulebase.ModuleCommandTree{
	{
		RootCommand: "store",
		Function: store,
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

func store(cmd *modulebase.ModuleCommand) (string, error){
    storeHelpString := `**BIT STORE**

    `
    var bitstring string;
    user := cmd.Message.Author
    storeHelpString=storeHelpString+`**SILENCE SOUND EFFECTS**
    **usage:** !!s silence *duration*
    Prevents any sound clips from being played for *duration* minutes

    **COST:** `;
    if bits.GetBits(cmd.Guild.ID, user.ID) < 100 {
        bitstring = "~~100 Bits per minute~~\n\n"
    } else {
        bitstring = "100 Bits per minute\n**YOU CAN AFFORD:** "+fmt.Sprintf("%v",bits.GetBits(cmd.Guild.ID, user.ID)/100)+" Minutes\n\n"
    }
    storeHelpString=storeHelpString+bitstring

    storeHelpString=storeHelpString+`**SET CUSTOM WELCOME MESSAGE**
    **usage:** !!greet myvoice *action* *<collection>* *<sound>*
    sets your greeting sound to *<collection>* *<sound>* (see !!s help)

    **COST:** `;
    if bits.GetBits(cmd.Guild.ID, user.ID) < 150 {
        bitstring = "~~150 Bits~~\n\n"
    } else {
        bitstring = "150 Bits\n\n"
    }
    storeHelpString=storeHelpString+bitstring

    storeHelpString=storeHelpString+`**CHANGE NICKNAME**
    **usage:** !!nick *nickname*
    Changes user's nickname to *nickname*

    **COST:** `;
    if bits.GetBits(cmd.Guild.ID, user.ID) < 200 {
        bitstring = "~~200 Bits~~\n\n"
    } else {
        bitstring = "200 Bits\n\n"
    }
    storeHelpString=storeHelpString+bitstring

    storeHelpString=storeHelpString+`**UPLOAD CUSTOM SOUND EFFECT**
    **usage:** put !!s upload *collection* *soundname* in the comments of an audio file attachment
    Processes the attached soundfile and adds it to the soundboard as *soundname* in the collection *collection*

    **COST:** `;
    if bits.GetBits(cmd.Guild.ID, user.ID) < 300 {
        bitstring = "~~300 Bits~~\n\n"
    } else {
        bitstring = "300 Bits\n\n"
    }
    storeHelpString=storeHelpString+bitstring

    storeHelpString=storeHelpString+`**ADD CUSTOM TITLE**
    **usage:** !!role create *rolename* *color*
    Creates role with name *rolename* and color indicated by *color*. Then, gives user that role.
        **color names:** red, orange, yellow, green, blue, purple

    **COST:** `;
    if bits.GetBits(cmd.Guild.ID, user.ID) < 650 {
        bitstring = "~~650 Bits~~\n\n"
    } else {
        bitstring = "650 Bits\n\n"
    }
    storeHelpString=storeHelpString+bitstring

    return storeHelpString, nil
}
