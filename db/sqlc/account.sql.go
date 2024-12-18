// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: account.sql

package db

import (
	"context"
)

const addAccountBalence = `-- name: AddAccountBalence :one
UPDATE accounts
SET balence = balence + $1
WHERE id = $2
RETURNING id, owner, balence, currency, "createdAt"
`

type AddAccountBalenceParams struct {
	Amount int64 `json:"amount"`
	ID     int64 `json:"id"`
}

func (q *Queries) AddAccountBalence(ctx context.Context, arg AddAccountBalenceParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, addAccountBalence, arg.Amount, arg.ID)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balence,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const createAccount = `-- name: CreateAccount :one
INSERT INTO accounts (
  owner,
  balence,
  currency
) VALUES (
  $1, $2, $3
) RETURNING id, owner, balence, currency, "createdAt"
`

type CreateAccountParams struct {
	Owner    string `json:"owner"`
	Balence  int64  `json:"balence"`
	Currency string `json:"currency"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, createAccount, arg.Owner, arg.Balence, arg.Currency)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balence,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const deleteAccount = `-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1
`

func (q *Queries) DeleteAccount(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteAccount, id)
	return err
}

const getAccount = `-- name: GetAccount :one
SELECT id, owner, balence, currency, "createdAt" FROM accounts
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetAccount(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccount, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balence,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const getAccountForUpdate = `-- name: GetAccountForUpdate :one
SELECT id, owner, balence, currency, "createdAt" FROM accounts
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE
`

func (q *Queries) GetAccountForUpdate(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccountForUpdate, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balence,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const listAccounts = `-- name: ListAccounts :many
SELECT id, owner, balence, currency, "createdAt" FROM accounts
WHERE owner = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListAccountsParams struct {
	Owner  string `json:"owner"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func (q *Queries) ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error) {
	rows, err := q.db.QueryContext(ctx, listAccounts, arg.Owner, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Account{}
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Balence,
			&i.Currency,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAccount = `-- name: UpdateAccount :one
UPDATE accounts
SET balence = $2
WHERE id = $1
RETURNING id, owner, balence, currency, "createdAt"
`

type UpdateAccountParams struct {
	ID      int64 `json:"id"`
	Balence int64 `json:"balence"`
}

func (q *Queries) UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, updateAccount, arg.ID, arg.Balence)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balence,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}
