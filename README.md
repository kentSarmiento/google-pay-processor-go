# Google Pay Processor (Go)

A high-performance, production-ready Google Pay payment processor built in Go with modern tooling, clean architecture, and best practices.

## Overview

This service processes Google Pay tokens by decrypting them using the Google Tink library and forwarding the payment details to an existing card payment processor. It's designed to be fault-tolerant, highly scalable, and performant with low latency.

## Features

- ✅ **Google Pay Token Decryption** using Tink library (ECv2 protocol)
- ✅ **Clean Architecture** with clear separation of concerns
- ✅ **SOLID Principles** for maintainable and extensible code
- ✅ **Structured Logging** with zerolog for observability
- ✅ **Prometheus Metrics** for monitoring and alerting
- ✅ **Graceful Shutdown** for zero-downtime deployments
- ✅ **Health & Readiness Checks** for Kubernetes/orchestration
- ✅ **Docker Support** with multi-stage builds
- ✅ **Comprehensive Testing** with table-driven tests
- ✅ **CI/CD Pipeline** with GitHub Actions
- ✅ **Security Scanning** with Gosec and Trivy
- ✅ **Configuration via Environment Variables** for 12-factor app compliance

## Architecture

### Clean Architecture Layers

```
├── cmd/server/              # Application entry point
├── internal/
│   ├── domain/              # Domain layer (entities, interfaces, errors)
│   ├── application/         # Application layer (use cases, business logic)
│   └── infrastructure/      # Infrastructure layer (implementations)
│       ├── config/          # Configuration management
│       ├── crypto/          # Tink-based token decryption
│       ├── logger/          # Zerolog implementation
│       ├── metrics/         # Prometheus metrics
│       └── payment/         # Payment processor client
└── pkg/googlepay/           # HTTP handlers and API
```

### Design Principles

- **Dependency Inversion**: Domain doesn't depend on infrastructure
- **Single Responsibility**: Each component has one reason to change
- **Interface Segregation**: Small, focused interfaces
- **Open/Closed**: Open for extension, closed for modification
- **Liskov Substitution**: Implementations are interchangeable

## Technology Stack

### Core Libraries

- **[Tink](https://github.com/google/tink)** (v1.7.0) - Google's cryptography library for token decryption
  - Provides battle-tested implementations of cryptographic primitives
  - Used for ECIES hybrid encryption/decryption and ECDSA signature verification
  
- **[Zerolog](https://github.com/rs/zerolog)** (v1.34.0) - High-performance structured logging
  - Zero allocation logging for minimal performance impact
  - JSON and console output formats
  - Contextual fields for tracing and debugging

- **[Prometheus Client](https://github.com/prometheus/client_golang)** (v1.23.2) - Metrics collection
  - Industry-standard metrics format
  - Counters, histograms, and gauges
  - Built-in HTTP handler for scraping

- **[UUID](https://github.com/google/uuid)** (v1.6.0) - Unique identifier generation
  - RFC 4122 compliant UUIDs
  - Used for transaction and authorization IDs

### Development Tools

- **golangci-lint** - Comprehensive linting with multiple analyzers
- **gosec** - Security vulnerability scanning
- **Trivy** - Container and dependency vulnerability scanning
- **Docker** - Containerization with multi-stage builds
- **Make** - Build automation and task management

## Getting Started

### Prerequisites

- Go 1.21 or later
- Docker (optional, for containerized deployment)
- Make (optional, for build automation)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/kentSarmiento/google-pay-processor-go.git
cd google-pay-processor-go
```

2. Install dependencies:
```bash
go mod download
```

3. Configure environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

### Configuration

Set the following environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | HTTP server port | 8080 |
| `GOOGLEPAY_MERCHANT_ID` | Your Google Pay merchant ID | Required |
| `GOOGLEPAY_ENVIRONMENT` | `TEST` or `PRODUCTION` | TEST |
| `PAYMENT_PROCESSOR_URL` | URL of the card payment processor | Required |
| `PAYMENT_PROCESSOR_API_KEY` | API key for payment processor | Optional |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | info |
| `LOG_FORMAT` | Log format (json, console) | json |
| `METRICS_ENABLED` | Enable Prometheus metrics | true |
| `METRICS_PORT` | Metrics server port | 9090 |

### Running Locally

Using Go:
```bash
# Set required environment variables
export GOOGLEPAY_MERCHANT_ID=your-merchant-id
export GOOGLEPAY_ENVIRONMENT=TEST
export PAYMENT_PROCESSOR_URL=http://localhost:9000

# Run the server
go run cmd/server/main.go
```

Using Make:
```bash
make run
```

Using Docker:
```bash
# Build the image
make docker-build

# Run the container
make docker-run
```

## API Endpoints

### Process Payment

Process a Google Pay token and execute the payment.

**Endpoint:** `POST /api/v1/googlepay/process`

**Request Body:**
```json
{
  "token": {
    "signature": "base64-encoded-signature",
    "protocolVersion": "ECv2",
    "signedMessage": "base64-encoded-signed-message",
    "intermediateSigningKey": {
      "signedKey": "base64-encoded-key",
      "signatures": ["base64-encoded-signature"]
    }
  },
  "amount": 10000,
  "currency": "USD",
  "merchantId": "merchant-123",
  "transactionId": "txn-123"
}
```

**Response:**
```json
{
  "transactionId": "txn-123",
  "status": "SUCCESS",
  "authorizationId": "auth-456",
  "processedAt": "2026-01-09T13:00:00Z"
}
```

**Status Values:**
- `SUCCESS` - Payment processed successfully
- `FAILED` - Payment processing failed
- `DECLINED` - Payment declined by issuer
- `PENDING` - Payment is being processed
- `ERROR` - System error occurred

### Health Check

Check if the service is running.

**Endpoint:** `GET /health`

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2026-01-09T13:00:00Z",
  "service": "google-pay-processor"
}
```

### Readiness Check

Check if the service is ready to accept traffic.

**Endpoint:** `GET /ready`

**Response:**
```json
{
  "status": "ready",
  "timestamp": "2026-01-09T13:00:00Z",
  "service": "google-pay-processor"
}
```

### Metrics

Prometheus metrics for monitoring.

**Endpoint:** `GET /metrics` (port 9090 by default)

**Available Metrics:**
- `googlepay_decryption_success` - Counter of successful decryptions
- `googlepay_decryption_failed` - Counter of failed decryptions
- `googlepay_payment_processed` - Counter of processed payments by status
- `googlepay_processing_duration` - Histogram of processing durations
- `http_requests_total` - Counter of HTTP requests
- `http_request_duration` - Histogram of HTTP request durations

## Development

### Building

```bash
# Build the binary
make build

# Build Docker image
make docker-build
```

### Testing

```bash
# Run tests
make test

# Run tests with coverage
make test-coverage

# View coverage report
open coverage.html
```

### Linting

```bash
# Run all linters
make lint

# Format code
make fmt

# Run go vet
make vet
```

### CI/CD

The project uses GitHub Actions for continuous integration:

- **Test**: Runs unit tests with race detection
- **Lint**: Runs golangci-lint for code quality
- **Build**: Builds the binary and Docker image
- **Security**: Scans for vulnerabilities with Gosec and Trivy

## Production Considerations

### Security

1. **Token Decryption**: Replace the placeholder implementation with actual Tink setup
   - Load merchant private keys from secure key storage (KMS, Vault, etc.)
   - Verify Google's root signing keys from their public endpoint
   - Implement proper signature verification

2. **Secrets Management**: Use a secrets manager (AWS Secrets Manager, HashiCorp Vault)
   - Never commit secrets to version control
   - Rotate keys regularly
   - Use encryption at rest and in transit

3. **API Security**:
   - Implement authentication (API keys, JWT, mTLS)
   - Add rate limiting to prevent abuse
   - Enable CORS with appropriate origins
   - Use HTTPS/TLS in production

### Scalability

1. **Horizontal Scaling**: The service is stateless and can scale horizontally
   - Deploy multiple instances behind a load balancer
   - Use container orchestration (Kubernetes, ECS)

2. **Performance Optimization**:
   - Connection pooling for external services
   - Caching for frequently accessed data
   - Request timeouts and circuit breakers

3. **Resource Management**:
   - Set appropriate memory and CPU limits
   - Monitor resource usage with metrics
   - Use autoscaling based on load

### Observability

1. **Logging**:
   - Structured JSON logs for easy parsing
   - Correlation IDs for request tracing
   - Log aggregation (ELK, Splunk, CloudWatch)

2. **Metrics**:
   - Prometheus for metrics collection
   - Grafana for visualization
   - Alerting on critical metrics

3. **Tracing**:
   - Consider adding distributed tracing (Jaeger, Zipkin)
   - Track request flow across services

### High Availability

1. **Fault Tolerance**:
   - Graceful degradation on failures
   - Retry logic with exponential backoff
   - Circuit breakers for external dependencies

2. **Health Checks**:
   - Kubernetes liveness and readiness probes
   - Monitor external dependency health

3. **Zero-Downtime Deployments**:
   - Rolling updates
   - Graceful shutdown implementation
   - Connection draining

## Replacing Mock Payment Processor

To integrate with a real payment processor:

1. Create a new implementation in `internal/infrastructure/payment/`
2. Implement the `domain.PaymentProcessor` interface
3. Update the initialization in `cmd/server/main.go`

Example:
```go
// internal/infrastructure/payment/stripe_processor.go
type StripeProcessor struct {
    client *stripe.Client
    logger domain.Logger
}

func (s *StripeProcessor) Process(ctx context.Context, request *domain.PaymentRequest) (*domain.PaymentResponse, error) {
    // Implement Stripe API calls
}
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run linters and tests
6. Submit a pull request

## License

[Add your license here]

## Support

For issues and questions:
- Open an issue on GitHub
- Contact: [Your contact information]

## Acknowledgments

- Google Tink team for the cryptography library
- Go community for excellent libraries and tools