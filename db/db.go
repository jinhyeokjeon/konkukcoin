package db

import (
	"github.com/jinhyeokjeon/konkukcoin/utils"
	bolt "go.etcd.io/bbolt"
)

const (
	dbName       = "saved_blockchain"
	blocksBucket = "blocks"
	newestBucket = "newest"
	newestHash   = "newestHash"
)

var db *bolt.DB

func DB() *bolt.DB {
	if db == nil {
		dbPointer, err := bolt.Open(dbName, 0600, nil)
		db = dbPointer
		utils.HandleErr(err)
		err = db.Update(func(tx *bolt.Tx) error {
			_, err = tx.CreateBucketIfNotExists([]byte(newestBucket))
			_, err = tx.CreateBucketIfNotExists([]byte(blocksBucket))
			return err
		})
		utils.HandleErr(err)
	}
	return db
}

func CloseDB() {
	DB().Close()
}

func SaveBlock(hash string, data []byte) {
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte(hash), data)
		return err
	})
	utils.HandleErr(err)
}

func SaveNewestHash(hash []byte) {
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(newestBucket))
		err := bucket.Put([]byte(newestHash), hash)
		return err
	})
	utils.HandleErr(err)
}

func GetBlock(hash string) []byte {
	var data []byte
	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		data = bucket.Get([]byte(hash))
		return nil
	})
	return data
}

func GetNewestHash() []byte {
	var hash []byte
	DB().View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(newestBucket))
		hash = bucket.Get([]byte(newestHash))
		return nil
	})
	return hash
}

func EmptyBlocks() {
	DB().Update(func(t *bolt.Tx) error {
		utils.HandleErr(t.DeleteBucket([]byte(blocksBucket)))
		_, err := t.CreateBucket([]byte(blocksBucket))
		utils.HandleErr(err)
		return nil
	})
}
