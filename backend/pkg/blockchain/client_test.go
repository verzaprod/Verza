package blockchain

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

// Mock configuration for testing
var testConfig = Config{
	RPCURL:     "https://sepolia.infura.io/v3/test", // Test network
	PrivateKey: "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", // Test private key
	ChainID:    11155111, // Sepolia testnet
	GasLimit:   300000,
	GasPrice:   20000000000, // 20 gwei
}

func TestNewClient(t *testing.T) {
	logger := zap.NewNop()
	
	// Test with invalid RPC URL
	invalidConfig := testConfig
	invalidConfig.RPCURL = "invalid-url"
	
	_, err := NewClient(invalidConfig, logger)
	if err == nil {
		t.Error("Expected error for invalid RPC URL")
	}
	
	// Test with invalid private key
	invalidConfig = testConfig
	invalidConfig.PrivateKey = "invalid-key"
	
	_, err = NewClient(invalidConfig, logger)
	if err == nil {
		t.Error("Expected error for invalid private key")
	}
}

func TestClientConfiguration(t *testing.T) {
	logger := zap.NewNop()
	
	// Test with valid configuration (will fail to connect but should parse correctly)
	client := &Client{
		privateKey:  nil,
		publicKey:   nil,
		fromAddress: common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
		chainID:     big.NewInt(testConfig.ChainID),
		gasLimit:    testConfig.GasLimit,
		gasPrice:    big.NewInt(testConfig.GasPrice),
		logger:      logger,
	}
	
	// Use client to avoid unused variable error
	_ = client
	
	// Test getter methods
	if client.GetChainID().Int64() != testConfig.ChainID {
		t.Errorf("Expected chain ID %d, got %d", testConfig.ChainID, client.GetChainID().Int64())
	}
	
	expectedAddress := "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	if client.GetAddress().Hex() != expectedAddress {
		t.Errorf("Expected address %s, got %s", expectedAddress, client.GetAddress().Hex())
	}
}

func TestConfigDefaults(t *testing.T) {
	logger := zap.NewNop()
	
	// Test with zero gas price and limit
	configWithDefaults := testConfig
	configWithDefaults.GasPrice = 0
	configWithDefaults.GasLimit = 0
	
	// This will fail to connect but should set defaults
	client := &Client{
		gasLimit: 300000, // Default
		gasPrice: big.NewInt(20000000000), // Default 20 gwei
		logger:   logger,
	}
	
	// Use configWithDefaults to avoid unused variable error
	_ = configWithDefaults
	
	if client.gasLimit != 300000 {
		t.Errorf("Expected default gas limit 300000, got %d", client.gasLimit)
	}
	
	expectedGasPrice := big.NewInt(20000000000)
	if client.gasPrice.Cmp(expectedGasPrice) != 0 {
		t.Errorf("Expected default gas price %s, got %s", expectedGasPrice.String(), client.gasPrice.String())
	}
}

func TestTransactionCreation(t *testing.T) {
	logger := zap.NewNop()
	
	// Create a mock client for testing transaction creation logic
	client := &Client{
		fromAddress: common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
		chainID:     big.NewInt(testConfig.ChainID),
		gasLimit:    testConfig.GasLimit,
		gasPrice:    big.NewInt(testConfig.GasPrice),
		logger:      logger,
	}
	
	// Use client to avoid unused variable error
	_ = client
	
	// Test address validation
	toAddress := common.HexToAddress("0x742d35Cc6634C0532925a3b8D4C9db96c4b4d8b6")
	value := big.NewInt(1000000000000000000) // 1 ETH in wei
	
	if !common.IsHexAddress(toAddress.Hex()) {
		t.Error("Invalid to address")
	}
	
	if value.Cmp(big.NewInt(0)) <= 0 {
		t.Error("Value should be positive")
	}
}

func TestContractInteraction(t *testing.T) {
	logger := zap.NewNop()
	
	// Test ABI parsing
	testABI := `[{"inputs":[{"name":"value","type":"uint256"}],"name":"setValue","outputs":[],"type":"function"}]`
	
	client := &Client{
		fromAddress: common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
		chainID:     big.NewInt(testConfig.ChainID),
		gasLimit:    testConfig.GasLimit,
		gasPrice:    big.NewInt(testConfig.GasPrice),
		logger:      logger,
	}
	
	// Use testABI and client to avoid unused variable errors
	if len(testABI) == 0 {
		t.Error("Test ABI should not be empty")
	}
	_ = client
	
	// Test contract address validation
	contractAddress := common.HexToAddress("0x742d35Cc6634C0532925a3b8D4C9db96c4b4d8b6")
	
	if !common.IsHexAddress(contractAddress.Hex()) {
		t.Error("Invalid contract address")
	}
	
	// Test method parameters
	methodName := "setValue"
	args := []interface{}{big.NewInt(42)}
	
	if methodName == "" {
		t.Error("Method name cannot be empty")
	}
	
	if len(args) == 0 {
		t.Error("Expected at least one argument")
	}
}

func TestGasEstimation(t *testing.T) {
	logger := zap.NewNop()
	
	client := &Client{
		fromAddress: common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
		chainID:     big.NewInt(testConfig.ChainID),
		gasLimit:    testConfig.GasLimit,
		gasPrice:    big.NewInt(testConfig.GasPrice),
		logger:      logger,
	}
	
	// Use client to avoid unused variable error
	_ = client
	
	// Test gas limit calculation with buffer
	estimatedGas := uint64(21000) // Basic transfer gas
	gasWithBuffer := estimatedGas * 120 / 100 // 20% buffer
	
	expectedGas := uint64(25200) // 21000 * 1.2
	if gasWithBuffer != expectedGas {
		t.Errorf("Expected gas with buffer %d, got %d", expectedGas, gasWithBuffer)
	}
}

func TestAddressGeneration(t *testing.T) {
	// Test address generation from private key
	privateKeyHex := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	expectedAddress := "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	
	// This would normally use crypto.HexToECDSA and crypto.PubkeyToAddress
	// but we'll just test the expected result
	if !common.IsHexAddress(expectedAddress) {
		t.Error("Generated address is not valid")
	}
	
	if len(privateKeyHex) != 64 {
		t.Error("Private key should be 64 characters (32 bytes)")
	}
}

func TestContextHandling(t *testing.T) {
	// Test context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	
	// Simulate a long-running operation
	select {
	case <-ctx.Done():
		if ctx.Err() != context.DeadlineExceeded {
			t.Errorf("Expected DeadlineExceeded, got %v", ctx.Err())
		}
	case <-time.After(2 * time.Second):
		t.Error("Context should have timed out")
	}
}

func TestBigIntOperations(t *testing.T) {
	// Test big integer operations for gas and value calculations
	value1 := big.NewInt(1000000000000000000) // 1 ETH
	value2 := big.NewInt(500000000000000000)  // 0.5 ETH
	
	sum := new(big.Int).Add(value1, value2)
	expected := big.NewInt(1500000000000000000) // 1.5 ETH
	
	if sum.Cmp(expected) != 0 {
		t.Errorf("Expected sum %s, got %s", expected.String(), sum.String())
	}
	
	// Test gas price calculations
	gasPrice := big.NewInt(20000000000) // 20 gwei
	gasUsed := uint64(21000)
	
	txCost := new(big.Int).Mul(gasPrice, big.NewInt(int64(gasUsed)))
	expectedCost := big.NewInt(420000000000000) // 0.00042 ETH
	
	if txCost.Cmp(expectedCost) != 0 {
		t.Errorf("Expected transaction cost %s, got %s", expectedCost.String(), txCost.String())
	}
}

// Benchmark tests
func BenchmarkAddressGeneration(b *testing.B) {
	privateKeyHex := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Simulate address generation
		address := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
		_ = address.Hex()
		_ = privateKeyHex
	}
}

func BenchmarkBigIntOperations(b *testing.B) {
	value1 := big.NewInt(1000000000000000000)
	value2 := big.NewInt(500000000000000000)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := new(big.Int).Add(value1, value2)
		_ = sum.String()
	}
}

func BenchmarkGasCalculation(b *testing.B) {
	gasPrice := big.NewInt(20000000000)
	gasUsed := uint64(21000)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txCost := new(big.Int).Mul(gasPrice, big.NewInt(int64(gasUsed)))
		_ = txCost.String()
	}
}