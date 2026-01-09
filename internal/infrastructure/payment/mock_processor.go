package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kentSarmiento/google-pay-processor-go/internal/domain"
)

// MockCardProcessor is a mock implementation of the card payment processor
// In production, this would integrate with an actual payment processor
type MockCardProcessor struct {
	logger domain.Logger
}

// NewMockCardProcessor creates a new mock card payment processor
func NewMockCardProcessor(logger domain.Logger) domain.PaymentProcessor {
	return &MockCardProcessor{
		logger: logger,
	}
}

// Process processes a payment request (mock implementation)
func (m *MockCardProcessor) Process(ctx context.Context, request *domain.PaymentRequest) (*domain.PaymentResponse, error) {
	m.logger.Info("Processing card payment", map[string]interface{}{
		"transaction_id": request.TransactionID,
		"amount":         request.Amount,
		"currency":       request.Currency,
		"masked_card":    maskCardNumber(request.CardNumber),
	})

	// Validate request
	if err := m.validateRequest(request); err != nil {
		return &domain.PaymentResponse{
			TransactionID: request.TransactionID,
			Status:        domain.PaymentStatusError,
			ErrorMessage:  err.Error(),
			ProcessedAt:   time.Now(),
		}, err
	}

	// Simulate payment processing
	// In production, this would make API calls to the actual payment processor
	
	// Generate authorization ID
	authID := uuid.New().String()

	// Simulate some processing time
	time.Sleep(50 * time.Millisecond)

	// Mock response based on card number (for testing purposes)
	status := m.determineStatus(request.CardNumber)

	response := &domain.PaymentResponse{
		TransactionID:   request.TransactionID,
		Status:          status,
		AuthorizationID: authID,
		ProcessedAt:     time.Now(),
	}

	if status == domain.PaymentStatusFailed || status == domain.PaymentStatusDeclined {
		response.ErrorMessage = "Payment declined by issuer"
		response.ErrorCode = "DECLINED"
	}

	m.logger.Info("Payment processed", map[string]interface{}{
		"transaction_id":   request.TransactionID,
		"authorization_id": authID,
		"status":           status,
	})

	return response, nil
}

func (m *MockCardProcessor) validateRequest(request *domain.PaymentRequest) error {
	if request.CardNumber == "" {
		return fmt.Errorf("%w: card number is required", domain.ErrMissingRequiredField)
	}

	if len(request.CardNumber) < 13 || len(request.CardNumber) > 19 {
		return fmt.Errorf("%w: invalid card number length", domain.ErrPaymentProcessingFailed)
	}

	if request.ExpiryMonth < 1 || request.ExpiryMonth > 12 {
		return fmt.Errorf("%w: invalid expiry month", domain.ErrPaymentProcessingFailed)
	}

	if request.ExpiryYear < time.Now().Year() {
		return fmt.Errorf("%w: card has expired", domain.ErrPaymentProcessingFailed)
	}

	if request.Amount <= 0 {
		return fmt.Errorf("%w: amount must be positive", domain.ErrPaymentProcessingFailed)
	}

	return nil
}

func (m *MockCardProcessor) determineStatus(cardNumber string) domain.PaymentStatus {
	// For testing: use last digit to determine status
	// In production, this would be determined by the actual processor response
	lastDigit := cardNumber[len(cardNumber)-1]

	switch lastDigit {
	case '0':
		return domain.PaymentStatusDeclined
	case '1', '2', '3', '4', '5', '6', '7', '8':
		return domain.PaymentStatusSuccess
	case '9':
		return domain.PaymentStatusPending
	default:
		return domain.PaymentStatusSuccess
	}
}

func maskCardNumber(cardNumber string) string {
	if len(cardNumber) <= 4 {
		return "****"
	}
	lastFour := cardNumber[len(cardNumber)-4:]
	return "************" + lastFour
}
