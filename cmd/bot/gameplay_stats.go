package main

import (
    "bytes"
    "fmt"
    "text/tabwriter"

    "github.com/bwmarrin/discordgo"
)

type GameTrackEntry struct {
    UserID string
    Game string
    Time int
}

func gpPrintStats(guild *discordgo.Guild, user *discordgo.User, args []string) string { 
    db := dbGetSession(guild.ID)
    entries := db.GameTrackGetStats(user.ID, 10)

    if len(entries) == 0 {
        return "No stats. Git gud scrub."
    }

    w := &tabwriter.Writer{}
    buf := &bytes.Buffer{}

    w.Init(buf, 0, 4, 0, ' ', 0)
    // fmt.Fprintf(w, "%s Game-Time Stats:\n", ) // Not sure how to get nicknames...
    fmt.Fprintf(w, "```\n")

    for _, entry := range(entries) {
        fmt.Fprintf(w, "%s: \t %.2f Hours\n", entry.Game, float64(entry.Time) / (60.0*60.0))
    }
    
    fmt.Fprintf(w, "```\n")
    w.Flush()
    return buf.String()
}