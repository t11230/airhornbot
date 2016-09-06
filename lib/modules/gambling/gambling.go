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

const ConfigName = "gambling"

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
				Function: awardBits,
			},
			"take": modulebase.CN{
				Function: takeBits,
			},
		},
		Function: handleRootCommand,
	},
}

// Called to initialize this module
func SetupFunc(config *modulebase.ModuleConfig) (*modulebase.ModuleSetupInfo, error) {
	return &modulebase.ModuleSetupInfo{
		Events:   nil,
		Commands: &commandTree,
		DBStart:  handleDbStart,
	}, nil
}

func handleDbStart() error {
	perms.CreatePerm("bits-admin")
	return nil
}

func handleRootCommand(cmd *modulebase.ModuleCommand) (string, error) {
	return "", nil
}

func rollDice(cmd *modulebase.ModuleCommand) (string, error) {
	roll_help := `**roll usage:** roll *dietype (optional)*
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
			return roll_help, nil
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
