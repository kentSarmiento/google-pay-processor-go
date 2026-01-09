package application

import (
	"context"
	"fmt"
	"time"

	"github.com/kentSarmiento/google-pay-processor-go/internal/domain"
)

// GooglePayService handles Google Pay payment processing
type GooglePayService struct {
	decryptor        domain.TokenDecryptor
	paymentProcessor domain.PaymentProcessor
	logger           domain.Logger
	metrics          domain.MetricsCollector
}

// NewGooglePayService creates a new Google Pay service
func NewGooglePayService(
	decryptor domain.TokenDecryptor,
	paymentProcessor domain.PaymentProcessor,
	logger domain.Logger,
	metrics domain.MetricsCollector,
) *GooglePayService {
	return &GooglePayService{
		decryptor:        decryptor,
		paymentProcessor: paymentProcessor,
		logger:           logger,
		metrics:          metrics,
	}
}

// ProcessGooglePayToken processes a Google Pay token and executes the payment
func (s *GooglePayService) ProcessGooglePayToken(
	ctx context.Context,
	token *domain.GooglePayToken,
	amount int64,
	currency string,
	merchantID string,
	transactionID string,
) (*domain.PaymentResponse, error) {
	startTime := time.Now()
	
	s.logger.Info("Processing Google Pay token", map[string]interface{}{
		"transaction_id": transactionID,
		"merchant_id":    merchantID,
		"amount":         amount,
		"currency":       currency,
	})

	// Validate input
	if err := s.validateInput(token, amount, currency, merchantID, transactionID); err != nil {
		s.metrics.IncrementCounter("googlepay_validation_failed", map[string]string{
			"merchant_id": merchantID,
		})
		s.logger.Error("Validation failed", map[string]interface{}{
			"error":          err.Error(),
			"transaction_id": transactionID,
		})
		return nil, err
	}

	// Decrypt the Google Pay token
	paymentDetails, err := s.decryptor.Decrypt(ctx, token)
	if err != nil {
		s.metrics.IncrementCounter("googlepay_decryption_failed", map[string]string{
			"merchant_id": merchantID,
		})
		s.logger.Error("Token decryption failed", map[string]interface{}{
			"error":          err.Error(),
			"transaction_id": transactionID,
		})
		return nil, fmt.Errorf("%w: %v", domain.ErrDecryptionFailed, err)
	}

	s.metrics.IncrementCounter("googlepay_decryption_success", map[string]string{
		"merchant_id": merchantID,
	})

	// Build payment request
	paymentRequest := &domain.PaymentRequest{
		CardNumber:     paymentDetails.PAN,
		ExpiryMonth:    paymentDetails.ExpirationMonth,
		ExpiryYear:     paymentDetails.ExpirationYear,
		Amount:         amount,
		Currency:       currency,
		MerchantID:     merchantID,
		TransactionID:  transactionID,
		Timestamp:      time.Now(),
	}

	// Process payment through the card payment processor
	response, err := s.paymentProcessor.Process(ctx, paymentRequest)
	if err != nil {
		s.metrics.IncrementCounter("googlepay_payment_failed", map[string]string{
			"merchant_id": merchantID,
		})
		s.logger.Error("Payment processing failed", map[string]interface{}{
			"error":          err.Error(),
			"transaction_id": transactionID,
		})
		return nil, fmt.Errorf("%w: %v", domain.ErrPaymentProcessingFailed, err)
	}

	duration := time.Since(startTime).Seconds()
	s.metrics.RecordDuration("googlepay_processing_duration", duration, map[string]string{
		"merchant_id": merchantID,
		"status":      string(response.Status),
	})

	s.metrics.IncrementCounter("googlepay_payment_processed", map[string]string{
		"merchant_id": merchantID,
		"status":      string(response.Status),
	})

	s.logger.Info("Payment processed successfully", map[string]interface{}{
		"transaction_id":   transactionID,
		"status":           response.Status,
		"authorization_id": response.AuthorizationID,
		"duration_seconds": duration,
	})

	return response, nil
}

func (s *GooglePayService) validateInput(
	token *domain.GooglePayToken,
	amount int64,
	currency string,
	merchantID string,
	transactionID string,
) error {
	if token == nil {
		return fmt.Errorf("%w: token is nil", domain.ErrInvalidToken)
	}

	if token.ProtocolVersion == "" {
		return fmt.Errorf("%w: protocol version is empty", domain.ErrInvalidToken)
	}

	if token.SignedMessage == "" {
		return fmt.Errorf("%w: signed message is empty", domain.ErrInvalidToken)
	}

	if amount <= 0 {
		return fmt.Errorf("%w: amount must be positive", domain.ErrMissingRequiredField)
	}

	if currency == "" {
		return fmt.Errorf("%w: currency is required", domain.ErrMissingRequiredField)
	}

	if merchantID == "" {
		return fmt.Errorf("%w: merchant ID is required", domain.ErrMissingRequiredField)
	}

	if transactionID == "" {
		return fmt.Errorf("%w: transaction ID is required", domain.ErrMissingRequiredField)
	}

	return nil
}
