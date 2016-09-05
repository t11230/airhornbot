package gambling

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/t11230/ramenbot/lib/bits"
	"github.com/t11230/ramenbot/lib/modules/modulebase"
	"github.com/t11230/ramenbot/lib/utils"
	"math/rand"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

const (
	RoundTime = 30
	MinAnte   = 5
)

func betRoll(cmd *modulebase.ModuleCommand) (string, error) {
	log.Debug("Placing Betroll")

	w := &tabwriter.Writer{}
	buf := &bytes.Buffer{}
	w.Init(buf, 0, 4, 0, ' ', 0)

	betroll_help := `**betroll usage:** betroll *ante* *dietype (optional)*
    This command initiates a bet on a dice roll. The second argument is the ante that all participants must pay into the pool.
    The third optional argument specifies a type of die for the roll.
    **Die Types**
    **d6 (default):** 6-sided die.
    **d4:** 4-sided die.
    **d8:** 8-sided die.
    **d10:** 10-sided die.
    **d12:** 12-sided die.
    **d20:** 20-sided die.
    **other:** random integer generator between 1 and input.`

	event := getActiveBetRoll(cmd.Guild.ID)
	if event != nil {
		fmt.Fprintf(w, "**ERROR:** BetRoll Event already in progress.")
		w.Flush()
		return buf.String(), nil
	}

	if (len(cmd.Args) < 1) || (len(cmd.Args) > 2) {
		fmt.Fprintf(w, betroll_help)
		w.Flush()
		return buf.String(), nil
	}

	//6-sided die is default
	maxnum := 6

	if len(cmd.Args) > 1 {
		var err error
		maxnum, err = strconv.Atoi(strings.Replace(cmd.Args[1], "d", "", 1))
		if err != nil {
			//user entered non-numerical number of die sides
			fmt.Fprintf(w, "**ERROR:** Non-numerical dice submitted.  Please don't be a smartass.")
			w.Flush()
			return buf.String(), nil
		}
	}

	if maxnum <= 0 {
		// User entered an invalid number of die sides
		fmt.Fprintf(w, "**ERROR:** Invalid dice submitted: Dice number less than or equal to zero.")
		w.Flush()
		return buf.String(), nil
	}

	ante, err := strconv.Atoi(cmd.Args[0])
	if err != nil {
		//user entered non-numerical ante
		fmt.Fprintf(w, "**ERROR:** Non-numerical ante submitted.  Please don't be a smartass.")
		w.Flush()
		return buf.String(), nil
	}

	go doBetRollRound(cmd, maxnum, ante, RoundTime)

	return buf.String(), nil
}

func bid(cmd *modulebase.ModuleCommand) (string, error) {
	user := cmd.Message.Author
	event := getActiveBetRoll(cmd.Guild.ID)
	w := &tabwriter.Writer{}
	buf := &bytes.Buffer{}
	var err error

	w.Init(buf, 0, 4, 0, ' ', 0)

	event_error_msg := "**ERROR:** No BetRoll Event currently in progress."
	bid_error_msg := "**ERROR:** Non-numerical bid submitted.  Please don't be a smartass."
	db_error_msg := "**ERROR:** Database error."
	success_msg := "*Bet Successfully Placed* :ok_hand:"
	ante_error_msg := "Not enough bits for ante :slight_frown:"
	bid_help := `**bid usage:** bid *number*
    This command bids on a bet roll. The second argument is the result that you are bidding on.`

	if event == nil {
		fmt.Fprintf(w, event_error_msg)
		w.Flush()
		return buf.String(), nil
	}
	if len(cmd.Args) != 1 {
		fmt.Fprintf(w, bid_help)
		w.Flush()
		return buf.String(), nil
	}
	for _, b := range event.Players {
		if b.UserID == user.ID {
			fmt.Fprintf(w, "You can only bid once!")
			w.Flush()
			return buf.String(), nil
		}
	}
	var me Player
	me.UserID = user.ID
	me.Bid, err = strconv.Atoi(cmd.Args[0])
	if err != nil {
		fmt.Fprintf(w, bid_error_msg)
		w.Flush()
		return buf.String(), nil
	}
	if bits.GetBits(cmd.Guild.ID, user.ID) < event.Ante {
		fmt.Fprintf(w, ante_error_msg)
		w.Flush()
		return buf.String(), nil
	}
	bits.RemoveBits(cmd.Session, cmd.Guild.ID, user.ID, event.Ante, "BetRoll bid")
	err = betRollAddPlayer(cmd.Guild.ID, &me)
	if err != nil {
		fmt.Fprintf(w, db_error_msg)
		w.Flush()
		return buf.String(), nil
	}
	fmt.Fprintf(w, success_msg)
	w.Flush()
	return buf.String(), nil
}

func printBetRollTime(time int, ante int) string {
	w := &tabwriter.Writer{}
	buf := &bytes.Buffer{}
	alert := ""
	w.Init(buf, 0, 4, 0, ' ', 0)
	if time == 30 {
		alert = "Dice Roll in **" + strconv.Itoa(time) + " seconds**.  Ante is **" + strconv.Itoa(ante) + " bits**. !!bid *result* to bid"
	} else if time != 0 {
		alert = "Dice Roll in **" + strconv.Itoa(time) + " seconds**."
	} else {
		alert = "**Dice Roll starting now!**"
	}
	fmt.Fprintf(w, alert)
	w.Flush()
	return buf.String()
}

func doBetRollRound(cmd *modulebase.ModuleCommand, maxnum int, ante int, roundtime int) {
	log.Info("Starting BetRoll Round")
	var winnerIDs []string
	win_result := "Winner(s):\n"
	payout_result := "Payout: "
	err := betRollOpen(cmd.Guild.ID)
	if err != nil {
		log.Error("Failed to open BetRoll")
		return
	}
	if ante < MinAnte {
		ante = MinAnte
		cmd.Session.ChannelMessageSend(cmd.Message.ChannelID, "**Ante Below Minimum (5):** Ante has been set to 5")
	}
	err = setBetRollAnte(cmd.Guild.ID, ante)
	payout := 0
	if err != nil {
		log.Error("Failed to set BetRoll Ante")
		// TODO: cleanup betroll entry
		betRollClose(cmd.Guild.ID)
		return
	}
	for roundtime > 0 {
		cmd.Session.ChannelMessageSend(cmd.Message.ChannelID, printBetRollTime(roundtime, ante))
		time.Sleep(10 * time.Second)
		roundtime -= 10
	}
	//send message that round is starting at roundtime == 0
	cmd.Session.ChannelMessageSend(cmd.Message.ChannelID, printBetRollTime(roundtime, ante))
	players := getBetRollPlayers(cmd.Guild.ID)
	pool := ante * len(players)
	r := rand.Intn(maxnum) + 1
	result := getDieString(maxnum, r)
	for _, player := range players {
		if player.Bid == r {
			winnerIDs = append(winnerIDs, player.UserID)
		}
	}
	if len(winnerIDs) > 0 {
		payout = pool / len(winnerIDs)
	}
	payout_result = payout_result + strconv.Itoa(payout) + " bits"
	for _, winner := range winnerIDs {
		bits.AddBits(cmd.Session, cmd.Guild.ID, winner, payout, "BetRoll win", true)
		win_result = win_result + utils.GetPreferredName(cmd.Guild, winner) + "\n"
	}
	err = betRollClose(cmd.Guild.ID)
	if err != nil {
		log.Error("Failed to close BetRoll")
		return
	}
	w := &tabwriter.Writer{}
	buf := &bytes.Buffer{}
	w.Init(buf, 0, 4, 0, ' ', 0)
	fmt.Fprintf(w, "```\n")
	fmt.Fprintf(w, result)
	fmt.Fprintf(w, "```\n")
	fmt.Fprintf(w, "```\n")
	fmt.Fprintf(w, win_result)
	fmt.Fprintf(w, payout_result)
	fmt.Fprintf(w, "```\n")
	w.Flush()
	cmd.Session.ChannelMessageSend(cmd.Message.ChannelID, buf.String())
}
