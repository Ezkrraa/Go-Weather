// services used:
//
//	OpenWeather API (openweathermap.org)
//	Discord API (discord.gg)
//	DiscordGo package (github.com/bwmarrin/discordgo)
package main

import (
	"discord-weather-bot/bot"
	"log"
	"os"

	_ "github.com/jpfuentes2/go-env/autoload"
)

func main() {

	if _, err := os.Stat(".env"); err != nil {
		file, err2 := os.Create(".env")
		if err2 != nil {
			log.Fatal("Failed to make an empty file, try adding .env to source")
		}
		file.WriteString("export BOT_TOKEN=<token here>\nexport OPENWEATHER_TOKEN=<token here>")
		file.Close()
	}
	botToken, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		log.Fatal("Must set Discord token as env variable: BOT_TOKEN")
	}
	openWeatherToken, ok := os.LookupEnv("OPENWEATHER_TOKEN")
	if !ok {
		log.Fatal("Must set Discord token as env variable: BOT_TOKEN")
	}

	bot.BotToken = botToken
	bot.OpenWeatherToken = openWeatherToken
	bot.Run()
}
