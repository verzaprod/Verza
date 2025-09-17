package blockchain

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

func TestNewVCRegistry(t *testing.T) {
	logger := zap.NewNop()
	
	// Create a mock client
	client := &Client{
		fromAddress: common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
		chainID:     big.NewInt(11155111),
		logger:      logger,
	}
	
	// Test with valid configuration
	config := VCRegistryConfig{
		ContractAddress: "0x742d35Cc6634C0532925a3b8D4C9db96c4b4d8b6",
		ContractABI:     DefaultVCRegistryABI,
	}
	
	registry, err := NewVCRegistry(client, config, logger)
	if err != nil {
		t.Fatalf("Failed to create VC registry: %v", err)
	}
	
	// Compare addresses in lowercase to handle case sensitivity
	if strings.ToLower(registry.GetContractAddress().Hex()) != strings.ToLower(config.ContractAddress) {
		t.Errorf("Expected contract address %s, got %s", config.ContractAddress, registry.GetContractAddress().Hex())
	}
	
	if registry.GetClientAddress() != client.GetAddress() {
		t.Errorf("Expected client address %s, got %s", client.GetAddress().Hex(), registry.GetClientAddress().Hex())
	}
}

func TestVCRegistryWithInvalidABI(t *testing.T) {
	logger := zap.NewNop()
	
	client := &Client{
		fromAddress: common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
		logger:      logger,
	}
	
	// Test with invalid ABI
	config := VCRegistryConfig{
		ContractAddress: "0x742d35Cc6634C0532925a3b8D4C9db96c4b4d8b6",
		ContractABI:     "invalid json",
	}
	
	_, err := NewVCRegistry(client, config, logger)
	if err == nil {
		t.Error("Expected error for invalid ABI")
	}
}

func TestVCRegistryWithDefaultABI(t *testing.T) {
	logger := zap.NewNop()
	
	client := &Client{
		fromAddress: common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
		logger:      logger,
	}
	
	// Test with empty ABI (should use default)
	config := VCRegistryConfig{
		ContractAddress: "0x742d35Cc6634C0532925a3b8D4C9db96c4b4d8b6",
		ContractABI:     "", // Empty, should use default
	}
	
	registry, err := NewVCRegistry(client, config, logger)
	if err != nil {
		t.Fatalf("Failed to create VC registry with default ABI: %v", err)
	}
	
	if registry.contractABI != DefaultVCRegistryABI {
		t.Error("Expected default ABI to be used")
	}
}

func TestDefaultVCRegistryABI(t *testing.T) {
	// Test that the default ABI is valid JSON and contains expected methods
	_, err := abi.JSON(strings.NewReader(DefaultVCRegistryABI))
	if err != nil {
		t.Fatalf("Default VC Registry ABI is invalid: %v", err)
	}
	
	// Check that ABI contains expected methods
	expectedMethods := []string{"anchorVC", "revokeVC", "getVCStatus", "isVCRevoked"}
	expectedEvents := []string{"VCAnchored", "VCRevoked"}
	
	for _, method := range expectedMethods {
		if !strings.Contains(DefaultVCRegistryABI, method) {
			t.Errorf("Default ABI missing method: %s", method)
		}
	}
	
	for _, event := range expectedEvents {
		if !strings.Contains(DefaultVCRegistryABI, event) {
			t.Errorf("Default ABI missing event: %s", event)
		}
	}
}

func TestVCStatusStruct(t *testing.T) {
	// Test VCStatus struct
	status := &VCStatus{
		Exists:    true,
		Anchored:  true,
		Revoked:   false,
		Timestamp: big.NewInt(1640995200), // 2022-01-01 00:00:00 UTC
		Issuer:    common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
	}
	
	if !status.Exists {
		t.Error("Expected VC to exist")
	}
	
	if !status.Anchored {
		t.Error("Expected VC to be anchored")
	}
	
	if status.Revoked {
		t.Error("Expected VC to not be revoked")
	}
	
	if status.Timestamp.Int64() != 1640995200 {
		t.Errorf("Expected timestamp 1640995200, got %d", status.Timestamp.Int64())
	}
	
	expectedIssuer := "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	if status.Issuer.Hex() != expectedIssuer {
		t.Errorf("Expected issuer %s, got %s", expectedIssuer, status.Issuer.Hex())
	}
}

func TestVCEventStruct(t *testing.T) {
	// Test VCEvent struct
	event := &VCEvent{
		VCID:        "vc-123",
		Issuer:      common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
		EventType:   "anchored",
		Timestamp:   big.NewInt(1640995200),
		BlockNumber: 12345678,
		TxHash:      common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
	}
	
	if event.VCID != "vc-123" {
		t.Errorf("Expected VC ID 'vc-123', got '%s'", event.VCID)
	}
	
	if event.EventType != "anchored" {
		t.Errorf("Expected event type 'anchored', got '%s'", event.EventType)
	}
	
	if event.BlockNumber != 12345678 {
		t.Errorf("Expected block number 12345678, got %d", event.BlockNumber)
	}
}

func TestBatchOperations(t *testing.T) {
	logger := zap.NewNop()
	
	// Create a mock client
	client := &Client{
		fromAddress: common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
		logger:      logger,
	}
	
	config := VCRegistryConfig{
		ContractAddress: "0x742d35Cc6634C0532925a3b8D4C9db96c4b4d8b6",
		ContractABI:     DefaultVCRegistryABI,
	}
	
	registry, err := NewVCRegistry(client, config, logger)
	if err != nil {
		t.Fatalf("Failed to create VC registry: %v", err)
	}
	
	// Use registry to avoid unused variable error
	_ = registry
	
	// Test batch anchor data structure
	batchAnchorData := []struct {
		VCID     string
		VCHash   string
		Metadata string
	}{
		{VCID: "vc-1", VCHash: "hash1", Metadata: "metadata1"},
		{VCID: "vc-2", VCHash: "hash2", Metadata: "metadata2"},
		{VCID: "vc-3", VCHash: "hash3", Metadata: "metadata3"},
	}
	
	if len(batchAnchorData) != 3 {
		t.Errorf("Expected 3 VCs in batch, got %d", len(batchAnchorData))
	}
	
	// Test batch revoke data structure
	batchRevokeData := []struct {
		VCID   string
		Reason string
	}{
		{VCID: "vc-1", Reason: "expired"},
		{VCID: "vc-2", Reason: "compromised"},
	}
	
	if len(batchRevokeData) != 2 {
		t.Errorf("Expected 2 VCs in revoke batch, got %d", len(batchRevokeData))
	}
	
	// Verify data structure
	for i, vc := range batchAnchorData {
		if vc.VCID == "" {
			t.Errorf("VC ID cannot be empty at index %d", i)
		}
		if vc.VCHash == "" {
			t.Errorf("VC Hash cannot be empty at index %d", i)
		}
	}
	
	for i, revoke := range batchRevokeData {
		if revoke.VCID == "" {
			t.Errorf("VC ID cannot be empty at index %d", i)
		}
		if revoke.Reason == "" {
			t.Errorf("Revoke reason cannot be empty at index %d", i)
		}
	}
}

func TestVCIDValidation(t *testing.T) {
	// Test VC ID validation
	validVCIDs := []string{
		"vc-123",
		"credential-456",
		"urn:uuid:12345678-1234-1234-1234-123456789012",
		"did:example:123#vc-1",
	}
	
	invalidVCIDs := []string{
		"",     // Empty
		"   ",  // Whitespace only
	}
	
	for _, vcID := range validVCIDs {
		if len(strings.TrimSpace(vcID)) == 0 {
			t.Errorf("Valid VC ID '%s' should not be empty after trimming", vcID)
		}
	}
	
	for _, vcID := range invalidVCIDs {
		if len(strings.TrimSpace(vcID)) > 0 {
			t.Errorf("Invalid VC ID '%s' should be empty after trimming", vcID)
		}
	}
}

func TestContractAddressValidation(t *testing.T) {
	// Test contract address validation
	validAddresses := []string{
		"0x742d35Cc6634C0532925a3b8D4C9db96c4b4d8b6",
		"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		"0x0000000000000000000000000000000000000000", // Zero address
	}
	
	invalidAddresses := []string{
		"0x742d35Cc6634C0532925a3b8D4C9db96c4b4d8b",   // Too short
		"0xGGGd35Cc6634C0532925a3b8D4C9db96c4b4d8b6",   // Invalid hex characters
		"",                                              // Empty
		"not-an-address",                               // Invalid format
	}
	
	for _, addr := range validAddresses {
		if !common.IsHexAddress(addr) {
			t.Errorf("Address '%s' should be valid", addr)
		}
	}
	
	for _, addr := range invalidAddresses {
		if common.IsHexAddress(addr) {
			t.Errorf("Address '%s' should be invalid", addr)
		}
	}
}

func TestTimestampHandling(t *testing.T) {
	// Test timestamp operations
	currentTime := big.NewInt(1640995200) // 2022-01-01 00:00:00 UTC
	futureTime := big.NewInt(1672531200)  // 2023-01-01 00:00:00 UTC
	pastTime := big.NewInt(1609459200)    // 2021-01-01 00:00:00 UTC
	
	// Test timestamp comparison
	if currentTime.Cmp(futureTime) >= 0 {
		t.Error("Current time should be less than future time")
	}
	
	if currentTime.Cmp(pastTime) <= 0 {
		t.Error("Current time should be greater than past time")
	}
	
	// Test timestamp difference
	diff := new(big.Int).Sub(futureTime, currentTime)
	expectedDiff := big.NewInt(31536000) // 1 year in seconds
	
	if diff.Cmp(expectedDiff) != 0 {
		t.Errorf("Expected time difference %s, got %s", expectedDiff.String(), diff.String())
	}
}

func TestErrorHandling(t *testing.T) {
	// Test error message formatting
	vcID := "vc-123"
	reason := "test error"
	
	errorMsg := fmt.Sprintf("failed to anchor VC %s: %s", vcID, reason)
	expectedMsg := "failed to anchor VC vc-123: test error"
	
	if errorMsg != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, errorMsg)
	}
	
	// Test error wrapping
	baseErr := fmt.Errorf("base error")
	wrappedErr := fmt.Errorf("wrapped: %w", baseErr)
	
	if !strings.Contains(wrappedErr.Error(), "base error") {
		t.Error("Wrapped error should contain base error message")
	}
}

// Benchmark tests
func BenchmarkVCStatusCreation(b *testing.B) {
	timestamp := big.NewInt(1640995200)
	issuer := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		status := &VCStatus{
			Exists:    true,
			Anchored:  true,
			Revoked:   false,
			Timestamp: timestamp,
			Issuer:    issuer,
		}
		_ = status
	}
}

func BenchmarkAddressConversion(b *testing.B) {
	addressStr := "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		addr := common.HexToAddress(addressStr)
		_ = addr.Hex()
	}
}

func BenchmarkBigIntComparison(b *testing.B) {
	value1 := big.NewInt(1640995200)
	value2 := big.NewInt(1672531200)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = value1.Cmp(value2)
	}
}