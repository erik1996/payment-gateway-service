package utils

import "errors"

type PaymentType string

const (
	PaymentTypeDeposit    PaymentType = "DEPOSIT"
	PaymentTypeWithdrawal PaymentType = "WITHDRAWAL"
)

type PaymentStatus string

const (
	PaymentStatusInitialized PaymentStatus = "INITIALIZED"
	PaymentStatusPending     PaymentStatus = "PENDING"
	PaymentStatusSuccess     PaymentStatus = "SUCCESS"
	PaymentStatusFailed      PaymentStatus = "FAILED"
)

// Define the error for invalid transaction type
var ErrInvalidTransactionType = errors.New("invalid transaction type")
