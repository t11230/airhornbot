package bits

import (
    "bytes"
    "fmt"
    "text/tabwriter"
    log "github.com/Sirupsen/logrus"

    "github.com/bwmarrin/discordgo"

    "github.com/t11230/ramenbot/lib/utils"
    "github.com/t11230/ramenbot/lib/rdb"
)

type BitStat struct {
    UserID string
    BitValue int
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
