package main

import (
    "bytes"
    "fmt"
    "text/tabwriter"

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
    return ""
}
