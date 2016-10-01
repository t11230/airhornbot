package blackjack

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/t11230/ramenbot/lib/bits"
)

// var pendingRoundTimers
var turnTimers map[string]*TurnTimer = map[string]*TurnTimer{}

// Async functions
func waitPlayers(guildId string, channelId string, minBet int) {
	log.Debugf("waitBlackjackPlayers: %v", minBet)

	if isRoundPending(guildId) {
		log.Error("Round already exists")
		return
	}

	pendingRound := &PendingRound{
		MinimumBet: minBet,
		Countdown:  3,
		Players:    []Player{},
		ChannelID:  channelId,
		GuildID:    guildId,
	}

	err := createPendingRound(guildId, pendingRound)
	if err != nil {
		log.Error("Error creating round: %v", err)
	}
}

func startRound(s *discordgo.Session, guildId string) {
	log.Debug("startBlackjackRound")
	pendingRound, err := getPendingRound(guildId)
	if err != nil {
		log.Errorf("Error getting pending round: %v", err)
		return
	}
	clearPendingRound(guildId)

	newRound := newRound(s, pendingRound)

	createRunningRound(guildId, newRound)

	log.Debug("Round started")
	go newRound.Start()
}

func checkPlayerTurn(guildId string, channelId string, userId string) bool {
	log.Debug("checkPlayerTurn")

	if t, ok := turnTimers[guildId]; ok {
		return t.CheckUser(userId)
	}
	return false
}

func handlePlayerAction(guildId string, channelId string, action Action) bool {
	log.Debug("handlePlayerAction")

	if t, ok := turnTimers[guildId]; ok {
		return t.SendAction(action)
	}
	return false
}

// Helper functions
func isRoundPending(guildId string) bool {
	r, _ := getPendingRound(guildId)
	return r != nil
}

func isRoundRunning(guildId string) bool {
	r, _ := getRunningRound(guildId)
	return r != nil
}

func getPendingPlayers(guildId string) []Player {
	log.Debug("blackjackPlayers")

	pendingRound, err := getPendingRound(guildId)
	if err != nil {
		log.Errorf("Error getting round: %v", err)
		return nil
	}

	if pendingRound == nil {
		log.Error("Round does not exist")
		return nil
	}

	return pendingRound.Players
}

func addPlayer(s *discordgo.Session, guildId string, userId string, bet int) (betPlaced bool, err error) {
	log.Debugf("addPlayer: GuildId: %v, UserId: %v", guildId, userId)

	pendingRound, err := getPendingRound(guildId)
	if err != nil {
		log.Errorf("Error getting round: %v", err)
		return
	}
	if pendingRound == nil {
		err = errors.New("Round does not exist")
		log.Error(err)
		return
	}

	newPlayer := Player{
		UserID:     userId,
		InitialBet: pendingRound.MinimumBet,
	}

	if bet > 0 {
		newPlayer.InitialBet = bet
	}

	if bits.GetBits(guildId, userId) < newPlayer.InitialBet {
		return
	}

	pendingRound.Players = append(pendingRound.Players, newPlayer)

	err = updatePendingRound(guildId, pendingRound)
	if err != nil {
		return
	}

	bits.RemoveBits(s, guildId, userId, newPlayer.InitialBet, "Blackjack bet")

	return true, nil
}
