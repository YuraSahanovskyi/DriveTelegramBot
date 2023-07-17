package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const commandStart = "start"

func (b *Bot) handleMessage(msg tgbotapi.Message) {
	if msg.IsCommand() {
		b.handleCommand(msg)
	} else {
		b.handleFile(msg)
	}
}

func (b *Bot) handleCommand(msg tgbotapi.Message) {
	switch msg.Command() {
	case commandStart:
		b.handleStart(msg.Chat.ID)
	default:
		msg := tgbotapi.NewMessage(msg.Chat.ID, "I don't understand that command")
		b.tg.Send(msg)
	}
}

func (b *Bot) handleStart(id int64) {
	response := fmt.Sprintf("follow auth link \n%v", b.drive.GetAuthUrl(id))
	msg := tgbotapi.NewMessage(id, response)
	b.tg.Send(msg)
}

func (b *Bot) handleFile(msg tgbotapi.Message) {
	//TODO: implement
}
