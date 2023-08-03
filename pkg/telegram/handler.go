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
		msg := tgbotapi.NewMessage(msg.Chat.ID, b.msgs.Command.Unknown)
		b.tg.Send(msg)
	}
}

func (b *Bot) handleStart(id int64) {
	//check if user is authenticated
	if b.checkAuth(id) {
		msg := tgbotapi.NewMessage(id, b.msgs.Auth.AlreadyAuthenticated)
		b.tg.Send(msg)
		return
	}
	//send auth link
	response := fmt.Sprintf(b.msgs.Auth.AuthLink, b.drive.GetAuthUrl(id))
	msg := tgbotapi.NewMessage(id, response)
	b.tg.Send(msg)
}

func (b *Bot) handleConfirm(id int64) {
	//check if user is authenticated
	if b.checkAuth(id) {
		msg := tgbotapi.NewMessage(id, b.msgs.Auth.AlreadyAuthenticated)
		b.tg.Send(msg)
		return
	}
	//get code from database
	code, err := b.repo.Get(database.Code, id)
	if err != nil {
		msg := tgbotapi.NewMessage(id, b.msgs.Auth.NoCodeInDB)
		b.tg.Send(msg)
		return
	}
	//exchange code for token
	token, err := b.drive.ExchangeCode(string(code))
	if err != nil {
		msg := tgbotapi.NewMessage(id, b.msgs.Auth.CannotExchangeToken)
		b.tg.Send(msg)
		return
	}
	//marshal token to []byte
	tokenBytes, err := json.Marshal(token)
	if err != nil {
		msg := tgbotapi.NewMessage(id, b.msgs.Auth.CannotMarshalToken)
		b.tg.Send(msg)
		return
	}
	//save token to database
	if err := b.repo.Put(database.Token, id, tokenBytes); err != nil {
		msg := tgbotapi.NewMessage(id, b.msgs.Auth.CannotSaveToken)
		b.tg.Send(msg)
		return
	}
	msg := tgbotapi.NewMessage(id, b.msgs.Auth.SuccessfulIn)
	b.tg.Send(msg)
	//delete code from database
	b.repo.Delete(database.Code, id)
}

func (b *Bot) handleLogOut(id int64) {
	if err := b.repo.Delete(database.Token, id); err != nil {
		msg := tgbotapi.NewMessage(id, b.msgs.Auth.SomethingWrong)
		b.tg.Send(msg)
		return
	}
	msg := tgbotapi.NewMessage(id, b.msgs.Auth.SuccessfulOut)
	b.tg.Send(msg)
}

func (b *Bot) handleFile(msg tgbotapi.Message) {
	//check if user if authenticated
	if !b.checkAuth(msg.Chat.ID) {
		msg := tgbotapi.NewMessage(msg.Chat.ID, b.msgs.Auth.NotAuthenticated)
		b.tg.Send(msg)
		return
	}
	//get oauth token from database
	token := b.getAuth(msg.Chat.ID)
	//get files info
	filesInfo := b.getFileInfo(msg)
	var text string
	if filesInfo.IsVoid() {
		text = b.msgs.File.NoFilesToSave
	} else {
		fileName, err := b.uploadFile(filesInfo, token)
		displayName := getDisplayName(fileName)
		if err != nil {
			text = fmt.Sprintf(b.msgs.File.FileError, displayName, err)
		} else {
			text = fmt.Sprintf(b.msgs.File.FileSuccess, displayName)
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
