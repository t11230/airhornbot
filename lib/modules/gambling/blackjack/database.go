package blackjack

import (
	log "github.com/Sirupsen/logrus"
	"github.com/t11230/ramenbot/lib/ramendb"
	"gopkg.in/mgo.v2"
)

const (
	collectionBaseName   = "blackjack"
	pendingRoundDBSuffix = "pending"
	runningRoundDBSuffix = "running"
)

func getPendingRoundCollection(guildId string) *mgo.Collection {
	return ramendb.GetCollection(guildId,
		collectionBaseName+"_"+pendingRoundDBSuffix)
}

func getRunningRoundCollection(guildId string) *mgo.Collection {
	return ramendb.GetCollection(guildId,
		collectionBaseName+"_"+runningRoundDBSuffix)
}

// Create round functions
func createPendingRound(guildId string, round *PendingRound) (err error) {
	c := getPendingRoundCollection(guildId)
	err = c.Insert(round)
	if err != nil {
		log.Errorf("Error adding pending round: %v", err)
		return
	}
	return
}

func createRunningRound(guildId string, round *Round) (err error) {
	c := getRunningRoundCollection(guildId)
	err = c.Insert(round)
	if err != nil {
		log.Errorf("Error adding pending round: %v", err)
		return
	}
	return
}

// Get round functions
func getPendingRound(guildId string) (round *PendingRound, err error) {
	c := getPendingRoundCollection(guildId)
	round = &PendingRound{}
	err = c.Find(nil).One(&round)
	if err != nil {
		round = nil
		return
	}
	return
}

func getRunningRound(guildId string) (round *Round, err error) {
	c := getRunningRoundCollection(guildId)
	round = &Round{}
	err = c.Find(nil).One(&round)
	if err != nil {
		round = nil
		return
	}
	return
}

// Update round functions
func updatePendingRound(guildId string, round *PendingRound) (err error) {
	c := getPendingRoundCollection(guildId)
	err = c.Remove(nil)
	if err != nil {
		log.Errorf("Error removing round: %v", err)
		return
	}

	err = c.Insert(round)
	if err != nil {
		log.Errorf("Error adding round: %v", err)
		return
	}
	return
}

func updateRunningRound(guildId string, round *Round) (err error) {
	c := getRunningRoundCollection(guildId)
	err = c.Remove(nil)
	if err != nil {
		log.Errorf("Error removing round: %v", err)
		return
	}

	err = c.Insert(round)
	if err != nil {
		log.Errorf("Error adding round: %v", err)
		return
	}
	return
}

// Clear round functions
func clearPendingRound(guildId string) (err error) {
	return clearRound(getPendingRoundCollection(guildId))
}

func clearRunningRound(guildId string) (err error) {
	return clearRound(getRunningRoundCollection(guildId))
}

func clearRound(c *mgo.Collection) (err error) {
	err = c.Remove(nil)
	if err != nil {
		log.Errorf("Error removing round: %v", err)
		return
	}
	return
}
