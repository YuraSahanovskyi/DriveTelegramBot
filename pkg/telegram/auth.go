package telegram

import (
	"encoding/json"
	"log"

	"github.com/YuraSahanovskyi/DriveTelegramBot/pkg/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/oauth2"
)

func (b *Bot) checkAuth(id int64) bool {
	tokenBytes, err := b.repo.Get(database.Token, id)
	if err != nil {
		return false
	}
	var token oauth2.Token
	if err := json.Unmarshal(tokenBytes, &token); err != nil {
		log.Printf("Error unmarshalling token: %v", err)
		return false
	}
	if !token.Valid() {
		return false
	}
	return true
}

func (b *Bot) getAuth(id int64) *oauth2.Token {
	//get token from database
	tokenBytes, err := b.repo.Get(database.Token, id)
	if err != nil {
		msg := tgbotapi.NewMessage(id, b.msgs.Auth.NoTokenInDB)
		b.tg.Send(msg)
		return nil
	}
	//unmarshal token
	var token oauth2.Token
	if err := json.Unmarshal(tokenBytes, &token); err != nil {
		log.Printf("Error unmarshalling token: %v", err)
		msg := tgbotapi.NewMessage(id, b.msgs.Auth.NoTokenInDB)
		b.tg.Send(msg)
		return nil
	}
	return &token
}
