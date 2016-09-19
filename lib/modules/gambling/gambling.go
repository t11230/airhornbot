package gambling

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/t11230/ramenbot/lib/modules/modulebase"
	"github.com/t11230/ramenbot/lib/perms"
	"github.com/t11230/ramenbot/lib/utils"
	"math/rand"
	"strconv"
	"strings"
	"text/tabwriter"
)

const (
	ConfigName = "gambling"
	helpString = "**!!$** : This module allows the user to use a number of gambling events and functions.\n"
	gamblingHelpString = `**$**
This module allows the user to use a number of gambling events and functions.

**usage:** !!$ *function* *args...*

**functions:**
    *roll:* This function allows the user to roll a die and prints the result.
	*betroll:* This function allows the user to start a betroll event.
	*bid:* This function allows the user to bid on an in-progress event.
	*bits:* This function allows the user to display the bit values of users in the server.
	*give:* This function allows the user to give their bits to another user.
		**permissions required:** bits-admin
	*award:* This function allows the user to create bits out of thin air and give them to another user.
		**permissions required:** bits-admin
	*take:* This function allows the user to take bits from another user.
		**permissions required:** bits-admin

For more info on using any of these functions, type **!!$ [function name] help**`

	rollHelpString = `**ROLL**
**usage:** !!$ roll *<dietype>*
	This command initiates a dice roll.
	The second optional argument specifies a type of die for the roll.
**Die Types**
	**d6 (default):** 6-sided die.
	**d4:** 4-sided die.
	**d8:** 8-sided die.
	**d10:** 10-sided die.
	**d12:** 12-sided die.
	**d20:** 20-sided die.
	**other:** random integer generator between 1 and input.`
)

var commandTree = []modulebase.ModuleCommandTree{
	{
		RootCommand: "$",
		SubKeys: modulebase.SK{
			"roll": modulebase.CN{
				Function: rollDice,
			},
			"betroll": modulebase.CN{
				Function: betRoll,
			},
			"bid": modulebase.CN{
				Function: bid,
			},
			"bits": modulebase.CN{
				Function: showBits,
			},
			"give": modulebase.CN{
				Function: giveBits,
			},
			"award": modulebase.CN{
				Function:    awardBits,
				Permissions: []perms.Perm{bitsAdminPerm},
			},
			"take": modulebase.CN{
				Function:    takeBits,
				Permissions: []perms.Perm{bitsAdminPerm},
			},
		},
		Function: handleRootCommand,
	},
}

var bitsAdminPerm = perms.Perm{"bits-admin"}

// Called to initialize this module
func SetupFunc(config *modulebase.ModuleConfig) (*modulebase.ModuleSetupInfo, error) {
	return &modulebase.ModuleSetupInfo{
		Events:   nil,
		Commands: &commandTree,
		DBStart:  handleDbStart,
		Help:     helpString,
	}, nil
}

func handleDbStart() error {
	err := perms.CreatePerm(bitsAdminPerm.Name)
	if err != nil {
		log.Errorf("Error creating perm: %v", err)
		return err
	}
	return nil
}

func handleRootCommand(cmd *modulebase.ModuleCommand) (string, error) {
	return gamblingHelpString, nil
}

func rollDice(cmd *modulebase.ModuleCommand) (string, error) {
	draw := false
	r := 0
	maxnum := 0
	w := &tabwriter.Writer{}
	buf := &bytes.Buffer{}
	w.Init(buf, 0, 4, 0, ' ', 0)
	var err error
	log.Debugf("Args was: %v", cmd.Args)
	if len(cmd.Args) > 0 {
		if (len(cmd.Args) > 1) || (cmd.Args[0] == "help") {
			return rollHelpString, nil
		}
		if strings.HasPrefix(cmd.Args[0], "d") {
			maxnum, err = strconv.Atoi(strings.Replace(cmd.Args[0], "d", "", 1))
			if err != nil {
				return "Invalid dietype", nil
			}
			if isValidDie(maxnum) {
				draw = true
			}
		} else {
			maxnum, err = strconv.Atoi(cmd.Args[0])
			if err != nil {
				return "Invalid number for dietype", nil
			}
		}
		r = rand.Intn(maxnum) + 1
	} else {
		maxnum = 6
		r = rand.Intn(6) + 1
		draw = true
	}
	result := ""
	if draw {
		if maxnum == 6 {
			result = drawD6(r)
		} else if (maxnum == 4) || (maxnum == 8) {
			result = drawD4_D8(r)
		} else if maxnum == 10 {
			result = drawD10(r)
		} else if maxnum == 12 {
			result = drawD12(r)
		} else if maxnum == 20 {
			result = drawD20(r)
		}
	} else {
		result = "The result is: " + strconv.Itoa(r)
	}
	fmt.Fprintf(w, "```\n")
	fmt.Fprintf(w, result)
	fmt.Fprintf(w, "```\n")
	w.Flush()
	return buf.String(), nil
}

func isValidDie(num int) bool {
	return utils.IntInSlice(num, []int{4, 6, 8, 10, 12, 20})
}

func drawD6(r int) string {
	C := "o "
	s := "---------\n| " + string(C[utils.BooltoInt(r <= 1)]) + "   " + string(C[utils.BooltoInt(r <= 3)]) + " |\n| " + string(C[utils.BooltoInt(r <= 5)])
	z := string(C[utils.BooltoInt(r <= 5)]) + " |\n| " + string(C[utils.BooltoInt(r <= 3)]) + "   " + string(C[utils.BooltoInt(r <= 1)]) + " |\n---------"
	return s + " " + string(C[utils.BooltoInt((r&1) == 0)]) + " " + z
}

func drawD4_D8(r int) string {
	return "      *\n     * *\n    *   *\n   *  " + strconv.Itoa(r) + "  *\n  *       *\n * * * * * *"
}

func drawD10(r int) string {
	numstring := strconv.Itoa(r)
	if r > 9 {
		return "        *\n       * *\n      *   *\n     * " + string(numstring[0]) + " " + string(numstring[1]) + " *\n      *   *\n        *"
	} else {
		return "        *\n       * *\n      *   *\n     *  " + numstring + "  *\n      *   *\n        *"
	}
}

func drawD12(r int) string {
	numstring := strconv.Itoa(r)
	if r > 9 {
		return "         *\n      *     *\n    *   " + string(numstring[0]) + " " + string(numstring[1]) + "   *\n     *       *\n      * * * *"
	} else {
		return "         *\n      *     *\n    *    " + numstring + "    *\n     *       *\n      * * * *"
	}
}

func drawD20(r int) string {
	numstring := strconv.Itoa(r)
	if r > 9 {
		return "      *\n     * *\n    *   *\n   * " + string(numstring[0]) + " " + string(numstring[1]) + " *\n  *       *\n * * * * * *"
	} else {
		return "      *\n     * *\n    *   *\n   *  " + numstring + "  *\n  *       *\n * * * * * *"
	}
}

func getDieString(maxnum int, r int) string {
	if isValidDie(maxnum) {
		if maxnum == 6 {
			return drawD6(r)
		} else if (maxnum == 4) || (maxnum == 8) {
			return drawD4_D8(r)
		} else if maxnum == 10 {
			return drawD10(r)
		} else if maxnum == 12 {
			return drawD12(r)
		} else if maxnum == 20 {
			return drawD20(r)
		}
	}

	return "The result is: " + strconv.Itoa(r)
}
