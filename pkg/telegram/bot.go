package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	tg *tgbotapi.BotAPI
}

func NewBot(tg *tgbotapi.BotAPI) *Bot {
	return &Bot{tg: tg}
}

func (b *Bot) Start() {
	log.Printf("Authorized on account %s", b.tg.Self.UserName)

	updates := b.getUpdatesChan()

	b.handleUpdates(updates)
}

func (b *Bot) getUpdatesChan() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.tg.GetUpdatesChan(u)
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "lalala")
			msg.ReplyToMessageID = update.Message.MessageID

			b.tg.Send(msg)
		}
	}
}
