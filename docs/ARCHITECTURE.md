# Architecture Documentation

## Overview

The Google Pay Processor follows Clean Architecture principles with clear separation of concerns and dependency inversion. This document explains the architectural decisions and patterns used.

## Architecture Layers

### 1. Domain Layer (`internal/domain/`)

The innermost layer containing enterprise business logic and rules.

**Components:**
- **Entities**: Core business objects (`PaymentRequest`, `PaymentResponse`, `GooglePayToken`)
- **Interfaces (Ports)**: Abstract interfaces for external dependencies
- **Errors**: Domain-specific error types
- **Business Rules**: Validation and invariants

**Characteristics:**
- No external dependencies
- Pure Go code
- Framework agnostic
- Most stable layer (changes least frequently)

**Key Files:**
- `payment.go` - Payment entities and value objects
- `ports.go` - Interface definitions for external systems
- `errors.go` - Domain error types

### 2. Application Layer (`internal/application/`)

Contains application-specific business logic and use cases.

**Components:**
- **Services**: Orchestrate domain objects and coordinate workflows
- **Use Cases**: Application-specific business flows

**Characteristics:**
- Depends only on domain layer
- Orchestrates domain objects
- Coordinates infrastructure components through interfaces
- Contains transaction boundaries

**Key Files:**
- `googlepay_service.go` - Main service orchestrating Google Pay processing

**Flow:**
```
ProcessGooglePayToken
  → Validate input
  → Decrypt token (via TokenDecryptor port)
  → Build payment request
  → Process payment (via PaymentProcessor port)
  → Record metrics
  → Return response
```

### 3. Infrastructure Layer (`internal/infrastructure/`)

Contains implementations of domain interfaces and external system integrations.

**Components:**

#### Config (`infrastructure/config/`)
- Environment variable loading
- Configuration validation
- 12-factor app compliance

#### Crypto (`infrastructure/crypto/`)
- Tink-based token decryption
- Signature verification
- Key management integration

#### Logger (`infrastructure/logger/`)
- Zerolog implementation
- Structured logging
- Multiple output formats

#### Metrics (`infrastructure/metrics/`)
- Prometheus implementation
- Counter, histogram, and gauge metrics
- Dynamic metric registration

#### Payment (`infrastructure/payment/`)
- Payment processor client
- Mock implementation for testing
- Extensible for different processors

**Characteristics:**
- Implements domain interfaces
- External system integration
- Framework-specific code
- Most volatile layer (changes most frequently)

### 4. Presentation Layer (`pkg/googlepay/`)

HTTP handlers and API endpoints.

**Components:**
- HTTP request/response handling
- Input validation and sanitization
- Error response formatting

**Characteristics:**
- Depends on application layer
- HTTP-specific code
- Request/response transformation

## Design Patterns

### 1. Dependency Injection

All dependencies are injected through constructors for testability and flexibility.

### 2. Repository Pattern

Payment processor acts as a repository for payment operations.

### 3. Adapter Pattern

Infrastructure implementations adapt external systems to domain interfaces.

### 4. Strategy Pattern

Different implementations can be swapped without changing business logic.

## Error Handling

Domain errors are defined and wrapped with context at each layer.

## Scalability Considerations

- Stateless design for horizontal scaling
- Connection pooling
- Configurable timeouts
- Resource limits

## Security Architecture

- Token signature verification
- Secure key management
- PCI compliance considerations
- API security (authentication, rate limiting)

## Testing Strategy

- Unit tests with mocks
- Table-driven tests
- Integration tests
- High test coverage

## Monitoring and Observability

- Prometheus metrics
- Structured logging with zerolog
- Health and readiness checks
