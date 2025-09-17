package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"
)

// VCRegistry represents a Verifiable Credential registry on the blockchain
type VCRegistry struct {
	client          *Client
	contractAddress common.Address
	contractABI     string
	logger          *zap.Logger
}

// VCRegistryConfig holds the configuration for the VC registry
type VCRegistryConfig struct {
	ContractAddress string `json:"contract_address"`
	ContractABI     string `json:"contract_abi"`
}

// VCStatus represents the status of a VC on the blockchain
type VCStatus struct {
	Exists    bool
	Anchored  bool
	Revoked   bool
	Timestamp *big.Int
	Issuer    common.Address
}

// VCEvent represents a VC-related event from the blockchain
type VCEvent struct {
	VCID        string
	Issuer      common.Address
	EventType   string // "anchored" or "revoked"
	Timestamp   *big.Int
	BlockNumber uint64
	TxHash      common.Hash
}

// Default VC Registry ABI (simplified version)
const DefaultVCRegistryABI = `[
	{
		"inputs": [
			{"internalType": "string", "name": "vcId", "type": "string"},
			{"internalType": "string", "name": "vcHash", "type": "string"},
			{"internalType": "string", "name": "metadata", "type": "string"}
		],
		"name": "anchorVC",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{"internalType": "string", "name": "vcId", "type": "string"},
			{"internalType": "string", "name": "reason", "type": "string"}
		],
		"name": "revokeVC",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{"internalType": "string", "name": "vcId", "type": "string"}
		],
		"name": "getVCStatus",
		"outputs": [
			{"internalType": "bool", "name": "exists", "type": "bool"},
			{"internalType": "bool", "name": "anchored", "type": "bool"},
			{"internalType": "bool", "name": "revoked", "type": "bool"},
			{"internalType": "uint256", "name": "timestamp", "type": "uint256"},
			{"internalType": "address", "name": "issuer", "type": "address"}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{"internalType": "string", "name": "vcId", "type": "string"}
		],
		"name": "isVCRevoked",
		"outputs": [
			{"internalType": "bool", "name": "", "type": "bool"}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"anonymous": false,
		"inputs": [
			{"indexed": true, "internalType": "string", "name": "vcId", "type": "string"},
			{"indexed": true, "internalType": "address", "name": "issuer", "type": "address"},
			{"indexed": false, "internalType": "string", "name": "vcHash", "type": "string"},
			{"indexed": false, "internalType": "string", "name": "metadata", "type": "string"},
			{"indexed": false, "internalType": "uint256", "name": "timestamp", "type": "uint256"}
		],
		"name": "VCAnchored",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{"indexed": true, "internalType": "string", "name": "vcId", "type": "string"},
			{"indexed": true, "internalType": "address", "name": "issuer", "type": "address"},
			{"indexed": false, "internalType": "string", "name": "reason", "type": "string"},
			{"indexed": false, "internalType": "uint256", "name": "timestamp", "type": "uint256"}
		],
		"name": "VCRevoked",
		"type": "event"
	}
]`

// NewVCRegistry creates a new VC registry instance
func NewVCRegistry(client *Client, config VCRegistryConfig, logger *zap.Logger) (*VCRegistry, error) {
	contractAddress := common.HexToAddress(config.ContractAddress)
	
	// Use default ABI if not provided
	contractABI := config.ContractABI
	if contractABI == "" {
		contractABI = DefaultVCRegistryABI
	}
	
	// Validate ABI
	_, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return nil, fmt.Errorf("invalid contract ABI: %w", err)
	}
	
	registry := &VCRegistry{
		client:          client,
		contractAddress: contractAddress,
		contractABI:     contractABI,
		logger:          logger,
	}
	
	logger.Info("VC Registry initialized",
		zap.String("contract_address", contractAddress.Hex()),
		zap.String("client_address", client.GetAddress().Hex()),
	)
	
	return registry, nil
}

// AnchorVC anchors a Verifiable Credential on the blockchain
func (r *VCRegistry) AnchorVC(ctx context.Context, vcID, vcHash, metadata string) (*types.Transaction, error) {
	r.logger.Info("Anchoring VC",
		zap.String("vc_id", vcID),
		zap.String("vc_hash", vcHash),
		zap.String("metadata", metadata),
	)
	
	// Check if VC is already anchored
	status, err := r.GetVCStatus(ctx, vcID)
	if err != nil {
		return nil, fmt.Errorf("failed to check VC status: %w", err)
	}
	
	if status.Anchored {
		return nil, fmt.Errorf("VC %s is already anchored", vcID)
	}
	
	// Send anchor transaction
	tx, err := r.client.SendContractTransaction(ctx, r.contractAddress, r.contractABI, "anchorVC", vcID, vcHash, metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to anchor VC: %w", err)
	}
	
	r.logger.Info("VC anchor transaction sent",
		zap.String("vc_id", vcID),
		zap.String("tx_hash", tx.Hash().Hex()),
	)
	
	return tx, nil
}

// RevokeVC revokes a Verifiable Credential on the blockchain
func (r *VCRegistry) RevokeVC(ctx context.Context, vcID, reason string) (*types.Transaction, error) {
	r.logger.Info("Revoking VC",
		zap.String("vc_id", vcID),
		zap.String("reason", reason),
	)
	
	// Check if VC exists and is not already revoked
	status, err := r.GetVCStatus(ctx, vcID)
	if err != nil {
		return nil, fmt.Errorf("failed to check VC status: %w", err)
	}
	
	if !status.Anchored {
		return nil, fmt.Errorf("VC %s is not anchored", vcID)
	}
	
	if status.Revoked {
		return nil, fmt.Errorf("VC %s is already revoked", vcID)
	}
	
	// Send revoke transaction
	tx, err := r.client.SendContractTransaction(ctx, r.contractAddress, r.contractABI, "revokeVC", vcID, reason)
	if err != nil {
		return nil, fmt.Errorf("failed to revoke VC: %w", err)
	}
	
	r.logger.Info("VC revoke transaction sent",
		zap.String("vc_id", vcID),
		zap.String("tx_hash", tx.Hash().Hex()),
	)
	
	return tx, nil
}

// GetVCStatus retrieves the status of a VC from the blockchain
func (r *VCRegistry) GetVCStatus(ctx context.Context, vcID string) (*VCStatus, error) {
	results, err := r.client.CallContract(ctx, r.contractAddress, r.contractABI, "getVCStatus", vcID)
	if err != nil {
		return nil, fmt.Errorf("failed to get VC status: %w", err)
	}
	
	if len(results) != 5 {
		return nil, fmt.Errorf("unexpected number of results: %d", len(results))
	}
	
	status := &VCStatus{
		Exists:    results[0].(bool),
		Anchored:  results[1].(bool),
		Revoked:   results[2].(bool),
		Timestamp: results[3].(*big.Int),
		Issuer:    results[4].(common.Address),
	}
	
	r.logger.Debug("Retrieved VC status",
		zap.String("vc_id", vcID),
		zap.Bool("exists", status.Exists),
		zap.Bool("anchored", status.Anchored),
		zap.Bool("revoked", status.Revoked),
		zap.String("timestamp", status.Timestamp.String()),
		zap.String("issuer", status.Issuer.Hex()),
	)
	
	return status, nil
}

// IsVCRevoked checks if a VC is revoked
func (r *VCRegistry) IsVCRevoked(ctx context.Context, vcID string) (bool, error) {
	results, err := r.client.CallContract(ctx, r.contractAddress, r.contractABI, "isVCRevoked", vcID)
	if err != nil {
		return false, fmt.Errorf("failed to check if VC is revoked: %w", err)
	}
	
	if len(results) != 1 {
		return false, fmt.Errorf("unexpected number of results: %d", len(results))
	}
	
	return results[0].(bool), nil
}

// WaitForAnchor waits for a VC anchor transaction to be mined and returns the receipt
func (r *VCRegistry) WaitForAnchor(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	receipt, err := r.client.WaitForTransaction(ctx, tx.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to wait for anchor transaction: %w", err)
	}
	
	if receipt.Status != 1 {
		return nil, fmt.Errorf("anchor transaction failed")
	}
	
	r.logger.Info("VC anchor transaction confirmed",
		zap.String("tx_hash", tx.Hash().Hex()),
		zap.Uint64("block_number", receipt.BlockNumber.Uint64()),
		zap.Uint64("gas_used", receipt.GasUsed),
	)
	
	return receipt, nil
}

// WaitForRevoke waits for a VC revoke transaction to be mined and returns the receipt
func (r *VCRegistry) WaitForRevoke(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	receipt, err := r.client.WaitForTransaction(ctx, tx.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to wait for revoke transaction: %w", err)
	}
	
	if receipt.Status != 1 {
		return nil, fmt.Errorf("revoke transaction failed")
	}
	
	r.logger.Info("VC revoke transaction confirmed",
		zap.String("tx_hash", tx.Hash().Hex()),
		zap.Uint64("block_number", receipt.BlockNumber.Uint64()),
		zap.Uint64("gas_used", receipt.GasUsed),
	)
	
	return receipt, nil
}

// GetContractAddress returns the contract address
func (r *VCRegistry) GetContractAddress() common.Address {
	return r.contractAddress
}

// GetClientAddress returns the client's address
func (r *VCRegistry) GetClientAddress() common.Address {
	return r.client.GetAddress()
}

// BatchAnchorVCs anchors multiple VCs in a single transaction (if supported by contract)
func (r *VCRegistry) BatchAnchorVCs(ctx context.Context, vcData []struct {
	VCID     string
	VCHash   string
	Metadata string
}) ([]*types.Transaction, error) {
	var transactions []*types.Transaction
	
	for _, vc := range vcData {
		tx, err := r.AnchorVC(ctx, vc.VCID, vc.VCHash, vc.Metadata)
		if err != nil {
			r.logger.Error("Failed to anchor VC in batch",
				zap.String("vc_id", vc.VCID),
				zap.Error(err),
			)
			continue
		}
		transactions = append(transactions, tx)
	}
	
	r.logger.Info("Batch anchor completed",
		zap.Int("total_vcs", len(vcData)),
		zap.Int("successful_anchors", len(transactions)),
	)
	
	return transactions, nil
}

// BatchRevokeVCs revokes multiple VCs
func (r *VCRegistry) BatchRevokeVCs(ctx context.Context, revokeData []struct {
	VCID   string
	Reason string
}) ([]*types.Transaction, error) {
	var transactions []*types.Transaction
	
	for _, revoke := range revokeData {
		tx, err := r.RevokeVC(ctx, revoke.VCID, revoke.Reason)
		if err != nil {
			r.logger.Error("Failed to revoke VC in batch",
				zap.String("vc_id", revoke.VCID),
				zap.Error(err),
			)
			continue
		}
		transactions = append(transactions, tx)
	}
	
	r.logger.Info("Batch revoke completed",
		zap.Int("total_vcs", len(revokeData)),
		zap.Int("successful_revokes", len(transactions)),
	)
	
	return transactions, nil
}