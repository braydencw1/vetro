package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var morthisId string = "186317976033558528"
var vetroId string = "1131832403581747381"
var channelId string = "209404729225248769"
var griefers []string = []string{}

func main() {
	// Load environment variables from the .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	token := os.Getenv("TOKEN")

	// Initialize Discord session
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session:", err)
	}
	discord.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	// Open a connection to Discord
	if err := discord.Open(); err != nil {
		log.Fatal("Error opening Discord connection:", err)
	}
	defer discord.Close()
	discord.AddHandler(onReady)
	// Register the messageCreate functfunc messageCreate(s *discordgo.Session, m *discordgo.MessageCreate)ion as a callback for the MessageCreate event
	// discord.AddMessageCreateHandler(messageCreate)
	discord.AddHandler(handleMessageEvents)
	discord.AddHandler(handleVoiceStateUpdate)
	// Keep the bot running
	log.Println("Bot is now running. Press Ctrl+C to exit.")
	select {} // Block the main goroutine indefinitely
}

func onReady(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Logged in as %s\n", event.User.String())
}

func handleMessageEvents(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	addGriefer(s, m)
	handleCommands(s, m)
}

func addGriefer(s *discordgo.Session, m *discordgo.MessageCreate) {
	parts := strings.Split(m.Content, " ")

	if parts[0] == "!grief" {
		if len(m.Mentions) == 0 {
			if len(griefers) == 0 {
				s.ChannelMessageSend(m.ChannelID, "Nobody is getting griefed!")

				return
			} else {
				myGriefees := []string{}

				for _, grief := range griefers {
					myGriefees = append(myGriefees, fmt.Sprintf("<@%s>", grief))
				}
				s.ChannelMessageSend(m.ChannelID, strings.Join(myGriefees, " "))

				return
			}
		}

		for _, mention := range m.Mentions {
			griefers = append(griefers, mention.ID)
		}

		s.ChannelMessageSend(m.ChannelID, "This brotha is getting griefed")

		return
	}
}

func handleCommands(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != s.State.User.ID {
		log.Printf(m.Author.Username + ": " + m.Content)
	}

	// Respond to messages
	if m.Content == "!hello" {
		// Reply with a message
		s.ChannelMessageSend(m.ChannelID, "Hello, "+m.Author.Username+"!")
	}

	if m.Content == fmt.Sprintf("<@%v>", vetroId) {
		s.ChannelMessageSend(m.ChannelID, "Hey, "+m.Author.Username+"!")

	}
	if m.Content == "https://imgur.com/a/XQ3pPTQ" {
		s.ChannelMessageSend(m.ChannelID, "Assemble!!!!!")
	}
}

func handleVoiceStateUpdate(s *discordgo.Session, m *discordgo.VoiceStateUpdate) {
	if m.ChannelID != channelId {
		return
	}

	if m.VoiceState.UserID == morthisId {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hello gaylord <@%v>", morthisId))
	}

	for _, griefee := range griefers {
		if m.VoiceState.UserID == griefee {
			s.GuildMemberMove(m.GuildID, griefee, &channelId)
		}
	}
}
