package telegram

import (
	"github.com/spf13/viper"
)

type Messages struct {
	Auth    AuthMessages
	File    FileMessages
	Command CommandMessages
}

type AuthMessages struct {
	NoTokenInDB          string `mapstructure:"no_token_in_db"`
	NoCodeInDB           string `mapstructure:"no_code_in_db"`
	AlreadyAuthenticated string `mapstructure:"already_authenticated"`
	AuthLink             string `mapstructure:"auth_link"`
	CannotExchangeToken  string `mapstructure:"cannot_exchange_token"`
	CannotMarshalToken   string `mapstructure:"cannot_marshal_token"`
	CannotSaveToken      string `mapstructure:"cannot_save_token"`
	SuccessfulIn         string `mapstructure:"successful_in"`
	SuccessfulOut        string `mapstructure:"successful_out"`
	SomethingWrong       string `mapstructure:"something_wrong"`
	NotAuthenticated     string `mapstructure:"not_authenticated"`
}

type FileMessages struct {
	NoFilesToSave string `mapstructure:"no_files_to_save"`
	FileError     string `mapstructure:"file_error"`
	FileSuccess   string `mapstructure:"file_success"`
}

type CommandMessages struct {
	Unknown string `mapstructure:"unknown"`
}

func InitMessages() (*Messages, error) {
	var messages Messages
	if err := viper.Unmarshal(&messages); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("messages.auth", &messages.Auth); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("messages.file", &messages.File); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("messages.command", &messages.Command); err != nil {
		return nil, err
	}
	return &messages, nil
}
