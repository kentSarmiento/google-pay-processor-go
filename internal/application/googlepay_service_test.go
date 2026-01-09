package application

import (
	"context"
	"errors"
	"testing"

	"github.com/kentSarmiento/google-pay-processor-go/internal/domain"
)

// Mock implementations for testing
type mockDecryptor struct {
	decryptFunc func(ctx context.Context, token *domain.GooglePayToken) (*domain.PaymentMethodDetails, error)
}

func (m *mockDecryptor) Decrypt(ctx context.Context, token *domain.GooglePayToken) (*domain.PaymentMethodDetails, error) {
	if m.decryptFunc != nil {
		return m.decryptFunc(ctx, token)
	}
	return &domain.PaymentMethodDetails{
		PAN:             "4111111111111111",
		ExpirationMonth: 12,
		ExpirationYear:  2025,
		AuthMethod:      "PAN_ONLY",
	}, nil
}

type mockPaymentProcessor struct {
	processFunc func(ctx context.Context, request *domain.PaymentRequest) (*domain.PaymentResponse, error)
}

func (m *mockPaymentProcessor) Process(ctx context.Context, request *domain.PaymentRequest) (*domain.PaymentResponse, error) {
	if m.processFunc != nil {
		return m.processFunc(ctx, request)
	}
	return &domain.PaymentResponse{
		TransactionID:   request.TransactionID,
		Status:          domain.PaymentStatusSuccess,
		AuthorizationID: "auth-123",
	}, nil
}

type mockLogger struct{}

func (m *mockLogger) Debug(msg string, fields map[string]interface{}) {}
func (m *mockLogger) Info(msg string, fields map[string]interface{})  {}
func (m *mockLogger) Warn(msg string, fields map[string]interface{})  {}
func (m *mockLogger) Error(msg string, fields map[string]interface{}) {}
func (m *mockLogger) Fatal(msg string, fields map[string]interface{}) {}

type mockMetrics struct{}

func (m *mockMetrics) IncrementCounter(name string, labels map[string]string)                 {}
func (m *mockMetrics) RecordDuration(name string, duration float64, labels map[string]string) {}
func (m *mockMetrics) SetGauge(name string, value float64, labels map[string]string)          {}

func TestGooglePayService_ProcessGooglePayToken(t *testing.T) {
	tests := []struct {
		name              string
		token             *domain.GooglePayToken
		amount            int64
		currency          string
		merchantID        string
		transactionID     string
		decryptFunc       func(ctx context.Context, token *domain.GooglePayToken) (*domain.PaymentMethodDetails, error)
		processFunc       func(ctx context.Context, request *domain.PaymentRequest) (*domain.PaymentResponse, error)
		expectedStatus    domain.PaymentStatus
		expectError       bool
		expectedErrorType error
	}{
		{
			name: "successful payment processing",
			token: &domain.GooglePayToken{
				ProtocolVersion: "ECv2",
				SignedMessage:   "signed-message",
				Signature:       "signature",
			},
			amount:         10000,
			currency:       "USD",
			merchantID:     "merchant-123",
			transactionID:  "txn-123",
			expectedStatus: domain.PaymentStatusSuccess,
			expectError:    false,
		},
		{
			name:              "nil token should fail validation",
			token:             nil,
			amount:            10000,
			currency:          "USD",
			merchantID:        "merchant-123",
			transactionID:     "txn-123",
			expectError:       true,
			expectedErrorType: domain.ErrInvalidToken,
		},
		{
			name: "empty protocol version should fail validation",
			token: &domain.GooglePayToken{
				ProtocolVersion: "",
				SignedMessage:   "signed-message",
			},
			amount:            10000,
			currency:          "USD",
			merchantID:        "merchant-123",
			transactionID:     "txn-123",
			expectError:       true,
			expectedErrorType: domain.ErrInvalidToken,
		},
		{
			name: "zero amount should fail validation",
			token: &domain.GooglePayToken{
				ProtocolVersion: "ECv2",
				SignedMessage:   "signed-message",
			},
			amount:            0,
			currency:          "USD",
			merchantID:        "merchant-123",
			transactionID:     "txn-123",
			expectError:       true,
			expectedErrorType: domain.ErrMissingRequiredField,
		},
		{
			name: "decryption failure",
			token: &domain.GooglePayToken{
				ProtocolVersion: "ECv2",
				SignedMessage:   "signed-message",
				Signature:       "signature",
			},
			amount:        10000,
			currency:      "USD",
			merchantID:    "merchant-123",
			transactionID: "txn-123",
			decryptFunc: func(ctx context.Context, token *domain.GooglePayToken) (*domain.PaymentMethodDetails, error) {
				return nil, errors.New("decryption failed")
			},
			expectError:       true,
			expectedErrorType: domain.ErrDecryptionFailed,
		},
		{
			name: "payment processing failure",
			token: &domain.GooglePayToken{
				ProtocolVersion: "ECv2",
				SignedMessage:   "signed-message",
				Signature:       "signature",
			},
			amount:        10000,
			currency:      "USD",
			merchantID:    "merchant-123",
			transactionID: "txn-123",
			processFunc: func(ctx context.Context, request *domain.PaymentRequest) (*domain.PaymentResponse, error) {
				return nil, errors.New("payment failed")
			},
			expectError:       true,
			expectedErrorType: domain.ErrPaymentProcessingFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decryptor := &mockDecryptor{decryptFunc: tt.decryptFunc}
			processor := &mockPaymentProcessor{processFunc: tt.processFunc}
			logger := &mockLogger{}
			metrics := &mockMetrics{}

			service := NewGooglePayService(decryptor, processor, logger, metrics)

			ctx := context.Background()
			response, err := service.ProcessGooglePayToken(
				ctx,
				tt.token,
				tt.amount,
				tt.currency,
				tt.merchantID,
				tt.transactionID,
			)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if tt.expectedErrorType != nil && !errors.Is(err, tt.expectedErrorType) {
					t.Errorf("Expected error type %v, got %v", tt.expectedErrorType, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if response == nil {
				t.Errorf("Expected response but got nil")
				return
			}

			if response.Status != tt.expectedStatus {
				t.Errorf("Expected status %v, got %v", tt.expectedStatus, response.Status)
			}
		})
	}
}
