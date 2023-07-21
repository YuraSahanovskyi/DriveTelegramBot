package telegram

import (
	"encoding/json"
	"fmt"

	"github.com/YuraSahanovskyi/DriveTelegramBot/pkg/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	commandStart   = "start"
	commandConfirm = "confirm"
	commandLogOut  = "logout"
)

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
	case commandConfirm:
		b.handleConfirm(msg.Chat.ID)
	case commandLogOut:
		b.handleLogOut(msg.Chat.ID)
	default:
		msg := tgbotapi.NewMessage(msg.Chat.ID, "I don't understand that command")
		b.tg.Send(msg)
	}
}

func (b *Bot) handleStart(id int64) {
	response := fmt.Sprintf("follow auth link \n%v \nand then call \"/confirm\" command", b.drive.GetAuthUrl(id))
	msg := tgbotapi.NewMessage(id, response)
	b.tg.Send(msg)
}

func (b *Bot) handleConfirm(id int64) {
	code, err := b.repo.Get(database.Code, id)
	if err != nil {
		//TODO: error
		msg := tgbotapi.NewMessage(id, fmt.Sprintf("read code error: %v", err))
		b.tg.Send(msg)
		return
	}
	token, err := b.drive.ExchangeCode(code)
	if err != nil {
		//TODO: error
		msg := tgbotapi.NewMessage(id, fmt.Sprintf("exchange token error: %v", err))
		b.tg.Send(msg)
		return
	}
	tok, err := json.Marshal(token)
	if err != nil {
		//TODO: error
		msg := tgbotapi.NewMessage(id, fmt.Sprintf("marshal token error: %v", err))
		b.tg.Send(msg)
		return
	}
	if err := b.repo.Put(database.Token, id, string(tok)); err != nil {
		//TODO: error
		msg := tgbotapi.NewMessage(id, fmt.Sprintf("write token error: %v", err))
		b.tg.Send(msg)
		return
	}
	msg := tgbotapi.NewMessage(id, "Authorization successful")
	b.tg.Send(msg)
}

func (b *Bot) handleLogOut(id int64) {
	if err := b.repo.Delete(database.Code, id); err != nil {
		//TODO: error
		msg := tgbotapi.NewMessage(id, fmt.Sprintf("error: %v", err))
		b.tg.Send(msg)
		return
	}
	if err := b.repo.Delete(database.Token, id); err != nil {
		//TODO: error
		msg := tgbotapi.NewMessage(id, fmt.Sprintf("error: %v", err))
		b.tg.Send(msg)
		return
	}
	msg := tgbotapi.NewMessage(id, "Logout successful")
	b.tg.Send(msg)
}

func (b *Bot) handleFile(msg tgbotapi.Message) {
	//TODO: implement
}
