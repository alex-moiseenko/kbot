/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	telebot "gopkg.in/telebot.v3"
)

var (
	// TeleToken bot
	TeleToken = os.Getenv("TELE_TOKEN")
)

type WeatherResponse struct {
	CurrentWeather struct {
		Temperature float32 `json:"temperature"`
		Windspeed   float32 `json:"windspeed"`
	} `json:"current_weather"`
}

func getWeather(lat, lon float32) (WeatherResponse, error) {
	var weather WeatherResponse

	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current_weather=true", lat, lon)
	resp, err := http.Get(url)
	if err != nil {
		return weather, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return weather, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&weather)
	if err != nil {
		return weather, err
	}

	return weather, nil
}

// kbotCmd represents the kbot command
var kbotCmd = &cobra.Command{
	Use:     "kbot",
	Aliases: []string{"start"},
	Short:   "Starts telegram bot",
	Long:    `Starts telegram bot server which can respond on commands`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("kbot %s started", appVersion)
		kbot, err := telebot.NewBot(telebot.Settings{
			URL:    "",
			Token:  TeleToken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})

		if err != nil {
			log.Fatalf("Please check TELE_TOKEN env variable. %s", err)
			return
		}

		kbot.Handle("/start", func(m telebot.Context) error {
			sender := m.Sender().FirstName + m.Sender().LastName
			return m.Send(fmt.Sprintf("Hello! %s!\n I'm Kbot %s. \n\n📍 You can send me your location and I will tell you the weather ⛅.\n🐱 Also i can send you a cat photo! Jut send me command /cat", sender, appVersion))
		})

		kbot.Handle("/cat", func(m telebot.Context) error {
			catURL := fmt.Sprintf("https://cataas.com/cat?timestamp=%d", time.Now().UnixNano())
			photo := &telebot.Photo{File: telebot.FromURL(catURL)}

			return m.Send(photo)
		})

		kbot.Handle(telebot.OnLocation, func(m telebot.Context) error {
			location := m.Message().Location

			log.Printf("Received location: Latitude: %f, Longitude: %f", location.Lat, location.Lng)

			weather, err := getWeather(location.Lat, location.Lng)
			if err != nil {
				log.Printf("Error fetching weather: %s", err)
				return m.Send("Sorry, I couldn't fetch the weather data.")
			}

			weatherMessage := fmt.Sprintf("The current temperature is %.2f°C with wind speed %.2fkm/h.",
				weather.CurrentWeather.Temperature, weather.CurrentWeather.Windspeed)

			return m.Send(weatherMessage)
		})

		kbot.Handle(telebot.OnText, func(m telebot.Context) error {
			// log.Print(m.Message().Payload, m.Text())

			payload := m.Message().Payload
			text := m.Text()
			sender := m.Sender().FirstName + m.Sender().LastName

			log.Printf("Payload: %s", payload)
			log.Printf("Text: %s", text)
			log.Printf("Sender: %s", sender)

			switch payload {
			case "hello":
				err = m.Send(fmt.Sprintf("Hello I'm Kbot %s!", appVersion))
			}

			return err
		})

		kbot.Start()
	},
}

func init() {
	rootCmd.AddCommand(kbotCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kbotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kbotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
