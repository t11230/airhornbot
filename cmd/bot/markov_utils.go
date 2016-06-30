package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
)

var (
	// Markov chains map
	chains map[string] *GuildChain = make(map[string] *GuildChain)
)

// Tracks which Guild this chain belongs to
type GuildChain struct {
	GuildID				string
	MarkovChain 		*Chain
}

func mkGetMessage(guild *discordgo.Guild, user *discordgo.User) string {
	// Check if we already have a connection to this guild
	guildChain, exists := chains[guild.ID]
	if !exists {
		log.Printf("No chain for guild %d", guild.ID)
		return ""
	}

	if len(guildChain.MarkovChain.Chain) == 0 {
		log.Printf("Empty markov chain...")
		return ""
	}
	return guildChain.MarkovChain.Generate(40)
}

func mkWriteMessage(guild *discordgo.Guild, content string) {
	guildChain, exists := chains[guild.ID]
	if !exists {
		log.Printf("No chain for guild %d", guild.ID)
		return
	}

	guildChain.MarkovChain.Write(utilParseText(content))
	log.Printf("Chain length: %d", len(guildChain.MarkovChain.Chain))
}

func mkLoadChain(cid string, guild *discordgo.Guild, user *discordgo.User, markovStatePath string) {
	chain := NewChain(2)

	// Rebuild the markov chain from state
	log.Printf("Loading state from '%s'.", markovStatePath)
	err := chain.Load(markovStatePath)
	if err != nil {
		//log.Fatal(err)
		log.Printf("Could not load from '%s'. This may be expected.", markovStatePath)
	} else {
		log.Printf("Loaded previous state from '%s' (%d suffixes).", markovStatePath, len(chain.Chain))
	}

	chains[guild.ID] = &GuildChain {
		GuildID: guild.ID,
		MarkovChain: chain,
	}
}

func mkGenerateChain(cid string, guild *discordgo.Guild, user *discordgo.User,
							chatLogPath string, markovStatePath string) {
	file, err := os.Open(chatLogPath)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    log.Printf("Starting import from %s", file)

    chain := NewChain(2)

	scanner := bufio.NewScanner(file)
	lastText := ""

    for scanner.Scan() {
        text := scanner.Text()
        text = utilParseText(text)
        if text == "" {
        	continue
        }

        if len(lastText) != 0 {
        	chain.Write(text)

        	if rand.Intn(100) < 50 {
				text = lastText + " " + utilLowerFirst(text)
			} else {
				text = lastText + ". " + text
			}
			
			lastText = ""
		}

		if rand.Intn(100) < 25 {
			lastText = text
		}

       	chain.Write(text)
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

    // Write the state file
	err = chain.Save(markovStatePath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Import complete. %d suffixes", len(chain.Chain))

	chains[guild.ID] = &GuildChain {
		GuildID: guild.ID,
		MarkovChain: chain,
	}
}

func mkFetchAndSaveMessages(cid string, guild *discordgo.Guild, user *discordgo.User, cs string,
							chatLogPath string) {
	count, _ := strconv.Atoi(cs)

	f, _ := os.Create(chatLogPath)
	defer f.Close()

	fetchAll := false
	if count == 0 {
		fetchAll = true
	}

	lastID := ""
	total := 0
	for count > 0 || fetchAll {
		toGet := count
		if count > 100 || fetchAll {
			toGet = 100
		}
		m, err := discord.ChannelMessages(cid, toGet, lastID, "")

		if err != nil || len(m) == 0 {
			break
		}

		for i := 1; i < len(m); i++ {
			fmt.Fprintf(f, "%s\n", m[i].Content)
		}
		lastID = m[len(m) - 1].ID
		count -= toGet
		total += toGet
	}
	
	f.Sync()

	discord.ChannelMessageSend(cid, fmt.Sprintf("Saved something around %d messages\n", total))
}