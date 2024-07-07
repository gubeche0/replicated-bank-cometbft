package model

import (
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger/v3"
)

type Client struct {
	Name    string `json:"name" badgerhold:"index"`
	Balance int64  `json:"balance"`
}

func getKeyName(name string) []byte {
	return []byte(fmt.Sprintf("client_%s", name))
}

func ListClients(db *badger.DB) ([]*Client, error) {
	clients := make([]*Client, 0)
	err := db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte("client_")
		for it.Rewind(); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				client := new(Client)
				if err := json.Unmarshal(val, &client); err != nil {
					return err
				}
				clients = append(clients, client)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return clients, err
}

func FindUserByName(txn *badger.Txn, name string) (*Client, error) {
	var client *Client
	item, err := txn.Get(getKeyName(name))
	if err != nil {
		return nil, err
	}

	err = item.Value(func(val []byte) error {
		return json.Unmarshal(val, &client)
	})
	return client, err
}
