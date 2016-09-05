package bits

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/t11230/ramenbot/lib/ramendb"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	stateCollection       = "bitstatus"
	transactionCollection = "bittransactions"
)

type BitStatus struct {
	UserID string `bson:",omitempty"`
	Value  *int   `bson:",omitempty"`
}

type BitTransaction struct {
	UserID string
	Amount int
	Reason string
	Time   int64
}

type bitsCollection struct {
	*mgo.Collection
}

func getCollections(guildId string) (bitsCollection, bitsCollection) {
	return bitsCollection{ramendb.GetCollection(guildId, stateCollection)},
		bitsCollection{ramendb.GetCollection(guildId, transactionCollection)}
}

func AddBits(session *discordgo.Session,
	guildId string,
	userId string,
	amount int,
	reason string,
	quiet bool) {

	s, t := getCollections(guildId)

	// Update or create the current bit state
	_, err := s.Upsert(BitStatus{
		UserID: userId,
	}, bson.M{"$inc": BitStatus{
		Value: &[]int{amount}[0],
	}})
	if err != nil {
		log.Errorf("Error updating BitStatus: %v", err)
		return
	}

	err = t.Insert(BitTransaction{
		UserID: userId,
		Amount: amount,
		Reason: reason,
		Time:   time.Now().UTC().Unix(),
	})
	if err != nil {
		log.Errorf("Error inserting BitTransaction: %v", err)
		return
	}

	if quiet {
		return
	}

	message := fmt.Sprintf("You received %v bits!", amount)

	channel, _ := session.UserChannelCreate(userId)
	session.ChannelMessageSend(channel.ID, message)
}

func RemoveBits(s *discordgo.Session, guildId string, userId string, amount int, reason string) {
	AddBits(s, guildId, userId, -amount, reason, true)
}

func GetBitsLeaderboard(guildId string, count int) []BitStatus {
	s, _ := getCollections(guildId)

	var result []BitStatus
	err := s.Find(nil).Sort("-value").Limit(count).All(&result)
	if err != nil {
		log.Errorf("Error getting top BitStatuses: %v", err)
		return nil
	}

	return result
}

func GetBits(guildId string, userId string) int {
	s, _ := getCollections(guildId)

	result := &BitStatus{}
	err := s.Find(BitStatus{
		UserID: userId,
	}).One(&result)
	if err != nil {
		log.Errorf("Error getting BitStatus: %v", err)
		return 0
	}

	if result.Value == nil {
		return 0
	}

	return *result.Value
}
