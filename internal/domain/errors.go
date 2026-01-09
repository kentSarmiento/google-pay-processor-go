package domain

import "errors"

var (
	// ErrInvalidToken indicates the Google Pay token is invalid or malformed
	ErrInvalidToken = errors.New("invalid google pay token")

	// ErrDecryptionFailed indicates failure to decrypt the token
	ErrDecryptionFailed = errors.New("failed to decrypt token")

	// ErrInvalidSignature indicates signature verification failed
	ErrInvalidSignature = errors.New("invalid signature")

	// ErrPaymentProcessingFailed indicates the payment processor failed
	ErrPaymentProcessingFailed = errors.New("payment processing failed")

	// ErrInvalidConfiguration indicates invalid configuration
	ErrInvalidConfiguration = errors.New("invalid configuration")

	// ErrMissingRequiredField indicates a required field is missing
	ErrMissingRequiredField = errors.New("missing required field")
)
