package blackjack

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/looplab/fsm"
	"github.com/t11230/ramenbot/lib/bits"
	"github.com/t11230/ramenbot/lib/modules/gambling/cards"
	"github.com/t11230/ramenbot/lib/utils"
	"image/png"
	"time"
)

func newRound(s *discordgo.Session, pending *PendingRound) (r *Round) {
	// Copy the nedded parameters from the pending game
	r = &Round{
		Session:    s,
		GuildID:    pending.GuildID,
		ChannelID:  pending.ChannelID,
		MinimumBet: pending.MinimumBet,
		Players:    pending.Players,
	}

	// Create a new shuffled deck
	r.Deck = cards.NewDeck(true)

	// Setup the dealer as the bot user
	user, err := s.User("@me")
	if err != nil {
		log.Errorf("Error getting bot user: %v", err)
		return nil
	}
	r.Dealer = Player{
		UserID:     user.ID,
		InitialBet: r.MinimumBet,
	}

	// Create ths FSM to track the game
	r.FSM = fsm.NewFSM("setup", fsm.Events{
		{Name: "deal", Src: []string{"setup"}, Dst: "dealt"},
		{Name: "showHands", Src: []string{"dealt"}, Dst: "handsShown"},
		{Name: "startPlayers", Src: []string{"handsShown"}, Dst: "playerGo"},
		{Name: "dealerBlackjack", Src: []string{"handsShown"}, Dst: "closeGame"},
		{Name: "startDealer", Src: []string{"playerGo"}, Dst: "dealerGo"},
		{Name: "dealerDone", Src: []string{"dealerGo"}, Dst: "closeGame"},
	}, fsm.Callbacks{
		"enter_dealt": func(e *fsm.Event) {
			go r.Deal()
		},
		"enter_handsShown": func(e *fsm.Event) {
			go func() {
				r.ShowTable()
				r.FSM.Event("startPlayers")
			}()
		},
		"enter_playerGo": func(e *fsm.Event) {
			go r.PlayerTurn()
		},
		"enter_dealerGo": func(e *fsm.Event) {
			go r.DealerTurn()
		},
		"enter_closeGame": func(e *fsm.Event) {
			go r.CloseGame()
		},
	})

	return
}

func (r *Round) Start() {
	log.Debug("Starting round")
	r.FSM.Event("deal")
}

func (r *Round) Deal() {
	log.Debug("Dealing round")

	// Draw all the cards we will use for now
	drawResult := r.Deck.Draw((len(r.Players) + 1) * 2)

	// Create the player's hands
	for i, _ := range r.Players {
		r.Players[i].AddHand(drawResult.Cards[0])
		drawResult.Cards = drawResult.Cards[1:]
	}

	// Now deal a card to the dealer
	drawResult.Cards[0].IsFaceDown = true
	r.Dealer.AddHand(drawResult.Cards[0])
	drawResult.Cards = drawResult.Cards[1:]

	// Deal the second card to each player
	for i, _ := range r.Players {
		r.Players[i].Hands[0].Pile.AddCards(drawResult.Cards[0])
		drawResult.Cards = drawResult.Cards[1:]
	}

	// Deal the second card to the dealer
	r.Dealer.Hands[0].Pile.AddCards(drawResult.Cards[0])
	drawResult.Cards = drawResult.Cards[1:]

	log.Debug("Cards dealt")

	// Check if the dealer has blackjack
	dealerHand := r.Dealer.Hands[0]
	if getScore(dealerHand.Pile) == 21 {
		log.Debug("Dealer blackjack")
		r.SendMessage("Dealer has blackjack")
		dealerHand.Blackjack = true
		dealerHand.Complete = true
		r.ShowTable()
		r.FSM.Event("dealerBlackjack")
		return
	}

	r.FSM.Event("showHands")
}

func (r *Round) PlayerTurn() {
	log.Debug("Running PlayerTurn")

	// Iterate over the players and each player's uncomplete hands
	for playerIdx, _ := range r.Players {
		player := &r.Players[playerIdx]
		log.Debugf("Processing turn for %v", player.UserID)
		for handIdx, _ := range player.Hands {
			hand := &player.Hands[handIdx]
			if hand.Complete {
				continue
			}

			// Play that hand and mark it as done
			r.PlayerHand(player, hand)
			hand.Complete = true
		}
	}

	log.Debug("Players done, moving onto dealer")
	r.FSM.Event("startDealer")
}

func (r *Round) PlayerHand(player *Player, hand *Hand) {
	log.Debugf("Running hand %v for %v", hand, player.UserID)

	// Get the guild object for getting nicknames
	guild, _ := r.Session.Guild(r.GuildID)
	if guild == nil {
		log.WithFields(log.Fields{
			"guild": r.GuildID,
		}).Warning("Failed to grab guild")
		return
	}

	// Get the user's nickname
	username := utils.GetPreferredName(guild, player.UserID)

	// Check for player blackjack, at this point there can only
	// be 2 cards in this hand, so 21 is a blackjack.
	if getScore(hand.Pile) == 21 {
		log.Debug("Player blackjack")
		r.SendMessage(fmt.Sprintf("%s got Blackjack!", username))
		hand.Blackjack = true
		hand.Complete = true
	}

	for !hand.Complete {
		// Handle prompt
		message := "%s's turn. Hit or Stay? (You have 30 seconds)"
		// TODO: Handle splits/double downs
		// message := "@%s's turn. Hit, Stay"
		// if checkHandCanSplit(player.Hands[0].Pile) {
		// 	message += ", Split"
		// }
		// message += " or Double down? You have 30 seconds"
		r.SendMessage(fmt.Sprintf(message, username))

		// Create a TurnTimer that will timeout after x seconds and
		// accept an action during that time.
		t := newTurnTimer(time.Second*30, []Action{
			ActionHit, ActionStay, ActionSurrender, ActionDoubleDown,
		}, player.UserID)

		var action Action

		// Store the current TurnTimer into a place where the
		// command handlers can get at it
		turnTimers[r.GuildID] = t

		// Process whether we timed-out or got a valid action
		select {
		case <-t.OnTimeout:
			// Clear the TurnTimer that was setup before
			delete(turnTimers, r.GuildID)
			r.SendMessage("No action was taken. Continuing with next hand")
			return
		case action = <-t.OnAction:
			break
		}

		// Clear the TurnTimer that was setup before
		delete(turnTimers, r.GuildID)

		log.Debugf("Taking action: %v", action.String())

		// Process the action the player chose
		switch action {
		case ActionHit:
			log.Debug("Hit")
			// Draw the new card
			drawResult := r.Deck.Draw(1)
			hand.Pile.AddCards(drawResult.Cards...)
			break
		case ActionStay:
			log.Debug("Stay")
			r.SendMessage("Stay")
			hand.Complete = true
			break
		case ActionSurrender:
			// TODO: Handle surrender
			break
		case ActionSplit:
			// TODO: Handle Split
			break
		case ActionDoubleDown:
			// TODO: Handle Doubledown
			break
		}

		r.ShowTable()

		// Check for bust hand
		if hand.CheckBust() {
			log.Debug("Hand bust")
			r.SendMessage("Bust")
			hand.Complete = true
		}
	}
}

func (r *Round) DealerTurn() {
	log.Debug("Running DealerTurn")
	hand := &r.Dealer.Hands[0]

	r.SendMessage("**Starting Dealer's turn**")

	// Flip dealer's cards to face up
	for idx := range hand.Pile.Cards {
		hand.Pile.Cards[idx].IsFaceDown = false
	}

	for !hand.Complete {
		// Get the action the dealer should take
		action := getDealerAction(hand.Pile.Cards)

		r.SendMessage(fmt.Sprintf("Dealer chooses %v", action.String()))

		// Process the action the player chose
		switch action {
		case ActionHit:
			log.Debug("Dealer Hit")
			// Draw the new card
			drawResult := r.Deck.Draw(1)
			hand.Pile.AddCards(drawResult.Cards...)
			break
		case ActionStay:
			log.Debug("Dealer stay")
			hand.Complete = true
			break
		}

		r.ShowTable()

		// Check for bust hand
		if hand.CheckBust() {
			log.Debug("Dealer hand bust")
			r.SendMessage("Dealer Bust")
			hand.Complete = true
		}
	}

	r.FSM.Event("dealerDone")
}

func (r *Round) CloseGame() {
	// Get the guild object for nicknames
	guild, _ := r.Session.Guild(r.GuildID)
	if guild == nil {
		log.WithFields(log.Fields{
			"guild": r.GuildID,
		}).Warning("Failed to grab guild")
		return
	}

	// Setup the payout ratio
	blackjackPayoutRatio := 1.6

	// Calculate the dealer score
	dealerBlackjack := r.Dealer.Hands[0].Blackjack
	dealerScore := getScore(r.Dealer.Hands[0].Pile)

	// Run over all the players
	for _, player := range r.Players {
		log.Debugf("Processing hand for %v", player.UserID)
		username := utils.GetPreferredName(guild, player.UserID)
		for _, hand := range player.Hands {
			if !hand.Complete {
				continue
			}
			// Get the score for the current hand
			handScore := getScore(hand.Pile)

			// Calculate the payout versus the dealer's hand
			payoutType := calculatePayout(hand.Blackjack, handScore,
				dealerBlackjack, dealerScore)

			// Handle player messaging
			if payoutType == PayoutLoss {
				log.Debug("Player lost")
				r.SendMessage(fmt.Sprintf("Player %v lost", username))
			} else if payoutType == PayoutPush {
				log.Debug("Player push")
				r.SendMessage(fmt.Sprintf("Player %v pushed", username))
			} else {
				log.Debug("Player won")
				r.SendMessage(fmt.Sprintf("Player %v won", username))
			}

			// Payout the bits
			payoutAmount := hand.Bet
			switch payoutType {
			case PayoutLoss:
				payoutAmount = 0
				continue
			case PayoutBlackjack:
				payoutAmount = int(float64(payoutAmount) * blackjackPayoutRatio)
			}

			// Return the bet amount
			bits.AddBits(r.Session, r.GuildID, player.UserID, hand.Bet,
				"Blackjack bet return", true)

			// If we have net zero, there is no notification of winnings
			if payoutType != PayoutPush {
				bits.AddBits(r.Session, r.GuildID, player.UserID, payoutAmount,
					"Blackjack win", false)
			}
		}
	}
	r.SendMessage("**Game complete, Thanks for playing!**")
	clearRunningRound(r.GuildID)

	// TODO: Track historical rounds
}

func (r *Round) ShowTable() {
	log.Debug("Running ShowTable")

	// Render the round
	img, err := r.Render()
	if err != nil {
		log.Error(err)
		return
	}

	// Send it to the channel
	w := &bytes.Buffer{}
	png.Encode(w, img)
	r.Session.ChannelFileSend(r.ChannelID, "png", w)
}

func (r *Round) SendMessage(content string) {
	r.Session.ChannelMessageSend(r.ChannelID, content)
}
