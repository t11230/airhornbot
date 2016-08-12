package main

import (
    // "bytes"
    // "fmt"
    // "text/tabwriter"

    "github.com/bwmarrin/discordgo"
)

type BitStat struct {
    UserID string
    BitValue int
}

func bitsPrintStats(guild *discordgo.Guild, user *discordgo.User, args []string) string {
    // //this will give bit values
    // //if(me)
    // //  bit := dbGetMyBitStats(user.ID)
    // //else
    // //  bits := dbGetBitLeaderboard()
    // me:= false
    // other:= false
    // if len(args)>1 {
    //     if args[1]=="me" {
    //         me:= true
    //     }
    //     // else{
    //     //     //functionality to look up user by nickname
    //     // }
    // }
    // bits := 0
    //
    //
    // w := &tabwriter.Writer{}
    // buf := &bytes.Buffer{}
    //
    // w.Init(buf, 0, 4, 0, ' ', 0)
    // // fmt.Fprintf(w, "%s Game-Time Stats:\n", ) // Not sure how to get nicknames...
    // fmt.Fprintf(w, "```\n")
    // if me {
    //     fmt.Fprintf(w, "%s: \t %d bits\n", user.Nick, bit.BitValue)
    // } else {
    //     for i, bit := range(bits) {
    //         name = utilGetPreferredName(guild, bit.UserID)
    //         fmt.Fprintf(w, "%s: \t %d bits\n", name, bit.BitValue)
    //     }
    // }
    //
    // fmt.Fprintf(w, "```\n")
    // w.Flush()
    // return buf.String()
    return ""
}
