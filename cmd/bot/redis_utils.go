package main

import (
    "bytes"
    "fmt"
    "runtime"
    "strconv"
    "text/tabwriter"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/bwmarrin/discordgo"
    "github.com/dustin/go-humanize"
    redis "gopkg.in/redis.v3"
)

func rdCalculateAirhornsPerSecond(cid string) {
    current, _ := strconv.Atoi(rcli.Get("airhorn:a:total").Val())
    time.Sleep(time.Second * 10)
    latest, _ := strconv.Atoi(rcli.Get("airhorn:a:total").Val())

    discord.ChannelMessageSend(cid, fmt.Sprintf("Current APS: %v", (float64(latest-current))/10.0))
}

func rdDisplayBotStats(cid string) {
    stats := runtime.MemStats{}
    runtime.ReadMemStats(&stats)

    users := 0
    for _, guild := range discord.State.Ready.Guilds {
        users += len(guild.Members)
    }

    w := &tabwriter.Writer{}
    buf := &bytes.Buffer{}

    w.Init(buf, 0, 4, 0, ' ', 0)
    fmt.Fprintf(w, "```\n")
    fmt.Fprintf(w, "Discordgo: \t%s\n", discordgo.VERSION)
    fmt.Fprintf(w, "Go: \t%s\n", runtime.Version())
    fmt.Fprintf(w, "Memory: \t%s / %s (%s total allocated)\n", humanize.Bytes(stats.Alloc), humanize.Bytes(stats.Sys), humanize.Bytes(stats.TotalAlloc))
    fmt.Fprintf(w, "Tasks: \t%d\n", runtime.NumGoroutine())
    fmt.Fprintf(w, "Servers: \t%d\n", len(discord.State.Ready.Guilds))
    fmt.Fprintf(w, "Users: \t%d\n", users)
    fmt.Fprintf(w, "```\n")
    w.Flush()
    discord.ChannelMessageSend(cid, buf.String())
}

func rdUtilSumRedisKeys(keys []string) int {
    results := make([]*redis.StringCmd, 0)

    rcli.Pipelined(func(pipe *redis.Pipeline) error {
        for _, key := range keys {
            results = append(results, pipe.Get(key))
        }
        return nil
    })

    var total int
    for _, i := range results {
        t, _ := strconv.Atoi(i.Val())
        total += t
    }

    return total
}

func rdDisplayUserStats(cid, uid string) {
    keys, err := rcli.Keys(fmt.Sprintf("airhorn:*:user:%s:sound:*", uid)).Result()
    if err != nil {
        return
    }

    totalAirhorns := rdUtilSumRedisKeys(keys)
    discord.ChannelMessageSend(cid, fmt.Sprintf("Total Airhorns: %v", totalAirhorns))
}

func rdDisplayServerStats(cid, sid string) {
    keys, err := rcli.Keys(fmt.Sprintf("airhorn:*:guild:%s:sound:*", sid)).Result()
    if err != nil {
        return
    }

    totalAirhorns := rdUtilSumRedisKeys(keys)
    discord.ChannelMessageSend(cid, fmt.Sprintf("Total Airhorns: %v", totalAirhorns))
}

func rdTrackSoundStats(play *Play) {
    if rcli == nil {
        return
    }

    _, err := rcli.Pipelined(func(pipe *redis.Pipeline) error {
        var baseChar string

        if play.Forced {
            baseChar = "f"
        } else {
            baseChar = "a"
        }

        base := fmt.Sprintf("airhorn:%s", baseChar)
        pipe.Incr("airhorn:total")
        pipe.Incr(fmt.Sprintf("%s:total", base))
        pipe.Incr(fmt.Sprintf("%s:sound:%s", base, play.Sound.Name))
        pipe.Incr(fmt.Sprintf("%s:user:%s:sound:%s", base, play.UserID, play.Sound.Name))
        pipe.Incr(fmt.Sprintf("%s:guild:%s:sound:%s", base, play.GuildID, play.Sound.Name))
        pipe.Incr(fmt.Sprintf("%s:guild:%s:chan:%s:sound:%s", base, play.GuildID, play.ChannelID, play.Sound.Name))
        pipe.SAdd(fmt.Sprintf("%s:users", base), play.UserID)
        pipe.SAdd(fmt.Sprintf("%s:guilds", base), play.GuildID)
        pipe.SAdd(fmt.Sprintf("%s:channels", base), play.ChannelID)
        return nil
    })

    if err != nil {
        log.WithFields(log.Fields{
            "error": err,
        }).Warning("Failed to track stats in redis")
    }
}