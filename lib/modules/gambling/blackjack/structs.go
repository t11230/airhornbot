package blackjack

import (
	"github.com/bwmarrin/discordgo"
	"github.com/looplab/fsm"
	"github.com/t11230/ramenbot/lib/modules/gambling/cards"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	ActionHit Action = iota
	ActionStay
	ActionDoubleDown
	ActionSplit
	ActionSurrender
)

const (
	PayoutLoss Payout = iota
	PayoutWin
	PayoutBlackjack
)

type Action int

type Payout int

type historyCollection struct{ *mgo.Collection }

type Hand struct {
	Cards     []cards.Card `bson:",omitempty"`
	Bet       int          `bson:",omitempty"`
	Complete  bool         `bson:",omitempty"`
	Blackjack bool         `bson:",omitempty"`
}

type Player struct {
	UserID       string `bson:",omitempty"`
	Hands        []Hand `bson:",omitempty"`
	InitialBet   int    `bson:",omitempty"`
	TurnComplete bool   `bson:",omitempty"`
}

type PendingRound struct {
	ID         bson.ObjectId `bson:"_id,omitempty"`
	Countdown  int           `bson:",omitempty"`
	GuildID    string        `bson:",omitempty"`
	ChannelID  string        `bson:",omitempty"`
	MinimumBet int           `bson:",omitempty"`

	Players []Player `bson:",omitempty"`
}

type Round struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	GuildID   string        `bson:",omitempty"`
	ChannelID string        `bson:",omitempty"`

	Session *discordgo.Session `bson:"-"`
	FSM     *fsm.FSM           `bson:"-"`

	Deck       *cards.Deck `bson:",omitempty"`
	Players    []Player    `bson:",omitempty"`
	MinimumBet int         `bson:",omitempty"`
	Dealer     Player      `bson:",omitempty"`

	StartTime time.Time `bson:",omitempty"`
}

type HistoricalRound struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	GuildID   string        `bson:",omitempty"`
	ChannelID string        `bson:",omitempty"`

	Players    []Player `bson:",omitempty"`
	MinimumBet int      `bson:",omitempty"`
	Dealer     Player   `bson:",omitempty"`

	WinnerID  string    `bson:",omitempty"`
	StartTime time.Time `bson:",omitempty"`
}

type TurnTimer struct {
	ValidActions []Action
	OnTimeout    <-chan time.Time
	OnAction     chan Action
	UserID       string
}
