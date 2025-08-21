package db

import (
	"context"
	"fmt"
)

type DepositTxParams struct {
	AccountID int64 `json:"account_id"`
	Amount    int64 `json:"amount"`
}

type DepositTxResult struct {
	Account Account `json:"account"`
	Entry   Entry   `json:"entry"`
}

// DepositTx performs a money deposit to an account.
// It creates an entry and updates the account balance within a database transaction
func (store *SQLStore) DepositTx(ctx context.Context, arg DepositTxParams) (DepositTxResult, error) {
	var result DepositTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Validate amount
		if arg.Amount <= 0 {
			return fmt.Errorf("deposit amount must be positive")
		}

		// Create deposit entry
		result.Entry, err = q.CreateEntry(ctx, CreateEntryParams(arg))
		if err != nil {
			return err
		}

		// Add money to account
		result.Account, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.AccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
