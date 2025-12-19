package ports

import "context"

type PaymentGateway interface {
	// RequestPayment initiates a payment and returns the payment URL and Authority/ID
	RequestPayment(ctx context.Context, amount int64, callbackURL string, description string, email string) (paymentURL string, authority string, err error)

	// VerifyPayment verifies a payment after the user returns from the gateway
	VerifyPayment(ctx context.Context, authority string, amount int64) (refID string, err error)
}

type CardToCardGateway interface {
	// SubmitReceipt allows a user to submit a transaction receipt for manual approval
	SubmitReceipt(ctx context.Context, userID string, amount int64, receiptImageURL string, description string) (transactionID string, err error)
}
