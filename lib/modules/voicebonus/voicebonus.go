package voicebonus

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/t11230/ramenbot/lib/bits"
	"github.com/t11230/ramenbot/lib/modules/modulebase"
	"github.com/t11230/ramenbot/lib/ramendb"
	"github.com/t11230/ramenbot/lib/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"strings"
	"time"
)

const ConfigName = "voicebonus"

const joinMessage = `
Welcome **%s**, you get %d bits for joining this week!
You now have **%d bits**`

// List of commands that this module accepts
var commandTree = []modulebase.ModuleCommandTree{
	{
		RootCommand: "vb",
		SubKeys: modulebase.SK{
			"set": modulebase.CN{
				SubKeys: modulebase.SK{
					"amount": modulebase.CN{
						Function: handleSetAmount,
					},
					"time": modulebase.CN{
						Function: handleSetTime,
					},
				},
			},
		},
		Function: handleSet,
	},
}

// Called to initialize this module
func SetupFunc(config *modulebase.ModuleConfig) (*modulebase.ModuleSetupInfo, error) {
	events := []interface{}{
		voiceStateUpdateCallback,
	}

	return &modulebase.ModuleSetupInfo{
		Events:   &events,
		Commands: &commandTree,
	}, nil
}

func handleSet(cmd *modulebase.ModuleCommand) error {
	log.Debug("Called handleSet")
	if len(cmd.Args) == 0 {
		err := errors.New("Missing Args")
		log.Error(err)
		return err
	}

	c := voicebonusCollection{ramendb.GetCollection(cmd.Guild.ID, ConfigName)}
	enable, err := utils.EnableToBool(cmd.Args[0])
	if err != nil {
		return err
	}

	return c.Enable(enable)
}

func handleSetAmount(cmd *modulebase.ModuleCommand) error {
	log.Debug("Called handleSetAmount")
	if len(cmd.Args) == 0 {
		err := errors.New("Missing Args")
		log.Error(err)
		return err
	}

	amount, err := strconv.Atoi(cmd.Args[0])
	if err != nil {
		log.Error(err)
		return err
	}

	c := voicebonusCollection{ramendb.GetCollection(cmd.Guild.ID, ConfigName)}
	return c.SetAmount(amount)
}

func handleSetTime(cmd *modulebase.ModuleCommand) error {
	log.Debug("Called handleSetTime")
	if len(cmd.Args) != 3 {
		err := errors.New("Missing or invalid args")
		log.Error(err)
		return err
	}

	weekday := utils.ToWeekday(cmd.Args[0])
	time := strings.Split(cmd.Args[1], ":")
	hour, err := strconv.Atoi(time[0])
	if err != nil {
		err := errors.New("Invalid hour")
		return err
	}
	minute, err := strconv.Atoi(time[1])
	if err != nil {
		err := errors.New("Invalid minute")
		return err
	}

	if !strings.Contains(cmd.Args[2], "h") {
		err := errors.New("Invalid duration")
		return err
	}
	duration, err := strconv.Atoi(strings.Replace(cmd.Args[2], "h", "", 1))
	if err != nil {
		err := errors.New("Invalid duration")
		return err
	}

	span := voicebonusTimespan{
		Weekday:  int(weekday),
		Hour:     hour,
		Minute:   minute,
		Duration: duration,
	}

	c := voicebonusCollection{ramendb.GetCollection(cmd.Guild.ID, ConfigName)}
	return c.SetTimespan(span)
}

func voiceStateUpdateCallback(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	log.Debugf("Voice bonus On voice state update: %v", v.VoiceState)

	// Check if it was a part
	if v.ChannelID == "" {
		return
	}

	guild, _ := s.State.Guild(v.GuildID)
	if guild == nil {
		log.WithFields(log.Fields{
			"guild": v.GuildID,
		}).Warning("Failed to grab guild")
		return
	}

	member, _ := s.State.Member(v.GuildID, v.UserID)
	if member == nil {
		log.WithFields(log.Fields{
			"member": member,
		}).Warning("Failed to grab member")
		return
	}

	if member.User.Bot {
		return
	}

	c := voicebonusCollection{ramendb.GetCollection(v.GuildID, ConfigName)}
	if !c.Enabled() {
		return
	}

	span := c.Timespan()
	if span == nil {
		log.Error("Timespan was nil for %v", guild.ID)
		return
	}

	log.Debugf("Timespan is: %v", span)

	currentTime := time.Now().UTC()

	year, month, day := currentTime.Date()
	weekday := currentTime.Weekday()
	nextDay := day + utils.GetDaysTillWeekday(int(weekday), span.Weekday)
	startDate := time.Date(year, month, nextDay, span.Hour, span.Minute, 0, 0, time.UTC)

	log.Debugf("Start Date is: %v", startDate)

	lastJoinTime := time.Unix(getUserLastJoin(v.GuildID, v.UserID), 0)

	if currentTime.After(startDate) &&
		time.Since(startDate).Hours() < float64(span.Duration) &&
		lastJoinTime.Before(startDate) {

		bits.AddBits(s, v.GuildID, v.UserID, c.Amount(), "Voice join bonus", true)

		username := utils.GetPreferredName(guild, v.UserID)
		message := fmt.Sprintf(joinMessage, username, c.Amount(),
			bits.GetBits(v.GuildID, v.UserID))

		channel, _ := s.UserChannelCreate(v.UserID)
		s.ChannelMessageSend(channel.ID, message)
	}

	updateUserLastJoin(v.GuildID, v.UserID, currentTime.Unix())
}

type voicebonusLastJoin struct {
	UserID   string `bson:",omitempty"`
	LastJoin *int64 `bson:",omitempty"`
}

func updateUserLastJoin(guildId string, userId string, time int64) {
	c := ramendb.GetCollection(guildId, ConfigName+"lastjoins")
	c.Upsert(voicebonusLastJoin{UserID: userId},
		bson.M{"$set": voicebonusLastJoin{LastJoin: &time}})
}

func getUserLastJoin(guildId string, userId string) int64 {
	c := ramendb.GetCollection(guildId, ConfigName+"lastjoins")
	result := &voicebonusLastJoin{}
	c.Find(voicebonusLastJoin{UserID: userId}).One(result)
	if result.LastJoin == nil {
		return 0
	}

	return *result.LastJoin
}

type voicebonusCollection struct {
	*mgo.Collection
}

type voicebonusTimespan struct {
	Weekday  int
	Hour     int
	Minute   int
	Duration int
}

type voicebonusConfig struct {
	Enabled  *bool               `bson:",omitempty"`
	Amount   *int                `bson:",omitempty"`
	Timespan *voicebonusTimespan `bson:",omitempty"`
}

func (c *voicebonusCollection) UpdateConfig(update interface{}) error {
	_, err := c.Upsert(nil, update)
	if err != nil {
		log.Errorf("Error updating voicebonusConfig %v", err)
		return err
	}
	return nil
}

func (c *voicebonusCollection) Enable(enable bool) error {
	updateData := bson.M{"$set": voicebonusConfig{
		Enabled: &enable,
	}}
	return c.UpdateConfig(updateData)
}

func (c *voicebonusCollection) SetAmount(amount int) error {
	updateData := bson.M{"$set": voicebonusConfig{
		Amount: &amount,
	}}
	return c.UpdateConfig(updateData)
}

func (c *voicebonusCollection) SetTimespan(span voicebonusTimespan) error {
	updateData := bson.M{"$set": voicebonusConfig{
		Timespan: &span,
	}}
	return c.UpdateConfig(updateData)
}

func (c *voicebonusCollection) Enabled() bool {
	config := voicebonusConfig{}
	c.Find(nil).One(&config)

	if config.Enabled != nil {
		return *config.Enabled
	}
	return false
}

func (c *voicebonusCollection) Amount() int {
	config := voicebonusConfig{}
	c.Find(nil).One(&config)

	if config.Amount != nil {
		return *config.Amount
	}
	return 0
}

func (c *voicebonusCollection) Timespan() *voicebonusTimespan {
	config := voicebonusConfig{}
	c.Find(nil).One(&config)

	return config.Timespan
}
