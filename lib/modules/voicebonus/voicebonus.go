package voicebonus

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/t11230/ramenbot/lib/bits"
	"github.com/t11230/ramenbot/lib/modules/modulebase"
	"github.com/t11230/ramenbot/lib/perms"
	"github.com/t11230/ramenbot/lib/ramendb"
	"github.com/t11230/ramenbot/lib/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"strings"
	"time"
)

const (

	ConfigName = "voicebonus"

	joinMessage = `
	Welcome **%s**, you get %d bits for joining this week!
	You now have **%d bits**`

	helpString = "**!!vb** : This module allows the user to control the bit bonus for joining a voice channel.\n"

	vbHelpString = `**VB**
This module allows the user to control the bit bonus for joining a voice channel.

**usage:** !!vb set *function* *args...*  **OR** !!vb set *status*

**permissions required:** voicebonus-control

If true or false is passed as the argument to !!vb set, then it will turn the voice bonus on or off accordingly.
Otherwise, it will call the function that is passed to it with the arguments following it.

**functions:**
    *amount:* This function allows the user to set the bit amount of the voice bonus.
	*time:* This function allows the user to when a user must join the call to receive the voice bonus.

For more info on using any of these functions, type **!!vb set [function name] help**`

	amountHelpString = `**AMOUNT**
This function allows the user to set the bit amount of the voice bonus.

**usage:** !!vb set amount *value*
	Sets the value of the voice bonus to the integer value specified by *value*`

	timeHelpString = `**TIME**
This function allows the user to set when a user must join the call to receive the voice bonus.

**usage:** !!vb set time *weekday* *starttime* *endtime*
	Sets the time range of the voice bonus from the UTC time *starttime* to the UTC time *endtime* every week on *weekday*`
)
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
		DBStart:  handleDbStart,
		Help:     helpString,
	}, nil
}

func handleDbStart() error {
	perms.CreatePerm("voicebonus-control")
	return nil
}

func handleSet(cmd *modulebase.ModuleCommand) (string, error) {
	log.Debug("Called handleSet")

	if len(cmd.Args) == 0 || cmd.Args[0]=="help" {
		return vbHelpString, nil
	}

	permsHandle := perms.GetPermsHandle(cmd.Guild.ID, ConfigName)
	if !permsHandle.CheckPerm(cmd.Message.Author.ID, "voicebonus-control") {
		return "Insufficient permissions", nil
	}

	c := voicebonusCollection{ramendb.GetCollection(cmd.Guild.ID, ConfigName)}
	enable, err := utils.EnableToBool(cmd.Args[0])
	if err != nil {
		return vbHelpString, nil
	}

	err = c.Enable(enable)
	if err != nil {
		errString := fmt.Sprintf("Error updating enable state: %v", err)
		return errString, nil
	}

	var enabledString string
	if enable {
		enabledString = "Enabled"
	} else {
		enabledString = "Disabled"
	}
	return fmt.Sprintf("Voice Bonus %v", enabledString), nil
}

func handleSetAmount(cmd *modulebase.ModuleCommand) (string, error) {
	log.Debug("Called handleSetAmount")

	permsHandle := perms.GetPermsHandle(cmd.Guild.ID, ConfigName)
	if !permsHandle.CheckPerm(cmd.Message.Author.ID, "voicebonus-control") {
		return "Insufficient permissions", nil
	}

	if len(cmd.Args) == 0 {
		return amountHelpString, nil
	}

	amount, err := strconv.Atoi(cmd.Args[0])
	if err != nil {
		return amountHelpString, nil
	}

	c := voicebonusCollection{ramendb.GetCollection(cmd.Guild.ID, ConfigName)}
	err = c.SetAmount(amount)
	if err != nil {
		errString := fmt.Sprintf("Error updating amount: %v", err)
		return errString, nil
	}

	return fmt.Sprintf("Voice Bonus amount set to %v", amount), nil
}

func handleSetTime(cmd *modulebase.ModuleCommand) (string, error) {
	log.Debug("Called handleSetTime")

	permsHandle := perms.GetPermsHandle(cmd.Guild.ID, ConfigName)
	if !permsHandle.CheckPerm(cmd.Message.Author.ID, "voicebonus-control") {
		return "Insufficient permissions", nil
	}

	if len(cmd.Args) != 3 || cmd.Args[0]=="help" {
		return timeHelpString, nil
	}

	weekday := utils.ToWeekday(cmd.Args[0])
	time := strings.Split(cmd.Args[1], ":")
	hour, err := strconv.Atoi(time[0])
	if err != nil {
		return "Invalid hour", nil
	}
	minute, err := strconv.Atoi(time[1])
	if err != nil {
		return "Invalid minute", nil
	}

	if !strings.Contains(cmd.Args[2], "h") {
		return "Missing 'h' in duration", nil
	}
	duration, err := strconv.Atoi(strings.Replace(cmd.Args[2], "h", "", 1))
	if err != nil {
		return "Invalid duration", nil
	}

	span := voicebonusTimespan{
		Weekday:  int(weekday),
		Hour:     hour,
		Minute:   minute,
		Duration: duration,
	}

	c := voicebonusCollection{ramendb.GetCollection(cmd.Guild.ID, ConfigName)}
	err = c.SetTimespan(span)

	if err != nil {
		errString := fmt.Sprintf("Error updating time: %v", err)
		return errString, nil
	}
	return fmt.Sprintf("Voice Bonus time set to %v:%v on %v for %v hours", hour, minute, cmd.Args[0], duration), nil
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

	log.Debugf("Last join date: %v, now: %v", lastJoinTime, currentTime)

	log.Debugf("%v %v %v", currentTime.After(startDate), time.Since(startDate).Hours() < float64(span.Duration), lastJoinTime.Before(startDate))

	if currentTime.After(startDate) &&
		time.Since(startDate).Hours() < float64(span.Duration) &&
		lastJoinTime.Before(startDate) {
		log.Debug("Giving bits for join")

		bits.AddBits(s, v.GuildID, v.UserID, c.Amount(), "Voice join bonus", true)

		username := utils.GetPreferredName(guild, v.UserID)
		message := fmt.Sprintf(joinMessage, username, c.Amount(),
			bits.GetBits(v.GuildID, v.UserID))

		channel, _ := s.UserChannelCreate(v.UserID)
		s.ChannelMessageSend(channel.ID, message)

		updateUserLastJoin(v.GuildID, v.UserID, currentTime.Unix())
	}
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
