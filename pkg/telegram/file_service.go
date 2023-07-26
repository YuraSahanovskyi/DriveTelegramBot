package telegram

import (
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
)

type FileInfo struct {
	FileID   string
	FileName string
	MimeType string
}

func (f *FileInfo) IsVoid() bool {
	return f.FileID == "" && f.FileName == "" && f.MimeType == ""
}

func (b *Bot) uploadFile(info FileInfo, token *oauth2.Token) (string, error) {
	tgFile, err := b.tg.GetFile(tgbotapi.FileConfig{FileID: info.FileID})
	if err != nil {
		return info.FileName, err
	}
	link := tgFile.Link(b.tg.Token)
	//download file from web
	resp, err := http.Get(link)
	if err != nil {
		return info.FileName, err
	}
	defer resp.Body.Close()
	file := &drive.File{
		Name:     info.FileName,
		MimeType: info.MimeType,
	}
	//upload file to drive
	if err := b.drive.UploadFile(token, file, resp.Body); err != nil {
		return info.FileName, err
	}
	return info.FileName, nil
}

func (b *Bot) getFileInfo(msg tgbotapi.Message) FileInfo {
	if msg.Document != nil {
		return FileInfo{
			FileID:   msg.Document.FileID,
			FileName: msg.Document.FileName,
			MimeType: msg.Document.MimeType,
		}
	} else if msg.Audio != nil {
		return FileInfo{
			FileID:   msg.Audio.FileID,
			FileName: msg.Audio.FileName,
			MimeType: msg.Audio.MimeType,
		}
	} else if msg.Animation != nil {
		return FileInfo{
			FileID:   msg.Animation.FileID,
			FileName: msg.Animation.FileName,
			MimeType: msg.Animation.MimeType,
		}
	} else if msg.Video != nil {
		return FileInfo{
			FileID:   msg.Video.FileID,
			FileName: msg.Video.FileName,
			MimeType: msg.Video.MimeType,
		}
	} else if msg.Photo != nil {
		photo := msg.Photo[len(msg.Photo)-1]
		return FileInfo{
			FileID:   photo.FileID,
			FileName: photo.FileID + ".jpg",
			MimeType: "image/jpeg",
		}
	}
	return FileInfo{}
}
