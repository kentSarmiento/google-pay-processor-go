package domain

import "context"

// TokenDecryptor defines the interface for decrypting Google Pay tokens
type TokenDecryptor interface {
	// Decrypt decrypts a Google Pay token and returns payment method details
	Decrypt(ctx context.Context, token *GooglePayToken) (*PaymentMethodDetails, error)
}

// PaymentProcessor defines the interface for processing card payments
type PaymentProcessor interface {
	// Process processes a payment request and returns the result
	Process(ctx context.Context, request *PaymentRequest) (*PaymentResponse, error)
}

// Logger defines the interface for structured logging
type Logger interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
	Fatal(msg string, fields map[string]interface{})
}

// MetricsCollector defines the interface for collecting metrics
type MetricsCollector interface {
	// IncrementCounter increments a counter metric
	IncrementCounter(name string, labels map[string]string)
	
	// RecordDuration records a duration metric
	RecordDuration(name string, duration float64, labels map[string]string)
	
	// SetGauge sets a gauge metric
	SetGauge(name string, value float64, labels map[string]string)
}
