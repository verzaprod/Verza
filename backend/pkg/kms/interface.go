package kms

import (
	"context"
	"crypto"
	"fmt"

	"go.uber.org/zap"
)

// KMS defines the interface for key management systems
type KMS interface {
	// CreateKey creates a new key with the specified ID and type
	CreateKey(ctx context.Context, keyID string, keyType KeyType) (*KeyInfo, error)

	// GetKeyInfo retrieves information about a key
	GetKeyInfo(ctx context.Context, keyID string) (*KeyInfo, error)

	// GetPublicKey retrieves the public key for a given key ID
	GetPublicKey(ctx context.Context, keyID string) (crypto.PublicKey, error)

	// Sign signs data using the specified key
	Sign(ctx context.Context, req SignRequest) (*SignResponse, error)

	// RotateKey rotates a key to a new version
	RotateKey(ctx context.Context, keyID string) error

	// DeleteKey deletes a key
	DeleteKey(ctx context.Context, keyID string) error

	// ListKeys lists all available keys
	ListKeys(ctx context.Context) ([]string, error)
}

// Provider represents different KMS providers
type Provider string

const (
	ProviderVault Provider = "vault"
	ProviderAWS   Provider = "aws"
	ProviderLocal Provider = "local"
)

// Config holds configuration for different KMS providers
type Config struct {
	Provider Provider     `envconfig:"KMS_PROVIDER" default:"vault"`
	Vault    VaultConfig  `envconfig:"VAULT"`
	// Future: AWS, Azure, GCP configs
}

// Factory creates KMS instances based on provider
type Factory struct{}

// NewFactory creates a new KMS factory
func NewFactory() *Factory {
	return &Factory{}
}

// Create creates a KMS instance based on the provider configuration
func (f *Factory) Create(logger interface{}, config Config) (KMS, error) {
	switch config.Provider {
	case ProviderVault:
		return NewVaultKMS(logger.(*zap.Logger), config.Vault)
	case ProviderLocal:
		return NewLocalKMS(logger.(*zap.Logger))
	default:
		return nil, fmt.Errorf("unsupported KMS provider: %s", config.Provider)
	}
}