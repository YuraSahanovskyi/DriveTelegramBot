package gdrive

import (
	"context"
	"io"
	"net/http"
	"strconv"

	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type Drive struct {
	config *oauth2.Config
	ctx    context.Context
}

func NewDrive() (*Drive, error) {
	config, err := Init()
	return &Drive{config: config, ctx: context.Background()}, err
}

// returns http client for oauth token
func (d *Drive) GetClient(token *oauth2.Token) *http.Client {
	return d.config.Client(d.ctx, token)
}

// exchanges code for oauth token
func (d *Drive) ExchangeCode(code string) (*oauth2.Token, error) {
	return d.config.Exchange(d.ctx, code)
}

// returns authorization url for user
func (d *Drive) GetAuthUrl(userID int64) string {
	return d.config.AuthCodeURL(strconv.FormatInt(userID, 10), oauth2.AccessTypeOnline)
}

func (d *Drive) UploadFile(token *oauth2.Token, file *drive.File, r io.Reader) error {
	client := d.GetClient(token)
	//create new service
	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	_, err = srv.Files.Create(file).Media(r).Do()
	return err
}
