package main

import (
	"bank-account/model"
	"context"
	"encoding/json"
	"fmt"
	"log"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/dgraph-io/badger/v3"
)

type BankApplication struct {
	// abcitypes.BaseApplication
	db           *badger.DB
	onGoingBlock *badger.Txn
}

var _ abcitypes.Application = (*BankApplication)(nil)

func NewBankApplication(db *badger.DB) *BankApplication {
	return &BankApplication{db: db}
}

func (app *BankApplication) Info(_ context.Context, info *abcitypes.RequestInfo) (*abcitypes.ResponseInfo, error) {
	return &abcitypes.ResponseInfo{}, nil
}

func (app *BankApplication) Query(_ context.Context, req *abcitypes.RequestQuery) (*abcitypes.ResponseQuery, error) {
	resp := abcitypes.ResponseQuery{Key: req.Data}

	log.Printf("Querying ABCI data: %v", req.Data)

	dbErr := app.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(req.Data)
		if err != nil {
			if err != badger.ErrKeyNotFound {
				return err
			}
			resp.Log = "key does not exist"
			return nil
		}

		return item.Value(func(val []byte) error {
			resp.Log = "exists"
			resp.Value = val
			return nil
		})
	})
	if dbErr != nil {
		log.Panicf("Error reading database, unable to execute query: %v", dbErr)
	}
	return &resp, nil
}

func (app *BankApplication) CheckTx(_ context.Context, check *abcitypes.RequestCheckTx) (*abcitypes.ResponseCheckTx, error) {
	var transaction model.Transaction

	if err := json.Unmarshal(check.Tx, &transaction); err != nil {
		fmt.Printf("failed to parse transaction message req: %v\n", err)

		return &abcitypes.ResponseCheckTx{Code: 1, Log: fmt.Sprintf("Failed to unmarshal: %v", err)}, nil
	}

	if err := transaction.ValidateBasic(); err != nil {
		fmt.Printf("failed to validate transaction: %v\n", err)

		return &abcitypes.ResponseCheckTx{Code: 1, Log: fmt.Sprintf("Failed to validate transaction: %v", err)}, nil
	}

	return &abcitypes.ResponseCheckTx{Code: 0}, nil
}

func (app *BankApplication) InitChain(_ context.Context, chain *abcitypes.RequestInitChain) (*abcitypes.ResponseInitChain, error) {
	return &abcitypes.ResponseInitChain{}, nil
}

func (app *BankApplication) PrepareProposal(_ context.Context, proposal *abcitypes.RequestPrepareProposal) (*abcitypes.ResponsePrepareProposal, error) {
	return &abcitypes.ResponsePrepareProposal{Txs: proposal.Txs}, nil
}

func (app *BankApplication) ProcessProposal(_ context.Context, proposal *abcitypes.RequestProcessProposal) (*abcitypes.ResponseProcessProposal, error) {
	return &abcitypes.ResponseProcessProposal{Status: abcitypes.ResponseProcessProposal_ACCEPT}, nil
}

// FinalizeBlock Deliver the decided block to the Application. The block is guaranteed to be stable and won't change anymore.
// Note: FinalizeBlock only prepares the update to be made and does not change the state of the application. The state change is actually committed in a later stage i.e. in commit phase.
func (app *BankApplication) FinalizeBlock(_ context.Context, req *abcitypes.RequestFinalizeBlock) (*abcitypes.ResponseFinalizeBlock, error) {
	fmt.Println("Executing Application FinalizeBlock")

	var txs = make([]*abcitypes.ExecTxResult, len(req.Txs))

	app.onGoingBlock = app.db.NewTransaction(true)
	for i, tx := range req.Txs {
		var transaction model.Transaction

		if err := json.Unmarshal(tx, &transaction); err != nil {
			fmt.Printf("failed to parse transaction message req: %v\n", err)

			txs[i] = &abcitypes.ExecTxResult{Code: 1, Log: fmt.Sprintf("Failed to unmarshal: %v", err)}
			continue
		}

		if err := transaction.Validate(app.onGoingBlock); err != nil {
			log.Printf("Error: invalid transaction index %v", i)
			txs[i] = &abcitypes.ExecTxResult{Code: 2, Log: fmt.Sprintf("Failed to validate transaction: %v", err)}
		} else {
			if err := transaction.Apply(app.onGoingBlock); err != nil {
				// Panic caso ocorra algum erro na escrita na aplicação. Não é possível tratar de forma determinística.
				// Quando o node se recuperar ele irá tentar novamente. Se o erro for transiente, o node irá se recuperar.
				log.Panicf("Error writing to database, unable to execute tx: %v", err)
			}

			txs[i] = &abcitypes.ExecTxResult{}
		}
	}

	return &abcitypes.ResponseFinalizeBlock{
		TxResults: txs,
	}, nil
}

func (app BankApplication) Commit(_ context.Context, commit *abcitypes.RequestCommit) (*abcitypes.ResponseCommit, error) {
	return &abcitypes.ResponseCommit{}, app.onGoingBlock.Commit()
	// return &abcitypes.ResponseCommit{}, nil
}

func (app *BankApplication) ListSnapshots(_ context.Context, snapshots *abcitypes.RequestListSnapshots) (*abcitypes.ResponseListSnapshots, error) {
	return &abcitypes.ResponseListSnapshots{}, nil
}

func (app *BankApplication) OfferSnapshot(_ context.Context, snapshot *abcitypes.RequestOfferSnapshot) (*abcitypes.ResponseOfferSnapshot, error) {
	return &abcitypes.ResponseOfferSnapshot{}, nil
}

func (app *BankApplication) LoadSnapshotChunk(_ context.Context, chunk *abcitypes.RequestLoadSnapshotChunk) (*abcitypes.ResponseLoadSnapshotChunk, error) {
	return &abcitypes.ResponseLoadSnapshotChunk{}, nil
}

func (app *BankApplication) ApplySnapshotChunk(_ context.Context, chunk *abcitypes.RequestApplySnapshotChunk) (*abcitypes.ResponseApplySnapshotChunk, error) {
	return &abcitypes.ResponseApplySnapshotChunk{Result: abcitypes.ResponseApplySnapshotChunk_ACCEPT}, nil
}

func (app BankApplication) ExtendVote(_ context.Context, extend *abcitypes.RequestExtendVote) (*abcitypes.ResponseExtendVote, error) {
	return &abcitypes.ResponseExtendVote{}, nil
}

func (app *BankApplication) VerifyVoteExtension(_ context.Context, verify *abcitypes.RequestVerifyVoteExtension) (*abcitypes.ResponseVerifyVoteExtension, error) {
	return &abcitypes.ResponseVerifyVoteExtension{}, nil
}
