package blockchain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// ContractConfig represents the structure of the contract configuration file
type ContractConfig struct {
	Networks map[string]NetworkConfig `json:"networks"`
}

// NetworkConfig represents the configuration for a specific network
type NetworkConfig struct {
	ChainID            int    `json:"chainId"`
	VCRegistry         string `json:"vcRegistry"`
	EscrowContract     string `json:"escrowContract"`
	VerifierMarketplace string `json:"verifierMarketplace"`
	DeployedAt         string `json:"deployedAt"`
	Deployer           string `json:"deployer"`
	AccountID          string `json:"accountId"`
}

// LoadContractConfig loads the contract configuration from the specified file
func LoadContractConfig(configPath string) (*ContractConfig, error) {
	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("contract config file not found: %s", configPath)
	}

	// Read file
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read contract config file: %w", err)
	}

	// Parse JSON
	var config ContractConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse contract config file: %w", err)
	}

	return &config, nil
}

// GetNetworkConfig returns the configuration for the specified network
func (c *ContractConfig) GetNetworkConfig(network string) (*NetworkConfig, error) {
	if config, ok := c.Networks[network]; ok {
		return &config, nil
	}
	return nil, fmt.Errorf("network configuration not found: %s", network)
}

// DefaultConfigPath returns the default path to the contract configuration file
func DefaultConfigPath() string {
	// Get the project root directory
	rootDir := os.Getenv("VERZA_ROOT_DIR")
	if rootDir == "" {
		// Fallback to relative path if environment variable is not set
		return filepath.Join("..", "..", "contracts", "contract-config.json")
	}
	return filepath.Join(rootDir, "contracts", "contract-config.json")
}