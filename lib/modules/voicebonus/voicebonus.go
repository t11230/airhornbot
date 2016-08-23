package voicebonus
    // startTime := time.Date(2016, time.August, 16, 23, 0, 0, 0, time.UTC)
    // endTime := time.Date(2016, time.August, 17, 5, 0, 0, 0, time.UTC)

    // if utils.InTimeSpan(startTime, endTime, time.Now().UTC()) {
    //     db := rdb.GetSession(guild.ID)
    //     e := db.GetVoiceJoinEntry(member.User.ID)

    //     if e != nil {
    //         for _, t := range e.Dates {
    //             p := time.Unix(t, 0).UTC()
    //             // log.Info(p)
    //             if utils.InTimeSpan(startTime, endTime, p) {
    //                 return
    //             }
    //         }
    //     }

    //     db.UpsertVoiceJoinEntry(member.User.ID)

    //     // // Give weekly bit bonus
    //     // message:= giveWeeklyBitBonus(guild, member.User.ID)
    //     // c,_ := s.UserChannelCreate(member.User.ID)
    //     // s.ChannelMessageSend(c.ID, message)
    // }