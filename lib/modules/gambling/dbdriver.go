package gambling

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/t11230/ramenbot/lib/ramendb"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	betRollDBSuffix = "betroll"
)

type gamblingCollection struct {
	*mgo.Collection
}

type Player struct {
	UserID string
	Bid    int
}

type BetRoll struct {
	Players []Player
	Ante    int
}

func getBetRollCollection(guildId string) *gamblingCollection {
	return &gamblingCollection{ramendb.GetCollection(guildId, ConfigName+betRollDBSuffix)}
}

func getActiveBetRoll(guildId string) *BetRoll {
	c := getBetRollCollection(guildId)
	result := &BetRoll{}
	err := c.Find(nil).One(&result)
	if err != nil {
		return nil
	}
	return result
}

func betRollAddPlayer(guildId string, player *Player) error {
	b := getActiveBetRoll(guildId)
	if b == nil {
		log.Error(errors.New("No active BetRoll"))
		return nil
	}

	b.Players = append(b.Players, *player)

	c := getBetRollCollection(guildId)

	update := bson.M{
		"$set": b,
	}

	err := c.Update(bson.M{}, update)
	if err != nil {
		return err
	}

	return nil
}

func betRollOpen(guildId string) error {
	b := getActiveBetRoll(guildId)
	if b != nil {
		return errors.New("An active BetRoll already exists")
	}

	c := getBetRollCollection(guildId)
	return c.Insert(BetRoll{})
}

func setBetRollAnte(guildId string, ante int) error {
	b := getActiveBetRoll(guildId)
	if b == nil {
		return errors.New("No active BetRoll")
	}

	b.Ante = ante

	c := getBetRollCollection(guildId)

	update := bson.M{
		"$set": b,
	}

	err := c.Update(bson.M{}, update)
	if err != nil {
		return err
	}

	return nil
}

func betRollClose(guildId string) error {
	b := getActiveBetRoll(guildId)
	if b == nil {
		return errors.New("No active BetRoll")
	}

	c := getBetRollCollection(guildId)
	return c.Remove(bson.M{})
}

func getBetRollPlayers(guildId string) []Player {
	b := getActiveBetRoll(guildId)
	if b == nil {
		log.Error(errors.New("No active BetRoll"))
		return nil
	}

	return b.Players
}
