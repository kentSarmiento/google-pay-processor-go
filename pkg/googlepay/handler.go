package googlepay

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/kentSarmiento/google-pay-processor-go/internal/application"
	"github.com/kentSarmiento/google-pay-processor-go/internal/domain"
)

// Handler handles HTTP requests for Google Pay processing
type Handler struct {
	service *application.GooglePayService
	logger  domain.Logger
	metrics domain.MetricsCollector
}

// NewHandler creates a new Google Pay handler
func NewHandler(service *application.GooglePayService, logger domain.Logger, metrics domain.MetricsCollector) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
		metrics: metrics,
	}
}

// ProcessPaymentRequest represents the incoming payment request
type ProcessPaymentRequest struct {
	Token         domain.GooglePayToken `json:"token"`
	Amount        int64                 `json:"amount"`
	Currency      string                `json:"currency"`
	MerchantID    string                `json:"merchantId"`
	TransactionID string                `json:"transactionId"`
}

// ProcessPaymentResponse represents the payment response
type ProcessPaymentResponse struct {
	TransactionID   string `json:"transactionId"`
	Status          string `json:"status"`
	AuthorizationID string `json:"authorizationId,omitempty"`
	ErrorMessage    string `json:"errorMessage,omitempty"`
	ErrorCode       string `json:"errorCode,omitempty"`
	ProcessedAt     string `json:"processedAt"`
}

// ProcessPayment handles Google Pay payment processing requests
func (h *Handler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	h.metrics.IncrementCounter("http_requests_total", map[string]string{
		"endpoint": "/api/v1/googlepay/process",
		"method":   r.Method,
	})

	// Only accept POST requests
	if r.Method != http.MethodPost {
		h.respondWithError(w, http.StatusMethodNotAllowed, "method not allowed", "METHOD_NOT_ALLOWED")
		return
	}

	// Parse request body
	var req ProcessPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Failed to parse request body", map[string]interface{}{
			"error": err.Error(),
		})
		h.respondWithError(w, http.StatusBadRequest, "invalid request body", "INVALID_REQUEST")
		return
	}

	// Process the payment
	ctx := r.Context()
	response, err := h.service.ProcessGooglePayToken(
		ctx,
		&req.Token,
		req.Amount,
		req.Currency,
		req.MerchantID,
		req.TransactionID,
	)

	duration := time.Since(startTime).Seconds()
	h.metrics.RecordDuration("http_request_duration", duration, map[string]string{
		"endpoint": "/api/v1/googlepay/process",
	})

	if err != nil {
		h.logger.Error("Payment processing failed", map[string]interface{}{
			"error":          err.Error(),
			"transaction_id": req.TransactionID,
		})
		h.respondWithError(w, http.StatusInternalServerError, err.Error(), "PROCESSING_FAILED")
		return
	}

	// Build response
	respData := ProcessPaymentResponse{
		TransactionID:   response.TransactionID,
		Status:          string(response.Status),
		AuthorizationID: response.AuthorizationID,
		ErrorMessage:    response.ErrorMessage,
		ErrorCode:       response.ErrorCode,
		ProcessedAt:     response.ProcessedAt.Format(time.RFC3339),
	}

	h.respondWithJSON(w, http.StatusOK, respData)
}

// HealthCheck handles health check requests
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "google-pay-processor",
	}
	h.respondWithJSON(w, http.StatusOK, health)
}

// ReadinessCheck handles readiness check requests
func (h *Handler) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	readiness := map[string]interface{}{
		"status":    "ready",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "google-pay-processor",
	}
	h.respondWithJSON(w, http.StatusOK, readiness)
}

func (h *Handler) respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode response", map[string]interface{}{
			"error": err.Error(),
		})
	}
}

func (h *Handler) respondWithError(w http.ResponseWriter, statusCode int, message, code string) {
	errorResp := map[string]string{
		"error":     message,
		"errorCode": code,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	h.respondWithJSON(w, statusCode, errorResp)
}
