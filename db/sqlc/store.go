package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Queries only support one operation on one table so it dont support transactions
// so we do *Queries : called composition in golang fi placet el wirata

// store provide all funcs to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// create a new store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx: execute a function within a db transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// tx, db := store.db.BeginTx(ctx, &sql.TxOptions{})
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams: input params
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult : result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx : money transfer from account to other
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		// add account entries
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// update balence
		// get account => update it balence
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err =
				addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err =
				addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalence(ctx, AddAccountBalenceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalence(ctx, AddAccountBalenceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return
}

// create entrey & update balance
type EntryBalenceTxParams struct {
	AccountId int64 `json:"account_id"`
	Amount    int64 `json:"amount"`
}

type EntryBalenceTxResult struct {
	Entry   Entry   `json:"entry"`
	Account Account `json:"account"`
}

func (store *Store) EntryBalenceTx(ctx context.Context, arg EntryBalenceTxParams) (EntryBalenceTxResult, error) {
	var result EntryBalenceTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Entry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.AccountId,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		result.Account, err = q.AddAccountBalence(ctx, AddAccountBalenceParams{
			Amount: arg.Amount,
			ID:     arg.AccountId,
		})
		if err != nil {
			return err
		}

		return err
	})

	return result, err
}
