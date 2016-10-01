package blackjack

import (
	"github.com/t11230/ramenbot/lib/modules/modulebase"
	"strconv"
)

// Command handlers
func HandleStart(cmd *modulebase.ModuleCommand) (response string, err error) {
	if !validateBlackjackStart(cmd) {
		return "Blackjack start help", nil
	}

	minBet, err := strconv.Atoi(cmd.Args[0])
	if err != nil {
		return "Invalid minimum bet", nil
	}

	go waitPlayers(cmd.Guild.ID, cmd.Message.ChannelID, minBet)

	return "Waiting for players", nil
}

func HandleDeal(cmd *modulebase.ModuleCommand) (string, error) {
	if len(cmd.Args) > 0 {
		return "Blackjack deal takes no arguments", nil
	}

	if !isRoundPending(cmd.Guild.ID) {
		return "A round has not been created", nil
	}

	if len(getPendingPlayers(cmd.Guild.ID)) == 0 {
		return "No players", nil
	}

	go startRound(cmd.Session, cmd.Guild.ID)

	return "Round started", nil
}

func HandleBet(cmd *modulebase.ModuleCommand) (string, error) {
	if len(cmd.Args) > 1 {
		return "Blackjack bet takes one or no arguments", nil
	}

	if !isRoundPending(cmd.Guild.ID) {
		return "A round has not been created", nil
	}

	if len(getPendingPlayers(cmd.Guild.ID)) == 4 {
		return "Max players", nil
	}

	bet := -1
	if len(cmd.Args) > 0 {
		var err error
		bet, err = strconv.Atoi(cmd.Args[0])
		if err != nil {
			return "Invalid bet amount. something...something smartass", nil
		}
	}

	betPlaced, err := addPlayer(cmd.Session, cmd.Guild.ID, cmd.Message.Author.ID, bet)
	if err != nil {
		return "Error adding player", nil
	}

	if !betPlaced {
		return "Insufficient bits", nil
	}

	return "Bet placed", nil
}

func HandleHit(cmd *modulebase.ModuleCommand) (string, error) {
	if len(cmd.Args) > 0 {
		return "Blackjack hit takes no arguments", nil
	}

	if !checkPlayerTurn(cmd.Guild.ID, cmd.Message.ChannelID, cmd.Message.Author.ID) {
		return "It is not your turn", nil
	}

	if !handlePlayerAction(cmd.Guild.ID, cmd.Message.ChannelID, ActionHit) {
		return "Invalid Action", nil
	}

	return "", nil
}

func HandleStay(cmd *modulebase.ModuleCommand) (string, error) {
	if len(cmd.Args) > 0 {
		return "Blackjack stay takes no arguments", nil
	}

	if !checkPlayerTurn(cmd.Guild.ID, cmd.Message.ChannelID, cmd.Message.Author.ID) {
		return "It is not your turn", nil
	}

	if !handlePlayerAction(cmd.Guild.ID, cmd.Message.ChannelID, ActionStay) {
		return "Invalid Action", nil
	}

	return "", nil
}

func HandleDoubleDown(cmd *modulebase.ModuleCommand) (string, error) {
	return "Not yet implemented", nil
}

func HandleSplit(cmd *modulebase.ModuleCommand) (string, error) {
	return "Not yet implemented", nil
}

// Validation functions
func validateBlackjackStart(cmd *modulebase.ModuleCommand) bool {
	if len(cmd.Args) != 1 {
		return false
	}

	if cmd.Args[0] == "help" {
		return false
	}

	return true
}
