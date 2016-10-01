package blackjack

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/looplab/fsm"
	"github.com/t11230/ramenbot/lib/bits"
	"github.com/t11230/ramenbot/lib/modules/gambling/cards"
	"image/png"
	"time"
)

func newRound(s *discordgo.Session, pending *PendingRound) (r *Round) {
	var err error
	r = &Round{
		Session:    s,
		GuildID:    pending.GuildID,
		ChannelID:  pending.ChannelID,
		MinimumBet: pending.MinimumBet,
		Players:    pending.Players,
	}

	r.Deck, err = cards.NewDeck(1)
	if err != nil {
		log.Errorf("Error creating deck: %v", err)
		return nil
	}

	user, err := s.User("@me")
	r.Dealer = Player{
		UserID:     user.ID,
		InitialBet: r.MinimumBet,
	}

	r.FSM = fsm.NewFSM("setup", fsm.Events{
		{Name: "deal", Src: []string{"setup"}, Dst: "dealt"},
		{Name: "showHands", Src: []string{"dealt"}, Dst: "handsShown"},
		{Name: "startPlayers", Src: []string{"handsShown"}, Dst: "playerGo"},
		{Name: "dealerBlackjack", Src: []string{"handsShown"}, Dst: "closeGame"},
		// {Name: "playerDone", Src: []string{"playerGo"}, Dst: "playerGo"},
		{Name: "startDealer", Src: []string{"playerGo"}, Dst: "dealerGo"},
		{Name: "dealerDone", Src: []string{"dealerGo"}, Dst: "closeGame"},
	}, fsm.Callbacks{
		"enter_dealt": func(e *fsm.Event) {
			go r.Deal()
		},
		"before_dealerBlackjack": func(e *fsm.Event) {
			// go r.ShowDealerBlackjack()
		},
		"enter_handsShown": func(e *fsm.Event) {
			go r.ShowTable()
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
	drawResult, err := r.Deck.Draw((len(r.Players) + 1) * 2)
	if err != nil {
		log.Errorf("Error drawing cards: %v", err)
		return
	}

	// Create the player's hands
	for i, _ := range r.Players {
		r.Players[i].AddHand(drawResult.Cards[0])
		drawResult.Cards = drawResult.Cards[1:]
	}

	// Now deal a card to the dealer
	r.Dealer.AddHand(drawResult.Cards[0])
	drawResult.Cards = drawResult.Cards[1:]

	// Deal the second card to each player
	for i, _ := range r.Players {
		r.Players[i].Hands[0].Cards = append(r.Players[i].Hands[0].Cards, drawResult.Cards[0])
		drawResult.Cards = drawResult.Cards[1:]
	}

	// Deal the second card to the dealer
	card := drawResult.Cards[0]
	card.IsFaceDown = true
	r.Dealer.Hands[0].Cards = append(r.Dealer.Hands[0].Cards, card)
	drawResult.Cards = drawResult.Cards[1:]

	log.Debug("Cards dealt")
	r.FSM.Event("showHands")
}

func (r *Round) PlayerTurn() {
	log.Debug("Running PlayerTurn")

	// Run over all the players and each player's uncomplete hands
	for playerIdx, _ := range r.Players {
		player := &r.Players[playerIdx]
		log.Debugf("Processing turn for %v", player.UserID)
		for handIdx, _ := range player.Hands {
			hand := &player.Hands[handIdx]
			if hand.Complete {
				continue
			}
			r.PlayerHand(player, hand)
			hand.Complete = true
		}
	}

	log.Debug("Players done, moving onto dealer")
	r.FSM.Event("startDealer")
}

func (r *Round) PlayerHand(player *Player, hand *Hand) {
	log.Debugf("Running hand %v for %v", hand, player.UserID)

	if getScore(hand.Cards) == 21 {
		log.Debug("Player blackjack")
		r.SendMessage("Blackjack!")
		hand.Blackjack = true
		hand.Complete = true
	}

	for !hand.Complete {
		// Handle prompt
		message := "@%s's turn. Hit, Stay"
		if checkHandCanSplit(player.Hands[0].Cards) {
			message += ", Split"
		}

		message += " or Double down? You have 30 seconds"
		r.SendMessage(fmt.Sprintf(message, player.UserID))

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
			drawResult, err := r.Deck.Draw(1)
			if err != nil {
				r.SendMessage(fmt.Sprintf("Error drawing card: %v", err))
				return
			}
			hand.Cards = append(hand.Cards, drawResult.Cards...)

			r.SendMessage(fmt.Sprintf("New card: %v of %v",
				drawResult.Cards[0].Value, drawResult.Cards[0].Suit))
			break
		case ActionStay:
			log.Debug("Stay")
			r.SendMessage("Stay")
			hand.Complete = true
			break
		case ActionSurrender:
			// Handle surrender
			break
		case ActionSplit:
			// Handle Split
			break
		case ActionDoubleDown:
			// Handle Doubledown
			break
		}

		// Check for bust hand
		if checkHandBust(hand.Cards) {
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

	if getScore(hand.Cards) == 21 {
		log.Debug("Dealer blackjack")
		r.SendMessage("Dealer has blackjack")
		hand.Blackjack = true
		hand.Complete = true
	}

	// TODO: Reveal dealer's other card

	for !hand.Complete {
		action := getDealerAction(hand.Cards)

		r.SendMessage(fmt.Sprintf("Dealer chooses %v", action.String()))

		// Process the action the player chose
		switch action {
		case ActionHit:
			log.Debug("Dealer Hit")
			// Draw the new card
			drawResult, err := r.Deck.Draw(1)
			if err != nil {
				r.SendMessage(fmt.Sprintf("Error drawing card: %v", err))
				return
			}
			hand.Cards = append(hand.Cards, drawResult.Cards...)

			r.SendMessage(fmt.Sprintf("New card: %v of %v",
				drawResult.Cards[0].Value, drawResult.Cards[0].Suit))
			break
		case ActionStay:
			log.Debug("Dealer stay")
			hand.Complete = true
			break
		}

		// Check for bust hand
		if checkHandBust(hand.Cards) {
			log.Debug("Dealer hand bust")
			r.SendMessage("Dealer Bust")
			hand.Complete = true
		}
	}

	r.FSM.Event("dealerDone")
}

func (r *Round) CloseGame() {
	blackjackPayoutRatio := 1.5

	// Calculate the dealer score
	dealerBlackjack := r.Dealer.Hands[0].Blackjack
	dealerScore := getScore(r.Dealer.Hands[0].Cards)

	// Run over all the players
	for _, player := range r.Players {
		log.Debugf("Processing hand for %v", player.UserID)
		for _, hand := range player.Hands {
			if !hand.Complete {
				continue
			}
			handScore := getScore(hand.Cards)

			payoutType := calculatePayout(hand.Blackjack, handScore,
				dealerBlackjack, dealerScore)

			if payoutType == PayoutLoss {
				log.Debug("Player lost")
				r.SendMessage(fmt.Sprintf("Player %v lost", player.UserID))
			} else {
				log.Debug("Player won")
				r.SendMessage(fmt.Sprintf("Player %v won", player.UserID))
			}

			payoutAmount := hand.Bet

			switch payoutType {
			case PayoutLoss:
				payoutAmount = 0
				continue
			case PayoutBlackjack:
				payoutAmount = int(float64(payoutAmount) * blackjackPayoutRatio)
			}

			bits.AddBits(r.Session, r.GuildID, player.UserID, payoutAmount,
				"Blackjack win", false)
		}
	}
	r.SendMessage("Game complete, Thanks for playing!")
	clearRunningRound(r.GuildID)
}

func (r *Round) ShowTable() {
	log.Debug("Running ShowTable")
	// Show the player's hands
	for _, player := range r.Players {
		img, err := cards.RenderCards(player.Hands[0].Cards)
		if err != nil {
			log.Error(err)
			return
		}

		w := &bytes.Buffer{}
		png.Encode(w, img)
		handString := fmt.Sprintf("%s's Hand is:", player.UserID)
		r.Session.ChannelMessageSend(r.ChannelID, handString)
		r.Session.ChannelFileSend(r.ChannelID, "png", w)
	}

	// Show dealers hand
	img, err := cards.RenderCards(r.Dealer.Hands[0].Cards)
	if err != nil {
		log.Error(err)
		return
	}

	w := &bytes.Buffer{}
	png.Encode(w, img)
	handString := fmt.Sprintf("Dealer's Hand is:")
	r.Session.ChannelMessageSend(r.ChannelID, handString)
	r.Session.ChannelFileSend(r.ChannelID, "png", w)

	r.FSM.Event("startPlayers")
}

func (r *Round) SendMessage(content string) {
	r.Session.ChannelMessageSend(r.ChannelID, content)
}
