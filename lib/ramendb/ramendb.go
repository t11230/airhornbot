package ramendb

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

const (
	MongoBaseDBName = "discord"
)

var (
	mongoDatabase *BotDatabase
)

type BotDatabase struct {
	Session *mgo.Session
	GuildID string
	Module  string
}

func MongoOpen(serverURL string) {
	db, err := mgo.Dial(serverURL)
	if err != nil {
		log.Errorf("MongoDB Open: %s\n", err)
		return
	}

	db.SetMode(mgo.Monotonic, true)

	mongoDatabase = &BotDatabase{
		Session: db,
	}
}

func GetCollection(guildID string, moduleName string) *mgo.Collection {
	s := mongoDatabase.Session.Copy()
	return s.DB(fmt.Sprintf("%s-%s",
		MongoBaseDBName,
		guildID)).C(moduleName)

}
