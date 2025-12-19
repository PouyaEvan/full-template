package cardtocard

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CardToCardAdapter handles manual receipt submissions
type CardToCardAdapter struct {
	// In a real app, this would likely interact with a repository to store the pending transaction
}

func NewCardToCardAdapter() *CardToCardAdapter {
	return &CardToCardAdapter{}
}

func (c *CardToCardAdapter) SubmitReceipt(ctx context.Context, userID string, amount int64, receiptImageURL string, description string) (string, error) {
	// 1. Validate inputs
	if amount <= 0 {
		return "", fmt.Errorf("invalid amount")
	}
	if receiptImageURL == "" {
		return "", fmt.Errorf("receipt image is required")
	}

	// 2. Generate a Transaction ID
	txID := uuid.New().String()

	// 3. Log or Store the transaction request (Mocking storage here)
	// In production: repo.SaveTransaction(Transaction{Status: "PENDING_APPROVAL", ...})
	fmt.Printf("[CardToCard] New Receipt Submitted: User=%s Amount=%d Image=%s TxID=%s\n", userID, amount, receiptImageURL, txID)

	// Simulate processing time
	time.Sleep(100 * time.Millisecond)

	return txID, nil
}
