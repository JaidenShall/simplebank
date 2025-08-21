package db

import (
	"context"
	"fmt"
)

type WithdrawTxParams struct {
	AccountID int64 `json:"account_id"`
	Amount    int64 `json:"amount"`
}

type WithdrawTxResult struct {
	Account Account `json:"account"`
	Entry   Entry   `json:"entry"`
}

// WithdrawTx performs a money withdrawal from an account.
// It creates an entry and updates the account balance within a database transaction
func (store *SQLStore) WithdrawTx(ctx context.Context, arg WithdrawTxParams) (WithdrawTxResult, error) {
	var result WithdrawTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Validate amount
		if arg.Amount <= 0 {
			return fmt.Errorf("withdrawal amount must be positive")
		}

		// Get current account to check balance
		account, err := q.GetAccountForUpdate(ctx, arg.AccountID)
		if err != nil {
			return err
		}

		// Check sufficient balance
		if account.Balance < arg.Amount {
			return fmt.Errorf("insufficient balance: current=%d, requested=%d", account.Balance, arg.Amount)
		}

		// Create withdrawal entry
		result.Entry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.AccountID,
			Amount:    -arg.Amount, // Negative amount for withdrawal
		})
		if err != nil {
			return err
		}

		// Subtract money from account
		result.Account, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.AccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
