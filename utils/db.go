package utils

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/dgraph-io/badger/v2"
)

// Database wraps badger.DB with our application methods
type Database struct {
	db *badger.DB
}

// NewDatabase creates a new database connection
func NewDatabase(dbLocation string) (*Database, error) {
	opts := badger.DefaultOptions(dbLocation)
	opts.Logger = nil // Disable badger's verbose logging

	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	log.Printf("Database opened successfully at %s", dbLocation)
	return &Database{db: db}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	if d.db == nil {
		return nil
	}
	log.Println("Closing database connection")
	return d.db.Close()
}

// GetItem get Item from DB
func (d *Database) GetItem(sourceURL string) (Item, error) {
	var ret Item

	err := d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(sourceURL))
		if err != nil {
			return err
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		decoder := gob.NewDecoder(bytes.NewReader(val))
		if err := decoder.Decode(&ret); err != nil {
			return fmt.Errorf("failed to decode item: %w", err)
		}

		return nil
	})

	return ret, err
}

// GetItems get items from urls
func (d *Database) GetItems(listURL []string) ([]Item, error) {
	ret := make([]Item, 0, len(listURL))

	for _, url := range listURL {
		item, err := d.GetItem(url)
		if err != nil {
			// Return empty item for non-existent URLs
			ret = append(ret, Item{})
		} else {
			ret = append(ret, item)
		}
	}

	return ret, nil
}

// PutItem put Item to DB
func (d *Database) PutItem(sourceURL string, item Item) error {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(item); err != nil {
		return fmt.Errorf("failed to encode item: %w", err)
	}

	err := d.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(sourceURL), buf.Bytes())
	})

	if err != nil {
		return fmt.Errorf("failed to put item: %w", err)
	}

	return nil
}

// RunGC runs the garbage collector on the database
func (d *Database) RunGC(discardRatio float64) error {
	return d.db.RunValueLogGC(discardRatio)
}

// GetRawDB returns the underlying badger.DB (for backward compatibility)
func (d *Database) GetRawDB() *badger.DB {
	return d.db
}

// --- Legacy global DB support for backward compatibility ---
// This will be deprecated once all code is refactored

// DB the global connection to database (deprecated)
var DB *badger.DB

// GetDB open connection to badger db (deprecated - use NewDatabase instead)
func GetDB() *badger.DB {
	if DB != nil {
		return DB
	}

	log.Println("Warning: GetDB() is deprecated, use NewDatabase() instead")
	return DB
}

// GetItem get Item from DB using global DB (deprecated)
func GetItem(sourceURL string) (Item, error) {
	if DB == nil {
		return Item{}, fmt.Errorf("database not initialized")
	}

	var ret Item
	err := DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(sourceURL))
		if err != nil {
			return err
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		decoder := gob.NewDecoder(bytes.NewReader(val))
		if err := decoder.Decode(&ret); err != nil {
			return fmt.Errorf("failed to decode item: %w", err)
		}

		return nil
	})

	return ret, err
}

// GetItems get items from urls using global DB (deprecated)
func GetItems(listURL []string) ([]Item, error) {
	ret := make([]Item, 0, len(listURL))

	for _, url := range listURL {
		item, err := GetItem(url)
		if err != nil {
			ret = append(ret, Item{})
		} else {
			ret = append(ret, item)
		}
	}

	return ret, nil
}

// PutItem put Item to DB using global DB (deprecated)
func PutItem(sourceURL string, item Item) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(item); err != nil {
		return fmt.Errorf("failed to encode item: %w", err)
	}

	err := DB.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(sourceURL), buf.Bytes())
	})

	if err != nil {
		return fmt.Errorf("failed to put item: %w", err)
	}

	return nil
}
