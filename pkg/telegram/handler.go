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
	//check if user is authenticated
	if b.checkAuth(id) {
		msg := tgbotapi.NewMessage(id, "You are already authenticated")
		b.tg.Send(msg)
		return
	}
	//send auth link
	response := fmt.Sprintf("follow auth link \n%v \nand then call \"/confirm\" command", b.drive.GetAuthUrl(id))
	msg := tgbotapi.NewMessage(id, response)
	b.tg.Send(msg)
}

func (b *Bot) handleConfirm(id int64) {
	//check if user is authenticated
	if b.checkAuth(id) {
		msg := tgbotapi.NewMessage(id, "You are already authenticated")
		b.tg.Send(msg)
		return
	}
	//get code from database
	code, err := b.repo.Get(database.Code, id)
	if err != nil {
		msg := tgbotapi.NewMessage(id, "cannot fin your code in database, try to auth again")
		b.tg.Send(msg)
		return
	}
	//exchange code for token
	token, err := b.drive.ExchangeCode(string(code))
	if err != nil {
		msg := tgbotapi.NewMessage(id, "oops, cannot exchange your token, try again later")
		b.tg.Send(msg)
		return
	}
	//marshal token to []byte
	tokenBytes, err := json.Marshal(token)
	if err != nil {
		msg := tgbotapi.NewMessage(id, "cannot marshal token, try again later")
		b.tg.Send(msg)
		return
	}
	//save token to database
	if err := b.repo.Put(database.Token, id, tokenBytes); err != nil {
		msg := tgbotapi.NewMessage(id, "cannot save token to database")
		b.tg.Send(msg)
		return
	}
	msg := tgbotapi.NewMessage(id, "Authentication successful")
	b.tg.Send(msg)
	//delete code from database
	b.repo.Delete(database.Code, id)
}

func (b *Bot) handleLogOut(id int64) {
	if err := b.repo.Delete(database.Token, id); err != nil {
		msg := tgbotapi.NewMessage(id, "something went wrong, try again")
		b.tg.Send(msg)
		return
	}
	msg := tgbotapi.NewMessage(id, "Logout successful")
	b.tg.Send(msg)
}

func (b *Bot) handleFile(msg tgbotapi.Message) {
	//check if user if authenticated
	if !b.checkAuth(msg.Chat.ID) {
		msg := tgbotapi.NewMessage(msg.Chat.ID, "You are not authenticated, tap /start")
		b.tg.Send(msg)
		return
	}
	//get oauth token from database
	token := b.getAuth(msg.Chat.ID)
	//get files info
	filesInfo := b.getFileInfo(msg)
	var text string
	if filesInfo.IsVoid() {
		text = "No files to save"
	} else {
		fileName, err := b.uploadFile(filesInfo, token)
		displayName := getDisplayName(fileName)
		if err != nil {
			text = fmt.Sprintf("File %s: %v", displayName, err)
		} else {
			text = fmt.Sprintf("File %s uploaded", displayName)
		}
	}
	response := tgbotapi.NewMessage(msg.Chat.ID, text)
	b.tg.Send(response)
}

func getDisplayName(name string) string {
	if len(name) > 10 {
		return name[:10] + "..."
	} else {
		return name
	}
}
