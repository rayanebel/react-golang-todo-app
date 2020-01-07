package middleware

import (
	"fmt"
	"log"
	"path"

	badger "github.com/dgraph-io/badger"
)

const (
	DBRootPath   = "/tmp"
	DBFilename   = "data.json"
	DBBucketName = "todolist"
)

var Database *badger.DB

func InitBadgerDB() error {
	fmt.Println("Init badger database...")
	var err error
	Database, err = badger.Open(badger.DefaultOptions(path.Join(DBRootPath, DBFilename)))
	if err != nil {
		log.Fatal(err)
	}
	key := []byte("tasks")
	err = Database.Update(func(txn *badger.Txn) error {
		// check if key exist if not create empty one.
		_, err := txn.Get(key)
		if err != nil {
			err = txn.Set(key, []byte(`{}`))
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
		}
		return err
	})

	fmt.Println("Database has been initialized")
	return err
}

func GetItemByKey(key []byte) ([]byte, error) {
	var item *badger.Item
	var itemValue []byte
	err := Database.View(func(tx *badger.Txn) error {
		var err error
		item, err = tx.Get(key)
		if err != nil {
			log.Fatal(err)
			return err
		}

		err = item.Value(func(val []byte) error {
			itemValue, err = item.ValueCopy(nil)
			if err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return itemValue, nil
}

func NewItemForKey(key, value []byte) error {
	err := Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, value)
		return err
	})
	return err
}

func DeleteItemByKey(key []byte) error {
	err := Database.Update(func(txn *badger.Txn) error {
		err := txn.Delete(key)
		return err
	})
	return err
}
