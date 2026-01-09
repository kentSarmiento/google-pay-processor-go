package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectError bool
	}{
		{
			name: "valid configuration",
			envVars: map[string]string{
				"GOOGLEPAY_MERCHANT_ID":   "merchant-123",
				"GOOGLEPAY_ENVIRONMENT":   "TEST",
				"PAYMENT_PROCESSOR_URL":   "https://processor.example.com",
				"SERVER_PORT":             "8080",
			},
			expectError: false,
		},
		{
			name: "missing merchant ID",
			envVars: map[string]string{
				"GOOGLEPAY_ENVIRONMENT": "TEST",
				"PAYMENT_PROCESSOR_URL": "https://processor.example.com",
			},
			expectError: true,
		},
		{
			name: "invalid environment",
			envVars: map[string]string{
				"GOOGLEPAY_MERCHANT_ID": "merchant-123",
				"GOOGLEPAY_ENVIRONMENT": "INVALID",
				"PAYMENT_PROCESSOR_URL": "https://processor.example.com",
			},
			expectError: true,
		},
		{
			name: "missing processor URL",
			envVars: map[string]string{
				"GOOGLEPAY_MERCHANT_ID": "merchant-123",
				"GOOGLEPAY_ENVIRONMENT": "TEST",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			os.Clearenv()

			// Set test environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			cfg, err := LoadConfig()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if cfg == nil {
				t.Errorf("Expected config but got nil")
				return
			}

			// Verify some key values
			if merchantID, ok := tt.envVars["GOOGLEPAY_MERCHANT_ID"]; ok {
				if cfg.GooglePay.MerchantID != merchantID {
					t.Errorf("Expected merchant ID %s, got %s", merchantID, cfg.GooglePay.MerchantID)
				}
			}
		})
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
	}{
		{
			name: "valid config",
			config: &Config{
				Server: ServerConfig{Port: 8080},
				GooglePay: GooglePayConfig{
					MerchantID:  "merchant-123",
					Environment: "TEST",
				},
				Payment: PaymentConfig{
					ProcessorURL: "https://processor.example.com",
				},
			},
			expectError: false,
		},
		{
			name: "invalid port",
			config: &Config{
				Server: ServerConfig{Port: -1},
				GooglePay: GooglePayConfig{
					MerchantID:  "merchant-123",
					Environment: "TEST",
				},
				Payment: PaymentConfig{
					ProcessorURL: "https://processor.example.com",
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
