package auth

import (
	"log"
	"net/http"
	"strconv"

	"github.com/YuraSahanovskyi/DriveTelegramBot/pkg/database"
)

type AuthServer struct {
	repo database.Repository
}

func NewAuthServer(repo database.Repository) *AuthServer {
	return &AuthServer{repo: repo}
}

func (auth *AuthServer) Start() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		userIDParam := r.URL.Query().Get("state")
		if userIDParam == "" {
			//TODO: ?
			log.Println("cannot read state")
			return
		}
		userID, err := strconv.ParseInt(userIDParam, 10, 64)
		if err != nil {
			//TODO: ?
			log.Println("cannot convert user ID")
			return
		}
		code := r.URL.Query().Get("code")
		if code == "" {
			//TODO:?
			log.Println("cannot read code")
			return
		}
		if err := auth.repo.Put(database.Code, userID, code); err != nil {
			//TODO: ?
			log.Println("cannot save code")
			return
		}
		log.Printf("code for user %d saved", userID)
		//TODO: move link to config file
		http.Redirect(w, r, "https://t.me/gdriveclientbot", http.StatusFound)
	})
	return http.ListenAndServe(":8080", nil)
}
