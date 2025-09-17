package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

// Client represents an Ethereum blockchain client
type Client struct {
	ethClient    *ethclient.Client
	privateKey   *ecdsa.PrivateKey
	publicKey    *ecdsa.PublicKey
	fromAddress  common.Address
	chainID      *big.Int
	gasLimit     uint64
	gasPrice     *big.Int
	logger       *zap.Logger
}

// Config holds the configuration for the blockchain client
type Config struct {
	RPCURL     string `json:"rpc_url"`
	PrivateKey string `json:"private_key"`
	ChainID    int64  `json:"chain_id"`
	GasLimit   uint64 `json:"gas_limit"`
	GasPrice   int64  `json:"gas_price"` // in wei
}

// NewClient creates a new blockchain client
func NewClient(config Config, logger *zap.Logger) (*Client, error) {
	// Connect to Ethereum client
	ethClient, err := ethclient.Dial(config.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}
	
	// Parse private key
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(config.PrivateKey, "0x"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to cast public key to ECDSA")
	}
	
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	chainID := big.NewInt(config.ChainID)
	
	// Set default gas price if not provided
	gasPrice := big.NewInt(config.GasPrice)
	if gasPrice.Cmp(big.NewInt(0)) == 0 {
		gasPrice = big.NewInt(20000000000) // 20 gwei default
	}
	
	// Set default gas limit if not provided
	gasLimit := config.GasLimit
	if gasLimit == 0 {
		gasLimit = 300000 // Default gas limit
	}
	
	client := &Client{
		ethClient:   ethClient,
		privateKey:  privateKey,
		publicKey:   publicKeyECDSA,
		fromAddress: fromAddress,
		chainID:     chainID,
		gasLimit:    gasLimit,
		gasPrice:    gasPrice,
		logger:      logger,
	}
	
	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	blockNumber, err := ethClient.BlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get block number: %w", err)
	}
	
	logger.Info("Connected to Ethereum network",
		zap.String("address", fromAddress.Hex()),
		zap.Int64("chain_id", config.ChainID),
		zap.Uint64("latest_block", blockNumber),
	)
	
	return client, nil
}

// GetBalance returns the balance of the client's address
func (c *Client) GetBalance(ctx context.Context) (*big.Int, error) {
	balance, err := c.ethClient.BalanceAt(ctx, c.fromAddress, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	return balance, nil
}

// GetNonce returns the current nonce for the client's address
func (c *Client) GetNonce(ctx context.Context) (uint64, error) {
	nonce, err := c.ethClient.PendingNonceAt(ctx, c.fromAddress)
	if err != nil {
		return 0, fmt.Errorf("failed to get nonce: %w", err)
	}
	return nonce, nil
}

// EstimateGas estimates the gas needed for a transaction
func (c *Client) EstimateGas(ctx context.Context, to common.Address, data []byte) (uint64, error) {
	msg := ethereum.CallMsg{
		From: c.fromAddress,
		To:   &to,
		Data: data,
	}
	
	gasLimit, err := c.ethClient.EstimateGas(ctx, msg)
	if err != nil {
		return 0, fmt.Errorf("failed to estimate gas: %w", err)
	}
	
	return gasLimit, nil
}

// SendTransaction sends a transaction to the blockchain
func (c *Client) SendTransaction(ctx context.Context, to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	nonce, err := c.GetNonce(ctx)
	if err != nil {
		return nil, err
	}
	
	// Estimate gas if data is provided
	gasLimit := c.gasLimit
	if len(data) > 0 {
		estimatedGas, err := c.EstimateGas(ctx, to, data)
		if err != nil {
			c.logger.Warn("Failed to estimate gas, using default", zap.Error(err))
		} else {
			// Add 20% buffer to estimated gas
			gasLimit = estimatedGas * 120 / 100
		}
	}
	
	tx := types.NewTransaction(nonce, to, value, gasLimit, c.gasPrice, data)
	
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(c.chainID), c.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}
	
	err = c.ethClient.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}
	
	c.logger.Info("Transaction sent",
		zap.String("tx_hash", signedTx.Hash().Hex()),
		zap.String("to", to.Hex()),
		zap.String("value", value.String()),
		zap.Uint64("gas_limit", gasLimit),
		zap.String("gas_price", c.gasPrice.String()),
	)
	
	return signedTx, nil
}

// WaitForTransaction waits for a transaction to be mined
func (c *Client) WaitForTransaction(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := c.ethClient.TransactionReceipt(ctx, txHash)
		if err == nil {
			c.logger.Info("Transaction mined",
				zap.String("tx_hash", txHash.Hex()),
				zap.Uint64("block_number", receipt.BlockNumber.Uint64()),
				zap.Uint64("gas_used", receipt.GasUsed),
				zap.Uint64("status", receipt.Status),
			)
			return receipt, nil
		}
		
		if err.Error() != "not found" {
			return nil, fmt.Errorf("failed to get transaction receipt: %w", err)
		}
		
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(2 * time.Second):
			// Continue polling
		}
	}
}

// DeployContract deploys a smart contract
func (c *Client) DeployContract(ctx context.Context, contractABI string, bytecode []byte, constructorArgs ...interface{}) (common.Address, *types.Transaction, error) {
	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to parse contract ABI: %w", err)
	}
	
	// Pack constructor arguments
	var packedArgs []byte
	if len(constructorArgs) > 0 {
		packedArgs, err = parsedABI.Pack("", constructorArgs...)
		if err != nil {
			return common.Address{}, nil, fmt.Errorf("failed to pack constructor arguments: %w", err)
		}
	}
	
	// Combine bytecode with constructor arguments
	data := append(bytecode, packedArgs...)
	
	// Deploy contract (to address is zero for contract creation)
	tx, err := c.SendTransaction(ctx, common.Address{}, big.NewInt(0), data)
	if err != nil {
		return common.Address{}, nil, err
	}
	
	// Wait for deployment
	receipt, err := c.WaitForTransaction(ctx, tx.Hash())
	if err != nil {
		return common.Address{}, nil, err
	}
	
	if receipt.Status != 1 {
		return common.Address{}, nil, fmt.Errorf("contract deployment failed")
	}
	
	c.logger.Info("Contract deployed",
		zap.String("contract_address", receipt.ContractAddress.Hex()),
		zap.String("tx_hash", tx.Hash().Hex()),
		zap.Uint64("gas_used", receipt.GasUsed),
	)
	
	return receipt.ContractAddress, tx, nil
}

// CallContract calls a read-only contract method
func (c *Client) CallContract(ctx context.Context, contractAddress common.Address, contractABI string, method string, args ...interface{}) ([]interface{}, error) {
	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract ABI: %w", err)
	}
	
	// Pack method call
	data, err := parsedABI.Pack(method, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to pack method call: %w", err)
	}
	
	// Call contract
	msg := ethereum.CallMsg{
		From: c.fromAddress,
		To:   &contractAddress,
		Data: data,
	}
	
	result, err := c.ethClient.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %w", err)
	}
	
	// Unpack result
	outputs, err := parsedABI.Unpack(method, result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack result: %w", err)
	}
	
	return outputs, nil
}

// SendContractTransaction sends a transaction to a contract method
func (c *Client) SendContractTransaction(ctx context.Context, contractAddress common.Address, contractABI string, method string, args ...interface{}) (*types.Transaction, error) {
	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract ABI: %w", err)
	}
	
	// Pack method call
	data, err := parsedABI.Pack(method, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to pack method call: %w", err)
	}
	
	// Send transaction
	tx, err := c.SendTransaction(ctx, contractAddress, big.NewInt(0), data)
	if err != nil {
		return nil, err
	}
	
	c.logger.Info("Contract transaction sent",
		zap.String("contract_address", contractAddress.Hex()),
		zap.String("method", method),
		zap.String("tx_hash", tx.Hash().Hex()),
	)
	
	return tx, nil
}

// GetTransactionByHash retrieves a transaction by its hash
func (c *Client) GetTransactionByHash(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error) {
	tx, isPending, err := c.ethClient.TransactionByHash(ctx, txHash)
	if err != nil {
		return nil, false, fmt.Errorf("failed to get transaction: %w", err)
	}
	return tx, isPending, nil
}

// GetBlockByNumber retrieves a block by its number
func (c *Client) GetBlockByNumber(ctx context.Context, blockNumber *big.Int) (*types.Block, error) {
	block, err := c.ethClient.BlockByNumber(ctx, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get block: %w", err)
	}
	return block, nil
}

// Close closes the Ethereum client connection
func (c *Client) Close() {
	c.ethClient.Close()
	c.logger.Info("Blockchain client connection closed")
}

// GetAddress returns the client's Ethereum address
func (c *Client) GetAddress() common.Address {
	return c.fromAddress
}

// GetChainID returns the chain ID
func (c *Client) GetChainID() *big.Int {
	return c.chainID
}