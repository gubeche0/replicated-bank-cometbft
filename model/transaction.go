package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/dgraph-io/badger/v3"
)

type Transaction struct {
	ID     int64  `json:"id" badgerhold:"index"`
	Type   string `json:"type"`
	From   string `json:"from" badgerhold:"index"`
	To     string `json:"to,omitempty" badgerhold:"index"`
	Amount int64  `json:"amount"`
}

func (txn *Transaction) ValidateBasic() error {
	if txn.ID <= 0 {
		return errors.New("id must be positive non-zero")
	}

	if txn.From == "" {
		return errors.New("from account must not be empty")
	}

	if txn.Amount < 0 {
		return errors.New("amount must be positive")
	}

	if txn.Type != "deposit" && txn.Type != "withdraw" && txn.Type != "transfer" {
		return errors.New("type must be deposit, withdraw, or transfer")
	}

	if txn.Type != "deposit" {
		if txn.Amount == 0 {
			return errors.New("amount must be non-zero for withdraws and transfers")
		}
	}

	if txn.Type == "transfer" {
		if txn.From == txn.To {
			return errors.New("from and to accounts must be different")
		}

	}

	if txn.Type == "deposit" {
		if txn.To != "" {
			return errors.New("to account must be empty for deposits")
		}
	}
	return nil
}

func (txn *Transaction) Validate(dbTx *badger.Txn) error {

	if err := txn.ValidateBasic(); err != nil {
		return err
	}

	if txn.Type == "withdraw" || txn.Type == "transfer" {
		// Check if the from account exists
		fromClient, err := FindUserByName(dbTx, txn.From)
		if err != nil {
			return err
		}

		if fromClient.Balance < txn.Amount {
			log.Printf("Insufficient funds: %d < %d", fromClient.Balance, txn.Amount)
			return errors.New("insufficient funds")
		}
	}

	if txn.Type == "transfer" {
		// Check if the to account exists
		_, err := FindUserByName(dbTx, txn.To)
		if err != nil {
			return err
		}
	}

	return nil
}

func (txn *Transaction) Apply(dbTx *badger.Txn) error {
	if err := txn.Validate(dbTx); err != nil {
		fmt.Println("Transaction validation failed: ", err)
		return err
	}

	FromClient, err := FindUserByName(dbTx, txn.From)
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) && txn.Type == "deposit" {
			fmt.Println("from account does not exist! Creating it...")
			FromClient = &Client{
				Name:    txn.From,
				Balance: 0,
			}
		} else {
			return err
		}
	}

	if txn.Type == "deposit" {
		FromClient.Balance += txn.Amount
	} else {
		FromClient.Balance -= txn.Amount
	}

	FromClientBytes, err := json.Marshal(FromClient)
	if err != nil {
		return errors.New("failed to marshal user to JSON")
	}

	err = dbTx.Set(getKeyName(FromClient.Name), FromClientBytes)
	if err != nil {
		return err
	}

	if txn.Type == "transfer" {
		ToClient, err := FindUserByName(dbTx, txn.To)
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				fmt.Println("to account does not exist! Bug??????")
				return errors.New("to account does not exist")
			} else {
				return err
			}
		}

		ToClient.Balance += txn.Amount

		ToClientBytes, err := json.Marshal(ToClient)
		if err != nil {
			return errors.New("failed to marshal user to JSON")
		}

		err = dbTx.Set(getKeyName(ToClient.Name), ToClientBytes)
		if err != nil {
			return err
		}
	}

	if txn.From != "" {
		if err := appendToHistory(dbTx, txn.From, *txn); err != nil {
			return err
		}
	}
	if txn.To != "" {
		if err := appendToHistory(dbTx, txn.To, *txn); err != nil {
			return err
		}
	}

	return nil
}

func findAllTransactions(dbTx *badger.Txn, name string) ([]Transaction, error) {
	item, err := dbTx.Get([]byte(fmt.Sprintf("transaction_%s", name)))
	if err != nil {
		return nil, err
	}

	var transactions []Transaction
	err = item.Value(func(val []byte) error {
		return json.Unmarshal(val, &transactions)
	})

	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return transactions, nil
}

func appendToHistory(dbTx *badger.Txn, name string, transaction Transaction) error {
	item, err := dbTx.Get([]byte(fmt.Sprintf("transaction_%s", name)))
	if err != nil && !errors.Is(err, badger.ErrKeyNotFound) {
		return err
	}
	var transactions []Transaction
	if item != nil {
		err = item.Value(func(val []byte) error {
			return json.Unmarshal(val, &transactions)
		})

		if err != nil {
			return err
		}

	} else {
		transactions = make([]Transaction, 0)
	}

	transactions = append(transactions, transaction)
	transactionsB, err := json.Marshal(transactions)
	if err != nil {
		return err
	}

	err = dbTx.Set([]byte(fmt.Sprintf("transaction_%s", name)), transactionsB)
	return err
}
