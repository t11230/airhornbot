package blackjack

import (
	log "github.com/Sirupsen/logrus"
	"github.com/t11230/ramenbot/lib/modules/gambling/cards"
)

func (h *Hand) GetAvailableActions() (actions []Action) {
	actions = []Action{ActionSurrender, ActionStay, ActionHit}
	if len(h.Pile.Cards) == 2 &&
		h.Pile.Cards[0].NumericValue(cards.BlackjackAceHighMap) ==
			h.Pile.Cards[1].NumericValue(cards.BlackjackAceHighMap) {
		actions = append(actions, ActionSplit)
	}

	if len(h.Pile.Cards) == 2 {
		actions = append(actions, ActionDoubleDown)
	}
	return
}

func (p *Player) AddHand(c ...cards.Card) {
	log.Debug("Adding hand")
	if p.Hands == nil {
		p.Hands = []Hand{}
	}
	newHand := Hand{
		Bet:       p.InitialBet,
		Pile:      cards.Pile{c},
		Complete:  false,
		Blackjack: false,
	}
	p.Hands = append(p.Hands, newHand)

	return
}

func (h *Hand) CheckBust() (bust bool) {
	bust = false

	a := h.Pile.Sum(cards.BlackjackAceLowMap)
	log.Debugf("A Score: %v", a)

	if a <= 21 {
		return
	}

	b := h.Pile.Sum(cards.BlackjackAceHighMap)

	log.Debugf("B Score: %v", b)

	if b <= 21 {
		return
	}

	return true
}

// func checkHandBlackjack(c []cards.Card) (blackjack bool) {
// 	p := cards.Pile{Cards: c}

// 	blackjack = true
// 	if p.Sum(cards.BlackjackAceHighMap) == 21 {
// 		return
// 	}

// 	if p.Sum(cards.BlackjackAceLowMap) == 21 {
// 		return
// 	}

// 	return false
// }

func getDealerAction(c []cards.Card) (action Action) {
	action = ActionHit
	p := cards.Pile{Cards: c}

	// For now lets hit on soft 17
	if p.Sum(cards.BlackjackAceHighMap) <= 17 {
		return
	}

	if p.Sum(cards.BlackjackAceLowMap) < 17 {
		return
	}

	return ActionStay
}

func getScore(p cards.Pile) int {
	aceLowScore := p.Sum(cards.BlackjackAceLowMap)
	aceHighScore := p.Sum(cards.BlackjackAceHighMap)

	if aceHighScore > aceLowScore && aceHighScore <= 21 {
		return aceHighScore
	}

	if aceLowScore <= 21 {
		return aceLowScore
	}

	return 0
}

func checkHandCanSplit(p cards.Pile) (canSplit bool) {
	if len(p.Cards) != 2 {
		return false
	}
	if p.Cards[0].NumericValue(cards.BlackjackAceHighMap) ==
		p.Cards[1].NumericValue(cards.BlackjackAceHighMap) {
		return true
	}

	return false
}

func calculatePayout(pBlackjack bool, pScore int, dBlackjack bool, dScore int) Payout {
	// Check dealer blackjack
	if dBlackjack {
		if pBlackjack {
			// Return bet if its a push
			return PayoutPush
		}
		return PayoutLoss
	}

	// Player blackjack, and no dealer blackjack
	if pBlackjack {
		return PayoutBlackjack
	}

	if pScore > 21 {
		return PayoutLoss
	}

	// Check standard cases, no blackjack
	if pScore > dScore {
		return PayoutWin
	} else if pScore < dScore {
		return PayoutLoss
	}

	return PayoutPush
}

func (a *Action) String() string {
	switch *a {
	case ActionHit:
		return "Hit"
	case ActionStay:
		return "Stay"
	case ActionSplit:
		return "Split"
	case ActionDoubleDown:
		return "Double Down"
	case ActionSurrender:
		return "Surrender"
	}

	return "Unknown Action"
}
