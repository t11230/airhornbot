package main

import (
    "bytes"
    "fmt"
    "text/tabwriter"
    log "github.com/Sirupsen/logrus"

    "github.com/bwmarrin/discordgo"
)

type BitStat struct {
    UserID string
    BitValue int
}

func bitsPrintStats(guild *discordgo.Guild, message *discordgo.Message, args []string) string {
    user := message.Author
    db := dbGetSession(guild.ID)
    var bits BitStat
    var bitslist []BitStat

    me:= false
    // other:= false
    if len(args)>1 {
        if args[1]=="me" {
            me = true
        }
        // else{
        //     //functionality to look up user by nickname
        // }
    }

    //this will give bit values
    if(me) {
        b := db.GetBitStats(user.ID)
        if b == nil {
            bits = BitStat{UserID: user.ID, BitValue: 0}
            db.SetBitStats(user.ID, bits.BitValue)
        } else {
            bits = *b
        }
    } else {
        bitslist = db.GetTopBitStats(10)
    }



    w := &tabwriter.Writer{}
    buf := &bytes.Buffer{}

    w.Init(buf, 0, 4, 0, ' ', 0)
    // fmt.Fprintf(w, "%s Game-Time Stats:\n", ) // Not sure how to get nicknames...
    fmt.Fprintf(w, "```\n")
    if me {
        fmt.Fprintf(w, "%s: \t %d bits\n", utilGetPreferredName(guild, bits.UserID), bits.BitValue)
    } else {
        for _, bit := range(bitslist) {
            name := utilGetPreferredName(guild, bit.UserID)
            fmt.Fprintf(w, "%s: \t %d bits\n", name, bit.BitValue)
        }
    }

    fmt.Fprintf(w, "```\n")
    w.Flush()
    return buf.String()
}

func giveWeeklyBitBonus(guild *discordgo.Guild, userID string) string {
    db := dbGetSession(guild.ID)
    mybits := db.GetBitStats(userID)
    if mybits == nil {
        db.SetBitStats(userID, 50)
        mybits = db.GetBitStats(userID)
    } else {
        db.IncBitStats(userID, 50)
        mybits = db.GetBitStats(userID)
    }
    w := &tabwriter.Writer{}
    buf := &bytes.Buffer{}
    log.Info("Giving Weekly Bonus")
    w.Init(buf, 0, 4, 0, ' ', 0)
    fmt.Fprintf(w, "Welcome **%s**, you get 50 bits for joining this week!\n You now have **%d bits**\n", utilGetPreferredName(guild, userID), mybits.BitValue)
    w.Flush()
    return buf.String()
}
