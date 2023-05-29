package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	OpenWeatherToken string
	BotToken         string
)

func Run() {
	discord, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatal(err)
	}
	// add event handler for new message
	discord.AddHandler(newMessage)
	err2 := discord.Open()
	if err2 != nil {
		fmt.Println(err2)
	}
	defer discord.Close()

	fmt.Println("Bot running...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("\nClosing...")
}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	// fmt.Println(message.Author.Username, ": ", message.Content)
	if message.Author.ID == discord.State.User.ID {
		return
	}

	/*switch {
	case strings.Contains(message.Content, "weather"):
		discord.ChannelMessageSend(message.ChannelID, "I can help with that!")
	case strings.Contains(message.Content, "bot"):
		discord.ChannelMessageSend(message.ChannelID, "Hi there!")
	}*/
	if strings.HasPrefix(message.Content, "!weather") {
		currentWeather := getCurrentWeather(message.Content)
		discord.ChannelMessageSendComplex(message.ChannelID, currentWeather)
	}
}
