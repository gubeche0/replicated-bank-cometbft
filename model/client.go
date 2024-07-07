package model

import (
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger/v3"
)

type Client struct {
	Name string `json:"name" badgerhold:"index"`
	// PubKey  ed25519.PubKey `badgerhold:"index"` // this is just a wrapper around bytes
	Balance int64 `json:"balance"`
}

func getKeyName(name string) []byte {
	return []byte(fmt.Sprintf("client_%s", name))
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
