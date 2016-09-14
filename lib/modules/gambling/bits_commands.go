package gambling

import (
	"bytes"
	"fmt"
	"github.com/t11230/ramenbot/lib/bits"
	"github.com/t11230/ramenbot/lib/modules/modulebase"
	"github.com/t11230/ramenbot/lib/perms"
	"github.com/t11230/ramenbot/lib/utils"
	"strconv"
	"strings"
	"text/tabwriter"
)

const (
	showBitsHelpString = `**BITS**
This function allows the user to display the bit values of users in the server.

**usage:**!!$ bits *<me>*
	If no argument is passed to bits, the bit values of all users in the server are printed from highest to lowest.
	If *me* is passed as an argument, only the caller's bit value is displayed.`

	giveBitsHelpString = `**GIVE**
This function allows the user to give their bits to another user.

**usage:**!!$ give *amount* *user*
	Gives *user* *amount* bits.
	requires bits-admin permissions`

	awardBitsHelpString = `**AWARD**
This function allows the user to create bits out of thin air and give them to another user.

**usage:**!!$ award *amount* *user*
	Awards *user* *amount* bits.
	requires bits-admin permissions`

	takeBitsHelpString = `**TAKE**
This function allows the user to take bits from another user.

**usage:**!!$ take *amount* *user*
	Takes *amount* bits from *user*.
	requires bits-admin permissions`
)

func showBits(cmd *modulebase.ModuleCommand) (string, error) {
	w := &tabwriter.Writer{}
	buf := &bytes.Buffer{}
	w.Init(buf, 0, 3, 0, ' ', 0)
	fmt.Fprint(w, "```\n")

	if len(cmd.Args) == 0 {
		for _, b := range bits.GetBitsLeaderboard(cmd.Guild.ID, 10) {
			name := utils.GetPreferredName(cmd.Guild, b.UserID)
			if b.Value != nil && *b.Value != 0 {
				fmt.Fprintf(w, "%s: \t %d bits\n", name, *b.Value)
			}
		}
		fmt.Fprint(w, "```\n")
		w.Flush()
		return buf.String(), nil
	}

	var userId string
	if cmd.Args[0] == "me" {
		userId = cmd.Message.Author.ID
	} else if cmd.Args[0] == "help" {
		return showBitsHelpString, nil
	} else {
		return showBitsHelpString, nil
	}

	b := bits.GetBits(cmd.Guild.ID, userId)

	name := utils.GetPreferredName(cmd.Guild, userId)
	fmt.Fprintf(w, "%s: \t %d bits\n", name, b)

	fmt.Fprint(w, "```\n")
	w.Flush()
	return buf.String(), nil
}

func giveBits(cmd *modulebase.ModuleCommand) (string, error) {
	if len(cmd.Args) < 2 {
		return giveBitsHelpString, nil
	}

	amount, err := strconv.Atoi(cmd.Args[0])
	if err != nil {
		return "Invalid amount", nil
	}

	if amount <= 0 {
		return "Amount must be a positive nonzero integer", nil
	}

	userName := strings.Join(cmd.Args[1:], " ")
	user, err := utils.FindUser(cmd.Guild, userName)
	if err != nil {
		return "Unable to find user", nil
	}

	gifter := cmd.Message.Author
	if bits.GetBits(cmd.Guild.ID, gifter.ID) < amount {
		return "Not enough bits", nil
	}

	bits.RemoveBits(cmd.Session, cmd.Guild.ID, gifter.ID, amount, "Gift")
	bits.AddBits(cmd.Session, cmd.Guild.ID, user.ID, amount, "Gift", false)

	message := fmt.Sprintf("Transferred %d bits from %v to %v",
		amount,
		utils.GetPreferredName(cmd.Guild, gifter.ID),
		utils.GetPreferredName(cmd.Guild, user.ID))

	return message, nil
}

func awardBits(cmd *modulebase.ModuleCommand) (string, error) {
	permsHandle := perms.GetPermsHandle(cmd.Guild.ID, ConfigName)
	if !permsHandle.CheckPerm(cmd.Message.Author.ID, "bits-admin") {
		return "Insufficient permissions", nil
	}

	if len(cmd.Args) < 2 {
		return awardBitsHelpString, nil
	}

	amount, err := strconv.Atoi(cmd.Args[0])
	if err != nil {
		return "Invalid amount", nil
	}

	if amount <= 0 {
		return "Amount must be a positive nonzero integer", nil
	}

	userName := strings.Join(cmd.Args[1:], " ")
	user, err := utils.FindUser(cmd.Guild, userName)
	if err != nil {
		return "Unable to find user", nil
	}

	bits.AddBits(cmd.Session, cmd.Guild.ID, user.ID, amount, "Awarded bits", false)

	message := fmt.Sprintf("Awarded %d bits to %v",
		amount,
		utils.GetPreferredName(cmd.Guild, user.ID))

	return message, nil
}

func takeBits(cmd *modulebase.ModuleCommand) (string, error) {
	permsHandle := perms.GetPermsHandle(cmd.Guild.ID, ConfigName)
	if !permsHandle.CheckPerm(cmd.Message.Author.ID, "bits-admin") {
		return "Insufficient permissions", nil
	}

	if len(cmd.Args) < 2 {
		return takeBitsHelpString, nil
	}

	amount, err := strconv.Atoi(cmd.Args[0])
	if err != nil {
		return "Invalid amount", nil
	}

	if amount <= 0 {
		return "Amount must be a positive nonzero integer", nil
	}

	userName := strings.Join(cmd.Args[1:], " ")
	user, err := utils.FindUser(cmd.Guild, userName)
	if err != nil {
		return "Unable to find user", nil
	}

	bits.RemoveBits(cmd.Session, cmd.Guild.ID, user.ID, amount, "Took bits")

	message := fmt.Sprintf("Took %d bits from %v",
		amount,
		utils.GetPreferredName(cmd.Guild, user.ID))

	return message, nil
}
