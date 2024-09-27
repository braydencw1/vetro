package bot

import (
	"fmt"
	"log"
	"strings"
	"time"
	"venova/db"

	"github.com/bwmarrin/discordgo"
)

var tcGeneralId string = "209403061205073931"

var morthisId string = "186317976033558528"
var bettyId string = "641009995634180096"
var venovaId string = "1163950982259036302"
var blueId string = "202213189482446851"
var bangersRoleId string = "1079585245575270480"

var mcRoleId string = "1183228947874459668"
var frostedRoleId string = "618635064451923979"
var channelId string = "209403061205073931"
var griefers []string = []string{}
var joinableRolesList = []string{
	"bangers",
	"dorklock",
	"apes",
}

var joinableRolesMap = map[string]string{
	"apes":     "1250598584534175784",
	"dorklock": "1282817878244200488",
	"bangers":  "1079585245575270480",
}

func OnReady(discord *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Logged in as %s\n", event.User.String())
}

func HandleMessageEvents(discord *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID == discord.State.User.ID {
		return
	}

	if msg.Content == fmt.Sprintf("<@%v>", venovaId) {
		discord.ChannelMessageSend(msg.ChannelID, strings.ReplaceAll(db.DndMsgResponse(), "{nick}", msg.Author.Username))
	} else if msg.Content == fmt.Sprintf("<@&%s>", bangersRoleId) {
		discord.ChannelMessageSend(msg.ChannelID, "https://imgur.com/K7lTDGU")
	}

	AddGriefer(discord, msg)
	HandleCommands(discord, msg)
}

func GetUsernameFromID(session *discordgo.Session, userID string) (string, error) {
	user, err := session.User(userID)
	if err != nil {
		return "", err
	}
	return user.Username, nil
}

func AddGriefer(discord *discordgo.Session, msg *discordgo.MessageCreate) {
	parts := strings.Split(msg.Content, " ")

	if parts[0] == "!grief" {
		if len(msg.Mentions) == 0 {
			if len(griefers) == 0 {
				discord.ChannelMessageSend(msg.ChannelID, "Nobody is getting griefed!")

				return
			} else {
				myGriefees := []string{}

				for _, grief := range griefers {
					myGriefees = append(myGriefees, fmt.Sprintf("<@%s>", grief))
				}
				discord.ChannelMessageSend(msg.ChannelID, strings.Join(myGriefees, " "))

				return
			}
		}

		for _, mention := range msg.Mentions {
			griefers = append(griefers, mention.ID)
		}

		discord.ChannelMessageSend(msg.ChannelID, "This brotha is getting griefed")

		return
	}
}

func HandleVoiceStateUpdate(discord *discordgo.Session, msg *discordgo.VoiceStateUpdate) {
	if msg.ChannelID != channelId {
		return
	}

	if msg.VoiceState.UserID == morthisId {
		discord.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Hello cutie <@%s>", morthisId))
	}

	for _, griefee := range griefers {
		if msg.VoiceState.UserID == griefee {
			discord.GuildMemberMove(msg.GuildID, griefee, &channelId)
		}
	}
}

func PlayDateCheckRoutine(discord *discordgo.Session) {
	playDateCheck(discord)
	timer := time.NewTicker(24 * time.Hour)
	for range timer.C {
		playDateCheck(discord)
	}
}

func playDateCheck(discord *discordgo.Session) {
	nextDay := time.Now().Add(24 * time.Hour)
	res, tcId, roleId, err := db.GetPlayDates(nextDay)
	if err != nil {
		log.Printf("Failed to get play dates: %v", err)
		return
	}

	msg := fmt.Sprintf("Dnd is scheduled for tomorrow <@&%v>", roleId)
	if res {
		discord.ChannelMessageSend(fmt.Sprintf("%v", tcId), msg)
	}
}

func createTimer(timeLength string) (time.Time, error) {
	duration, err := time.ParseDuration(timeLength)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return time.Time{}, err
	}
	timer := time.Now().Add(duration)
	return timer, nil
}

func TimerCheckerRoutine(discord *discordgo.Session, timer time.Time, UserID string, errChan chan error) {
	ticker := time.NewTicker(1 * time.Minute) // Ticker to check every minute
	defer ticker.Stop()
	defer close(errChan)
	for {
		<-ticker.C
		if time.Now().After(timer) {
			dmChannel, err := discord.UserChannelCreate(UserID)
			if err != nil {
				errChan <- fmt.Errorf("error creating dm channel : %w", err)
				return
			}
			_, err = discord.ChannelMessageSend(dmChannel.ID, "Your timer is up!")
			if err != nil {
				errChan <- fmt.Errorf("error sending dm for time : %w", err)
			}
			return
		}
	}
}
