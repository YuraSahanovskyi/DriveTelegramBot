package gdrive

import (
	"errors"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

func Init() (*oauth2.Config, error) {
	clientID := viper.GetString("CLIENT_ID")
	if clientID == "" {
		return nil, errors.New("can't get client ID")
	}
	clientSecret := viper.GetString("CLIENT_SECRET")
	if clientSecret == "" {
		return nil, errors.New("can't get client secret")
	}
	redirectURL := viper.GetString("REDIRECT_URL")
	if redirectURL == "" {
		return nil, errors.New("can't get redirect URL")
	}
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{drive.DriveScope},
		Endpoint:     google.Endpoint,
	}, nil
}
