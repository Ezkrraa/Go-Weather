package bot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// use geocoding api to find lat and lon for any place
// http://api.openweathermap.org/geo/1.0/direct?q={place}&limit=1&appid={API key}

// use current weather data api to find weather information for those lat and lon

const geoURL string = "http://api.openweathermap.org/geo/1.0/direct?q="
const weatherURL string = "https://api.openweathermap.org/data/2.5/weather?lat="

type WeatherData struct {
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp       float64 `json:"temp"`
		Feels_Like float64 `json:"feels_like"`
		Humidity   int     `json:"humidity"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Name string `json:"name"`
}

type GeoData []struct {
	Name    string  `json:"name"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country"`
	State   string  `json:"state"`
}

func getCurrentWeather(message string) *discordgo.MessageSend {

	// get a url for figuring out location
	message = strings.TrimSpace(message)
	location := message[5:]
	geoURL := fmt.Sprintf("%s%s&limit=1&appid=%s", geoURL, location, OpenWeatherToken)

	// fmt.Println(geoURL)
	// if strings.HasSuffix(geoURL, "\r") {
	geoURL = geoURL[:len(geoURL)-1]
	// }

	// make a client for getting location
	geoClient := http.Client{Timeout: 5 * time.Second}
	// error handling for response
	geoResponse, err := geoClient.Get(geoURL)
	if err != nil {
		fmt.Println(geoURL, err)
		return &discordgo.MessageSend{
			Content: "Sorry, there was an error trying to get the location",
		}
	}

	// load the response
	geoBody, _ := io.ReadAll(geoResponse.Body)
	defer geoResponse.Body.Close()

	// load the response into a struct
	var data GeoData
	json.Unmarshal([]byte(geoBody), &data)

	if len(data) == 0 {
		fmt.Println(geoURL, err)
		return &discordgo.MessageSend{
			Content: "Sorry, there was an error trying to get the location",
		}
	}
	fmt.Println(data)
	// grab latitude and longitude
	lat := fmt.Sprintf("%f", data[0].Lat)
	lon := fmt.Sprintf("%f", data[0].Lon)

	// use lat and lon to get weather data for that location
	weatherURL := fmt.Sprintf("%s%s&lon=%s&appid=%s", weatherURL, lat, lon, OpenWeatherToken)
	weatherURL = weatherURL[:len(weatherURL)-1]
	weatherURL = weatherURL + "&units=metric"
	fmt.Println(weatherURL)
	weatherClient := http.Client{Timeout: 5 * time.Second}
	weatherResponse, err := weatherClient.Get(weatherURL)

	if err != nil {
		fmt.Println(err)
		return &discordgo.MessageSend{
			Content: "Sorry, there was an error trying to get the weather",
		}
	}
	weatherBody, _ := io.ReadAll(weatherResponse.Body)
	defer weatherResponse.Body.Close()

	var data2 WeatherData
	json.Unmarshal([]byte(weatherBody), &data2)

	fmt.Println(data2)
	// fmt.Println(data2.Weather)
	// fmt.Println(data2.Name)

	city := data2.Name
	conditions := data2.Weather[0].Description
	temperature := strconv.FormatFloat(data2.Main.Temp, 'f', 2, 64)
	feels_like := strconv.FormatFloat(data2.Main.Feels_Like, 'f', 2, 64)
	humidity := strconv.Itoa(data2.Main.Humidity)
	wind := strconv.FormatFloat(data2.Wind.Speed, 'f', 2, 64)
	caser := cases.Title(language.English)

	embed := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{{
			Type:        discordgo.EmbedTypeRich,
			Title:       "Current Weather",
			Description: city,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Conditions",
					Value:  caser.String(conditions),
					Inline: true,
				},
				{
					Name:   "Temperature",
					Value:  temperature + "°C",
					Inline: true,
				},
				{
					Name:   "Feels Like",
					Value:  feels_like + "°C",
					Inline: true,
				},
				{
					Name:   "Humidity",
					Value:  humidity + "%",
					Inline: true,
				},
				{
					Name:   "Wind",
					Value:  wind + " km/h",
					Inline: true,
				},
			},
		},
		},
	}

	return embed
}
