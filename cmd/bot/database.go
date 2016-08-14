package main

import (
    "errors"
    "fmt"

    log "github.com/Sirupsen/logrus"

    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

const (
    MongoBaseDBName = "discord"
    MongoBitsCollection = "bits"
    MongoGameTrackCollection = "gametrack"
    MongoBetRollCollection = "betroll"
)

var (
    mongoDatabase *BotDatabase
)

type BotDatabase struct {
    Session *mgo.Session
    GuildID string
}

/* New Mongo DB functions */
/* Example:
    log.Info("Updating bits")
    db := dbGetSession()
    db.SetBitStats("1", "2", int)
    bits := db.GetBitStats("1", "2")
 */
func dbMongoOpen(serverURL string) {
    db, err := mgo.Dial(serverURL)
    if err != nil {
        log.Fatalf("MongoDB Open: %s\n", err)
        return
    }

    db.SetMode(mgo.Monotonic, true)

    mongoDatabase = &BotDatabase {
        Session: db,
    }
}

func dbGetSession(guildID string) *BotDatabase {
    return &BotDatabase {
        Session: mongoDatabase.Session.Copy(),
        GuildID: guildID,
    }
}

/* Bit stats functions */
func (db *BotDatabase) GetBitsCollection() *mgo.Collection {
    return db.Session.DB(fmt.Sprintf("%s-%s",
                         MongoBaseDBName,
                         db.GuildID)).C(MongoBitsCollection)
}

func (db *BotDatabase) UpsertBitStats(userID string, update interface{}) {
    c := db.GetBitsCollection()

    _, err := c.Upsert(bson.M{"userid": userID}, update)
    if err != nil {
        log.Error(err)
    }
}

func (db *BotDatabase) GetBitStats(userID string) *BitStat {
    c := db.GetBitsCollection()

    // Retrieve the bits for the current user
    result := &BitStat{}
    err := c.Find(bson.M{"userid": userID}).One(&result)
    if err != nil {
        log.Info(err)
        return nil
    }

    return result
}

func (db *BotDatabase) SetBitStats(userID string, value int) {
    update := bson.M{
        "$set": BitStat{
            UserID: userID,
            BitValue: value,
        },
    }

    db.UpsertBitStats(userID, update)
}

func (db *BotDatabase) IncBitStats(userID string, value int) {
    update := bson.M{
        "$inc": bson.M{
            "bitvalue": value,
        },
    }

    db.UpsertBitStats(userID, update)
}

func (db *BotDatabase) DecBitStats(userID string, value int) {
    db.IncBitStats(userID, -value)
}

func (db *BotDatabase) DecCheckBitStats(userID string, value int) (error) {
    b := db.GetBitStats(userID)
    if b == nil {
        b = &BitStat{UserID: userID, BitValue: 0}
        db.SetBitStats(userID, b.BitValue)
    }

    if (b.BitValue - value < 0) {
        return errors.New("Not enough bits")
    }

    db.DecBitStats( userID, value)
    return nil
}

func (db *BotDatabase) GetTopBitStats(count int) []BitStat {
    c := db.GetBitsCollection()

    var result []BitStat
    err := c.Find(nil).Sort("-bitvalue").Limit(count).All(&result)
    if err != nil {
        log.Error(err)
        return nil
    }

    return result
}

/* Betting functions */
func (db *BotDatabase) GetBetRollCollection(guildID string) *mgo.Collection {
    return db.Session.DB(fmt.Sprintf("%s-%s",
                         MongoBaseDBName,
                         guildID)).C(MongoBetRollCollection)
}

func (db *BotDatabase) GetActiveBetRoll(guildID string) *BetRoll {
    c := db.GetBetRollCollection(guildID)

    result := &BetRoll{}
    err := c.Find(nil).One(&result)
    if err != nil {
        return nil
    }
    return result
}

func (db *BotDatabase) GetPlayers(guildID string) []Player {
    b := db.GetActiveBetRoll(guildID)
    if b == nil {
        log.Error(errors.New("No active BetRoll"))
        return nil
    }

    return b.Players
}

func (db *BotDatabase) SetBetRollAnte(guildID string, ante int) error {
    b := db.GetActiveBetRoll(guildID)
    if b == nil {
        return errors.New("No active BetRoll")
    }

    b.Ante = ante

    c := db.GetBetRollCollection(guildID)

    update := bson.M{
        "$set": b,
    }

    err := c.Update(bson.M{}, update)
    if err != nil {
        return err
    }

    return nil
}

func (db *BotDatabase) BetRollAddPlayer(guildID string, player Player) error {
    b := db.GetActiveBetRoll(guildID)
    if b == nil {
        log.Error(errors.New("No active BetRoll"))
        return nil
    }

    b.Players = append(b.Players, player)

    c := db.GetBetRollCollection(guildID)

    update := bson.M{
        "$set": b,
    }

    err := c.Update(bson.M{}, update)
    if err != nil {
        return err
    }


    return nil
}

func (db *BotDatabase) BetRollOpen(guildID string) error {
    b := db.GetActiveBetRoll(guildID)
    if b != nil {
        return errors.New("An active BetRoll already exists")
    }

    c := db.GetBetRollCollection(guildID)
    return c.Insert(BetRoll{})
}

func (db *BotDatabase) BetRollClose(guildID string) error {
    b := db.GetActiveBetRoll(guildID)
    if b == nil {
        return errors.New("No active BetRoll")
    }

    c := db.GetBetRollCollection(guildID)
    return c.Remove(bson.M{})
}


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