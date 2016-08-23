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

type VoiceJoinEntry struct {
	UserID string
	Dates  []int64
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

/* Weekly call tracking functions */
// func (db *BotDatabase) GetVoiceJoinCollection() *mgo.Collection {

// }

// func (db *BotDatabase) GetVoiceJoinEntry(userID string) *VoiceJoinEntry {
// 	c := db.GetVoiceJoinCollection()

// 	search := bson.M{
// 		"userid": userID,
// 	}

// 	var result VoiceJoinEntry
// 	err := c.Find(search).One(&result)
// 	if err != nil {
// 		log.Error(err)
// 		return nil
// 	}
// 	return &result
// }

// func (db *BotDatabase) UpsertVoiceJoinEntry(userID string) error {
// 	c := db.GetVoiceJoinCollection()

// 	search := bson.M{
// 		"userid": userID,
// 	}

// 	update := bson.M{
// 		"$push": bson.M{
// 			"dates": time.Now().UTC().Unix(),
// 		},
// 	}

// 	_, err := c.Upsert(search, update)
// 	if err != nil {
// 		log.Error(err)
// 		return err
// 	}
// 	return nil
// }
