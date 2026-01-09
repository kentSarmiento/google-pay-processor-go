# Google Pay Integration Guide

## Overview

This guide explains how to integrate with Google Pay and the token decryption process.

## Google Pay Token Flow

### 1. Client-Side Integration

Frontend collects the Google Pay token using Google's JavaScript SDK.

### 2. Token Structure

The token contains:
- `signature` - ECDSA signature
- `protocolVersion` - ECv2
- `signedMessage` - Encrypted payment details
- `intermediateSigningKey` - Google's intermediate signing key

### 3. Token Decryption Process

1. **Signature Verification** - Verify using Google's root keys
2. **Hybrid Decryption** - Decrypt using ECIES with merchant private key
3. **Extract Payment Details** - Parse decrypted JSON

## Setting Up Google Pay

### 1. Register with Google Pay

Create a merchant account in the Google Pay Business Console.

### 2. Generate Keys

- Download Google's root signing keys
- Generate merchant key pair for decryption

### 3. Configure Environment

Set required environment variables for merchant ID and environment.

## Production Setup

### 1. Secure Key Storage

Use a Key Management Service (AWS KMS, Google Cloud KMS, HashiCorp Vault).

### 2. Load Google's Root Keys

Fetch and cache Google's signing keys from their public endpoint.

### 3. Implement Proper Verification

Verify signatures before decryption using Google's root keys.

## Security Best Practices

- Never log full card numbers
- Implement PCI DSS compliance
- Rotate keys regularly
- Monitor decryption failures
- Implement replay protection

## Testing

Use Google's test environment and test card numbers for development.

## Troubleshooting

Common issues:
- Key mismatch between public and private keys
- Expired Google root keys
- Wrong protocol version
- Environment mismatch (TEST vs PRODUCTION)
