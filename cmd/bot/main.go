package main

import (
	"log"

	"github.com/YuraSahanovskyi/DriveTelegramBot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()
	apiKey := viper.GetString("TELEGRAM_API_KEY")
	if apiKey == "" {
		log.Fatal("API key not found in environment variables")
	}

	tg, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	bot := telegram.NewBot(tg)
	bot.Start()
}
