package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDepositTx(t *testing.T) {
	store := NewStore(testDB)

	account := createRandomAccount(t)
	fmt.Println(">> before deposit:", account.Balance)

	amount := int64(500)

	result, err := store.DepositTx(context.Background(), DepositTxParams{
		AccountID: account.ID,
		Amount:    amount,
	})

	require.NoError(t, err)
	require.NotEmpty(t, result)

	// Check account balance increased
	require.Equal(t, account.Balance+amount, result.Account.Balance)

	// Check entry was created
	require.Equal(t, account.ID, result.Entry.AccountID)
	require.Equal(t, amount, result.Entry.Amount)

	fmt.Println(">> after deposit:", result.Account.Balance)
}

func TestWithdrawTx(t *testing.T) {
	store := NewStore(testDB)

	account := createRandomAccount(t)
	// Ensure account has sufficient balance for test
	initialBalance := account.Balance
	if initialBalance < 1000 {
		// Add some money first
		depositResult, err := store.DepositTx(context.Background(), DepositTxParams{
			AccountID: account.ID,
			Amount:    1000,
		})
		require.NoError(t, err)
		account = depositResult.Account
	}

	fmt.Println(">> before withdrawal:", account.Balance)

	amount := int64(300)

	result, err := store.WithdrawTx(context.Background(), WithdrawTxParams{
		AccountID: account.ID,
		Amount:    amount,
	})

	require.NoError(t, err)
	require.NotEmpty(t, result)

	// Check account balance decreased
	require.Equal(t, account.Balance-amount, result.Account.Balance)

	// Check entry was created with negative amount
	require.Equal(t, account.ID, result.Entry.AccountID)
	require.Equal(t, -amount, result.Entry.Amount)

	fmt.Println(">> after withdrawal:", result.Account.Balance)
}

func TestWithdrawTxInsufficientBalance(t *testing.T) {
	store := NewStore(testDB)

	account := createRandomAccount(t)
	amount := account.Balance + 1000 // More than available

	_, err := store.WithdrawTx(context.Background(), WithdrawTxParams{
		AccountID: account.ID,
		Amount:    amount,
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "insufficient balance")
}
