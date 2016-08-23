package gametrack

import (
    "github.com/t11230/ramenbot/lib/ramendb"
)

/* Game tracking functions */
func (db *BotDatabase) GetGameTrackCollection() *mgo.Collection {
    return db.Session.DB(fmt.Sprintf("%s-%s",
                         MongoBaseDBName,
                         db.GuildID)).C(MongoGameTrackCollection)
}

func (db *BotDatabase) GameTrackIncGameEntry(userID string, game string, inc int) error {
    c := db.GetGameTrackCollection()

    search := bson.M{
        "userid" : userID,
        "game": game,
    }

    update := bson.M{
        "$inc": bson.M{
            "time": inc,
        },
    }

    _, err := c.Upsert(search, update)
    if err != nil {
        log.Error(err)
        return err
    }
    return nil
}

func (db *BotDatabase) GameTrackGetUserStats(userID string, count int) []GameTrackEntry {
    c := db.GetGameTrackCollection()

    var result []GameTrackEntry
    err := c.Find(bson.M{"userid" : userID}).Sort("-time").Limit(count).All(&result)
    if err != nil {
        log.Error(err)
        return nil
    }

    return result
}

func (db *BotDatabase) GameTrackGetTopTimes(count int) []GameTrackEntry {
    c := db.GetGameTrackCollection()

    var result []GameTrackEntry
    err := c.Find(nil).Sort("-time").Limit(count).All(&result)
    if err != nil {
        log.Error(err)
        return nil
    }

    return result
}