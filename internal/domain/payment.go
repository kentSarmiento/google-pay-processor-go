package domain

import (
	"time"
)

// PaymentRequest represents a decrypted payment request from Google Pay
type PaymentRequest struct {
	CardNumber     string
	ExpiryMonth    int
	ExpiryYear     int
	CVV            string
	CardholderName string
	BillingAddress *Address
	Amount         int64
	Currency       string
	MerchantID     string
	TransactionID  string
	Timestamp      time.Time
}

// Address represents a billing or shipping address
type Address struct {
	Line1      string
	Line2      string
	City       string
	State      string
	PostalCode string
	Country    string
}

// PaymentResponse represents the result of a payment processing attempt
type PaymentResponse struct {
	TransactionID   string
	Status          PaymentStatus
	AuthorizationID string
	ProcessedAt     time.Time
	ErrorMessage    string
	ErrorCode       string
}

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusSuccess  PaymentStatus = "SUCCESS"
	PaymentStatusFailed   PaymentStatus = "FAILED"
	PaymentStatusPending  PaymentStatus = "PENDING"
	PaymentStatusDeclined PaymentStatus = "DECLINED"
	PaymentStatusError    PaymentStatus = "ERROR"
)

// GooglePayToken represents the encrypted token from Google Pay
type GooglePayToken struct {
	Signature              string                  `json:"signature"`
	ProtocolVersion        string                  `json:"protocolVersion"`
	SignedMessage          string                  `json:"signedMessage"`
	IntermediateSigningKey *IntermediateSigningKey `json:"intermediateSigningKey,omitempty"`
}

// IntermediateSigningKey represents the intermediate signing key in the token
type IntermediateSigningKey struct {
	SignedKey  string   `json:"signedKey"`
	Signatures []string `json:"signatures"`
}

// SignedMessage represents the signed message within the Google Pay token
type SignedMessage struct {
	EncryptedMessage   string `json:"encryptedMessage"`
	EphemeralPublicKey string `json:"ephemeralPublicKey"`
	Tag                string `json:"tag"`
}

// PaymentMethodDetails represents the decrypted payment details
type PaymentMethodDetails struct {
	PAN             string `json:"pan"`
	ExpirationMonth int    `json:"expirationMonth"`
	ExpirationYear  int    `json:"expirationYear"`
	AuthMethod      string `json:"authMethod"`
}
