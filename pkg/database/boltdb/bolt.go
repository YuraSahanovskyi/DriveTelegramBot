package boltdb

import (
	"errors"
	"strconv"

	"github.com/boltdb/bolt"
)

const bucketName = "users"

type BoltRepository struct {
	db *bolt.DB
}

func (br *BoltRepository) Get(id int64) (string, error) {
	var value string
	err := br.db.View(func(tx *bolt.Tx) error {
		//get bucket or create it
		b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		//get value by id
		byteValue := b.Get(int64ToBytes(id))
		if byteValue == nil {
			return errors.New("no such value")
		}
		//convert value to string
		value = string(byteValue)
		return nil
	})
	if err != nil {
		return "", err
	}
	return value, nil
}

func (br *BoltRepository) Put(id int64, value string) error {
	return br.db.Update(func(tx *bolt.Tx) error {
		//get bucket or create it
		b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		//put value by id
		return b.Put(int64ToBytes(id), []byte(value))
	})
}

func int64ToBytes(id int64) []byte {
	return []byte(strconv.FormatInt(id, 10))
}
