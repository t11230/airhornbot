package blackjack

import (
	log "github.com/Sirupsen/logrus"
	"time"
)

func newTurnTimer(timeout time.Duration, validActions []Action, userId string) *TurnTimer {
	return &TurnTimer{
		OnTimeout:    time.After(timeout),
		OnAction:     make(chan Action),
		ValidActions: validActions,
		UserID:       userId,
	}
}

func (t *TurnTimer) SendAction(action Action) bool {
	log.Debug("SendAction")
	log.Debugf("Available Actions %v", t.ValidActions)
	log.Debugf("Action is %v", action)
	for _, a := range t.ValidActions {
		if action == a {
			t.OnAction <- action
			return true
		}
	}

	return false
}

func (t *TurnTimer) CheckUser(userId string) bool {
	log.Debug("CheckUser")
	if userId != t.UserID {
		log.Debug("Wrong user")
		return false
	}
	return true
}
