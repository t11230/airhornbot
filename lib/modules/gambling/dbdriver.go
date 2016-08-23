package gambling

import (
    "github.com/t11230/ramenbot/lib/ramendb"
)

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

