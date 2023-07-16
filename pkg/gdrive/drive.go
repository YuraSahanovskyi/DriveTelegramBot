package gdrive

import "golang.org/x/oauth2"

type Drive struct {
	config *oauth2.Config
}

func NewDrive() (*Drive, error) {
	config, err := Init()
	return &Drive{config: config}, err
}
