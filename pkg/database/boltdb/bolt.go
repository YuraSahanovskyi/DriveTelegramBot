package boltdb

import (
	"errors"
	"strconv"

	"github.com/YuraSahanovskyi/DriveTelegramBot/pkg/database"
	"github.com/boltdb/bolt"
)

type BoltRepository struct {
	db *bolt.DB
}

func (br *BoltRepository) Get(bucket database.Bucket, id int64) (string, error) {
	var value string
	err := br.db.View(func(tx *bolt.Tx) error {
		//get bucket or create it
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return errors.New("bucket not found")
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

func (br *BoltRepository) Put(bucket database.Bucket, id int64, value string) error {
	return br.db.Update(func(tx *bolt.Tx) error {
		//get bucket or create it
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		//put value by id
		return b.Put(int64ToBytes(id), []byte(value))
	})
}

func (br *BoltRepository) Delete(bucket database.Bucket, id int64) error {
	return br.db.Update(func(tx *bolt.Tx) error {
		//get bucket or create it
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		//delete value by id
		return b.Delete(int64ToBytes(id))
	})
}

func int64ToBytes(id int64) []byte {
	return []byte(strconv.FormatInt(id, 10))
}
