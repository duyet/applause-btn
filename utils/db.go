package utils

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"

	"github.com/dgraph-io/badger/v2"
)

// GetDB open connection to badger db
func GetDB() *badger.DB {
	badgerLocation := "/tmp/badger"
	val, ok := os.LookupEnv("DB_LOCATION")
	if ok {
		badgerLocation = val
	}

	db, err := badger.Open(badger.DefaultOptions(badgerLocation))
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// GetItem get Item from DB
func GetItem(sourceURL string) (ret Item, err error) {
	db := GetDB()
	defer db.Close()

	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(sourceURL))
		if err != nil {
			return err
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		d := gob.NewDecoder(bytes.NewReader(val))
		if err := d.Decode(&ret); err != nil {
			panic(err)
		}

		return nil
	})

	return
}

// GetItems get items from urls
func GetItems(listURL []string) (ret []Item, err error) {
	for _, url := range listURL {
		item, err := GetItem(url)
		if err != nil {
			ret = append(ret, Item{})
		} else {
			ret = append(ret, item)
		}
	}

	return
}

// PutItem put Item to DB
func PutItem(sourceURL string, item Item) error {
	db := GetDB()
	defer db.Close()

	var b bytes.Buffer
	e := gob.NewEncoder(&b)
	if err := e.Encode(item); err != nil {
		panic(err)
	}

	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(sourceURL), b.Bytes())
		return err
	})

	return err
}
