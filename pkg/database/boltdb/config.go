package boltdb

import (
	"errors"

	"github.com/boltdb/bolt"
	"github.com/spf13/viper"
)

func InitDB() (*BoltRepository, error) {
	dbName := viper.GetString("db_name")
	if dbName == "" {
		return nil, errors.New("can't read database name")
	}

	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		return nil, err
	}
	return &BoltRepository{db}, nil
}
