package main

import (
	"log"

	"github.com/YuraSahanovskyi/DriveTelegramBot/pkg/auth"
	"github.com/YuraSahanovskyi/DriveTelegramBot/pkg/database/boltdb"
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

	viper.SetConfigName("main")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read main.yml config file: %v", err)
	}
}

func main() {
	initViper()

	//read telegram API key
	apiKey := viper.GetString("TELEGRAM_API_KEY")
	if apiKey == "" {
		log.Fatal("API key not found in environment variables")
	}

	//create telegram client
	tg, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	//create database
	db, err := boltdb.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	//create drive client
	drive, err := gdrive.NewDrive()
	if err != nil {
		log.Fatal("Drive don't created: ", err)
	}
	//create bot
	bot := telegram.NewBot(tg, drive, db)
	go bot.Start()
	//create auth server
	authServer := auth.NewAuthServer(db)
	if err := authServer.Start(); err != nil {
		log.Fatal(err)
	}
}
