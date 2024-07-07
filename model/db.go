package model

import (
	"encoding/json"

	"github.com/dgraph-io/badger/v3"
)

type DB struct {
	db *badger.DB
}

// func (db *DB) CreateUser(client *Client) error {
// 	// Check if the user already exists
// 	err := db.db.View(func(txn *badger.Txn) error {
// 		_, err := txn.Get([]byte(client.Name))
// 		return err
// 	})
// 	if err == nil {
// 		return errors.New("Client already exists")
// 	}

// 	// Save the user to the database
// 	err = db.db.Update(func(txn *badger.Txn) error {
// 		userBytes, err := json.Marshal(client)
// 		if err != nil {
// 			return errors.New("failed to marshal user to JSON")
// 		}
// 		err = txn.Set([]byte(client.Name), userBytes)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// 	return err
// }

func (db *DB) FindUserByName(name string) (*Client, error) {
	var client *Client
	err := db.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(getKeyName(name))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, client)
		})
	})
	return client, err
}
