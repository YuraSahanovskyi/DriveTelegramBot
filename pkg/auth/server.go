package auth

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/YuraSahanovskyi/DriveTelegramBot/pkg/database"
	"github.com/spf13/viper"
)

var bot_link string
var server_port string

type AuthServer struct {
	repo database.Repository
}

func NewAuthServer(repo database.Repository) *AuthServer {
	if err := loadConfig(); err != nil {
		log.Fatal(err)
	}
	return &AuthServer{repo: repo}
}

func (auth *AuthServer) Start() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		userIDParam := r.URL.Query().Get("state")
		if userIDParam == "" {
			log.Println("cannot read state")
			return
		}
		userID, err := strconv.ParseInt(userIDParam, 10, 64)
		if err != nil {
			log.Printf("cannot convert user ID: %v\n", userIDParam)
			return
		}
		code := r.URL.Query().Get("code")
		if code == "" {
			log.Println("cannot read code")
			return
		}
		if err := auth.repo.Put(database.Code, userID, code); err != nil {
			log.Printf("cannot save code for user ID %v\n", userID)
			return
		}
		log.Printf("code for user %d saved", userID)
		http.Redirect(w, r, bot_link, http.StatusFound)
	})
	log.Printf("auth server started at port %v", server_port)
	return http.ListenAndServe(server_port, nil)

}

func loadConfig() error {
	bot_link = viper.GetString("bot_link")
	if bot_link == "" {
		return errors.New("bot_link is not set")
	}
	server_port = viper.GetString("server_port")
	if server_port == "" {
		return errors.New("server_port is not set")
	}
	return nil
}
