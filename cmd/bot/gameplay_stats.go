package main

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/bwmarrin/discordgo"
)

func gpPrintStats(cid string, user *discordgo.User) {	
	games, times := dbGetGameStats(user.ID)

	if len(games) == 0 {
		discord.ChannelMessageSend(cid, "No stats. Git gud scrub.")
		return
	}

	w := &tabwriter.Writer{}
	buf := &bytes.Buffer{}

	w.Init(buf, 0, 4, 0, ' ', 0)
	// fmt.Fprintf(w, "%s Game-Time Stats:\n", ) // Not sure how to get nicknames...
	fmt.Fprintf(w, "```\n")

	for i, game := range(games) {
		fmt.Fprintf(w, "%s: \t %.2f Hours\n", game, float64(times[i]) / (60.0*60.0))
	}
	
	fmt.Fprintf(w, "```\n")
	w.Flush()
	discord.ChannelMessageSend(cid, buf.String())
}