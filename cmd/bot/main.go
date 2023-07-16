package main

import (
	"log"

	"github.com/YuraSahanovskyi/DriveTelegramBot/pkg/gdrive"
	"github.com/YuraSahanovskyi/DriveTelegramBot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	env "github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func initViper() {
	if err := env.Load(); err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}
	viper.AutomaticEnv()
}

func main() {
	initViper()

	apiKey := viper.GetString("TELEGRAM_API_KEY")
	if apiKey == "" {
		log.Fatal("API key not found in environment variables")
	}

	tg, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		log.Fatal(err)
	}
	drive, err := gdrive.NewDrive()
	if err != nil {
		log.Fatal("Drive don't created: ", err)
	}
	log.Println(drive)
	bot := telegram.NewBot(tg)
	bot.Start()
}
