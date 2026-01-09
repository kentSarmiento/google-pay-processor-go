package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/kentSarmiento/google-pay-processor-go/internal/domain"
)

// Config holds application configuration
type Config struct {
	Server    ServerConfig
	GooglePay GooglePayConfig
	Payment   PaymentConfig
	Logging   LoggingConfig
	Metrics   MetricsConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
}

// GooglePayConfig holds Google Pay specific configuration
type GooglePayConfig struct {
	MerchantID      string
	MerchantName    string
	ProtocolVersion string
	SigningKeys     []string
	Environment     string // PRODUCTION or TEST
}

// PaymentConfig holds payment processor configuration
type PaymentConfig struct {
	ProcessorURL string
	APIKey       string
	Timeout      int
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string
	Format string // json or console
}

// MetricsConfig holds metrics configuration
type MetricsConfig struct {
	Enabled bool
	Port    int
	Path    string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port:         getEnvAsInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvAsInt("SERVER_READ_TIMEOUT", 15),
			WriteTimeout: getEnvAsInt("SERVER_WRITE_TIMEOUT", 15),
			IdleTimeout:  getEnvAsInt("SERVER_IDLE_TIMEOUT", 60),
		},
		GooglePay: GooglePayConfig{
			MerchantID:      getEnv("GOOGLEPAY_MERCHANT_ID", ""),
			MerchantName:    getEnv("GOOGLEPAY_MERCHANT_NAME", ""),
			ProtocolVersion: getEnv("GOOGLEPAY_PROTOCOL_VERSION", "ECv2"),
			Environment:     getEnv("GOOGLEPAY_ENVIRONMENT", "TEST"),
		},
		Payment: PaymentConfig{
			ProcessorURL: getEnv("PAYMENT_PROCESSOR_URL", ""),
			APIKey:       getEnv("PAYMENT_PROCESSOR_API_KEY", ""),
			Timeout:      getEnvAsInt("PAYMENT_PROCESSOR_TIMEOUT", 30),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		Metrics: MetricsConfig{
			Enabled: getEnvAsBool("METRICS_ENABLED", true),
			Port:    getEnvAsInt("METRICS_PORT", 9090),
			Path:    getEnv("METRICS_PATH", "/metrics"),
		},
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.GooglePay.MerchantID == "" {
		return fmt.Errorf("%w: GOOGLEPAY_MERCHANT_ID is required", domain.ErrInvalidConfiguration)
	}

	if c.GooglePay.Environment != "PRODUCTION" && c.GooglePay.Environment != "TEST" {
		return fmt.Errorf("%w: GOOGLEPAY_ENVIRONMENT must be PRODUCTION or TEST", domain.ErrInvalidConfiguration)
	}

	if c.Payment.ProcessorURL == "" {
		return fmt.Errorf("%w: PAYMENT_PROCESSOR_URL is required", domain.ErrInvalidConfiguration)
	}

	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("%w: invalid server port", domain.ErrInvalidConfiguration)
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
