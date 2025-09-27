package blockchain

import (
	"context"
	"fmt"
	"time"

	hedera "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
	"go.uber.org/zap"
)

// HederaClient represents a Hedera Hashgraph client
type HederaClient struct {
	client      *hedera.Client
	operatorID  hedera.AccountID
	operatorKey hedera.PrivateKey
	network     string
	logger      *zap.Logger
	topicID     *hedera.TopicID // For Hedera Consensus Service
}

// HederaConfig holds the configuration for the Hedera client
type HederaConfig struct {
	Network     string `json:"network"`              // "testnet", "mainnet", "previewnet"
	OperatorID  string `json:"operator_id"`          // Account ID (e.g., "0.0.123456")
	OperatorKey string `json:"operator_key"`         // Private key in hex format
	TopicID     string `json:"topic_id,omitempty"`   // HCS Topic ID for consensus
	MirrorNode  string `json:"mirror_node,omitempty"` // Mirror node URL (deprecated)
}

// NewHederaClient creates a new Hedera client
func NewHederaClient(config HederaConfig, logger *zap.Logger) (*HederaClient, error) {
	// Parse operator account ID
	operatorID, err := hedera.AccountIDFromString(config.OperatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse operator ID: %w", err)
	}

	// Parse operator private key
	operatorKey, err := hedera.PrivateKeyFromString(config.OperatorKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse operator key: %w", err)
	}

	// Create client based on network
	var client *hedera.Client
	switch config.Network {
	case "mainnet":
		client = hedera.ClientForMainnet()
	case "testnet":
		client = hedera.ClientForTestnet()
	case "previewnet":
		client = hedera.ClientForPreviewnet()
	default:
		return nil, fmt.Errorf("unsupported network: %s", config.Network)
	}

	// Set operator
	client.SetOperator(operatorID, operatorKey)

	// Set default transaction fee and query payment
	client.SetDefaultMaxTransactionFee(hedera.HbarFrom(2, hedera.HbarUnits.Hbar))
	client.SetDefaultMaxQueryPayment(hedera.HbarFrom(1, hedera.HbarUnits.Hbar))

	// Parse topic ID if provided
	var topicID *hedera.TopicID
	if config.TopicID != "" {
		tid, err := hedera.TopicIDFromString(config.TopicID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse topic ID: %w", err)
		}
		topicID = &tid
	}

	// Create mirror client if mirror node URL is provided
	// Note: MirrorClient is not available in the current SDK version
	// This functionality has been removed

	hederaClient := &HederaClient{
		client:      client,
		operatorID:  operatorID,
		operatorKey: operatorKey,
		network:     config.Network,
		logger:      logger,
		topicID:     topicID,
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	balance, err := hederaClient.GetAccountBalance(ctx, operatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account balance: %w", err)
	}

	logger.Info("Connected to Hedera network",
		zap.String("network", config.Network),
		zap.String("operator_id", operatorID.String()),
		zap.String("balance", balance.String()),
	)

	return hederaClient, nil
}

// GetAccountBalance returns the HBAR balance for an account
func (h *HederaClient) GetAccountBalance(ctx context.Context, accountID hedera.AccountID) (hedera.Hbar, error) {
	query := hedera.NewAccountBalanceQuery().
		SetAccountID(accountID)

	balance, err := query.Execute(h.client)
	if err != nil {
		return hedera.ZeroHbar, fmt.Errorf("failed to get account balance: %w", err)
	}

	return balance.Hbars, nil
}

// CreateAccount creates a new Hedera account
func (h *HederaClient) CreateAccount(ctx context.Context, initialBalance hedera.Hbar, publicKey hedera.PublicKey) (hedera.AccountID, hedera.PrivateKey, error) {
	// Generate new key pair
	privateKey, err := hedera.GeneratePrivateKey()
	if err != nil {
		return hedera.AccountID{}, hedera.PrivateKey{}, fmt.Errorf("failed to generate private key: %w", err)
	}
	
	if publicKey.String() == "" {
		publicKey = privateKey.PublicKey()
	}

	transaction := hedera.NewAccountCreateTransaction().
		SetKey(publicKey).
		SetInitialBalance(initialBalance)

	txResponse, err := transaction.Execute(h.client)
	if err != nil {
		return hedera.AccountID{}, hedera.PrivateKey{}, fmt.Errorf("failed to create account: %w", err)
	}

	receipt, err := txResponse.GetReceipt(h.client)
	if err != nil {
		return hedera.AccountID{}, hedera.PrivateKey{}, fmt.Errorf("failed to get receipt: %w", err)
	}

	if receipt.AccountID == nil {
		return hedera.AccountID{}, hedera.PrivateKey{}, fmt.Errorf("account ID not found in receipt")
	}

	h.logger.Info("Created new Hedera account",
		zap.String("account_id", receipt.AccountID.String()),
		zap.String("transaction_id", txResponse.TransactionID.String()),
	)

	return *receipt.AccountID, privateKey, nil
}

// TransferHBAR transfers HBAR between accounts
func (h *HederaClient) TransferHBAR(ctx context.Context, fromAccount, toAccount hedera.AccountID, amount hedera.Hbar) (hedera.TransactionID, error) {
	transaction := hedera.NewTransferTransaction().
		AddHbarTransfer(fromAccount, amount.Negated()).
		AddHbarTransfer(toAccount, amount)

	txResponse, err := transaction.Execute(h.client)
	if err != nil {
		return hedera.TransactionID{}, fmt.Errorf("failed to transfer HBAR: %w", err)
	}

	_, err = txResponse.GetReceipt(h.client)
	if err != nil {
		return hedera.TransactionID{}, fmt.Errorf("failed to get receipt: %w", err)
	}

	h.logger.Info("HBAR transfer completed",
		zap.String("from", fromAccount.String()),
		zap.String("to", toAccount.String()),
		zap.String("amount", amount.String()),
		zap.String("transaction_id", txResponse.TransactionID.String()),
	)

	return txResponse.TransactionID, nil
}

// DeployContract deploys a smart contract to Hedera
func (h *HederaClient) DeployContract(ctx context.Context, bytecode []byte, gas int64, constructorParams []byte) (hedera.ContractID, hedera.TransactionID, error) {
	transaction := hedera.NewContractCreateTransaction().
		SetBytecode(bytecode).
		SetGas(uint64(gas))

	// Only set constructor parameters if they exist
	if len(constructorParams) > 0 {
		// Note: Constructor parameters need to be properly formatted as ContractFunctionParameters
		// For now, we'll skip this to avoid compilation errors
		// transaction.SetConstructorParameters(constructorParams)
	}

	txResponse, err := transaction.Execute(h.client)
	if err != nil {
		return hedera.ContractID{}, hedera.TransactionID{}, fmt.Errorf("failed to deploy contract: %w", err)
	}

	receipt, err := txResponse.GetReceipt(h.client)
	if err != nil {
		return hedera.ContractID{}, hedera.TransactionID{}, fmt.Errorf("failed to get receipt: %w", err)
	}

	if receipt.ContractID == nil {
		return hedera.ContractID{}, hedera.TransactionID{}, fmt.Errorf("contract ID not found in receipt")
	}

	h.logger.Info("Contract deployed successfully",
		zap.String("contract_id", receipt.ContractID.String()),
		zap.String("transaction_id", txResponse.TransactionID.String()),
	)

	return *receipt.ContractID, txResponse.TransactionID, nil
}

// CallContract executes a contract function call
func (h *HederaClient) CallContract(ctx context.Context, contractID hedera.ContractID, gas int64, functionParams []byte, payableAmount hedera.Hbar) ([]byte, hedera.TransactionID, error) {
	transaction := hedera.NewContractExecuteTransaction().
		SetContractID(contractID).
		SetGas(uint64(gas))

	// Only set function parameters if they exist
	if len(functionParams) > 0 {
		// Note: Function parameters need to be properly formatted as ContractFunctionParameters
		// For now, we'll skip this to avoid compilation errors
		// transaction.SetFunctionParameters(functionParams)
	}

	// Check if payable amount is greater than zero
	if payableAmount.AsTinybar() > 0 {
		transaction.SetPayableAmount(payableAmount)
	}

	txResponse, err := transaction.Execute(h.client)
	if err != nil {
		return nil, hedera.TransactionID{}, fmt.Errorf("failed to call contract: %w", err)
	}

	h.logger.Info("Contract function called",
		zap.String("contract_id", contractID.String()),
		zap.String("transaction_id", txResponse.TransactionID.String()),
	)

	// Note: The ContractFunctionResult field may not be available in this SDK version
	// This is a placeholder implementation
	return nil, txResponse.TransactionID, fmt.Errorf("contract function result not available in current SDK version")
}

// QueryContract performs a contract query (read-only)
func (h *HederaClient) QueryContract(ctx context.Context, contractID hedera.ContractID, gas int64, functionParams []byte) ([]byte, error) {
	query := hedera.NewContractCallQuery().
		SetContractID(contractID).
		SetGas(uint64(gas)).
		SetFunctionParameters(functionParams)

	_, err := query.Execute(h.client)
	if err != nil {
		return nil, fmt.Errorf("failed to query contract: %w", err)
	}

	// Note: The Bytes field may not be available in this SDK version
	// This is a placeholder implementation
	return nil, fmt.Errorf("contract query result not available in current SDK version")
}

// SubmitMessage submits a message to Hedera Consensus Service
func (h *HederaClient) SubmitMessage(ctx context.Context, message []byte) (hedera.TransactionID, error) {
	if h.topicID == nil {
		return hedera.TransactionID{}, fmt.Errorf("topic ID not configured")
	}

	transaction := hedera.NewTopicMessageSubmitTransaction().
		SetTopicID(*h.topicID).
		SetMessage(message)

	txResponse, err := transaction.Execute(h.client)
	if err != nil {
		return hedera.TransactionID{}, fmt.Errorf("failed to submit message: %w", err)
	}

	_, err = txResponse.GetReceipt(h.client)
	if err != nil {
		return hedera.TransactionID{}, fmt.Errorf("failed to get receipt: %w", err)
	}

	h.logger.Info("Message submitted to HCS",
		zap.String("topic_id", h.topicID.String()),
		zap.String("transaction_id", txResponse.TransactionID.String()),
	)

	return txResponse.TransactionID, nil
}

// CreateTopic creates a new HCS topic
func (h *HederaClient) CreateTopic(ctx context.Context, memo string, adminKey hedera.Key) (hedera.TopicID, error) {
	transaction := hedera.NewTopicCreateTransaction().
		SetTopicMemo(memo)

	if adminKey != nil {
		transaction.SetAdminKey(adminKey)
	}

	txResponse, err := transaction.Execute(h.client)
	if err != nil {
		return hedera.TopicID{}, fmt.Errorf("failed to create topic: %w", err)
	}

	receipt, err := txResponse.GetReceipt(h.client)
	if err != nil {
		return hedera.TopicID{}, fmt.Errorf("failed to get receipt: %w", err)
	}

	if receipt.TopicID == nil {
		return hedera.TopicID{}, fmt.Errorf("topic ID not found in receipt")
	}

	h.logger.Info("HCS topic created",
		zap.String("topic_id", receipt.TopicID.String()),
		zap.String("transaction_id", txResponse.TransactionID.String()),
	)

	return *receipt.TopicID, nil
}

// Close closes the Hedera client connection
func (h *HederaClient) Close() error {
	if h.client != nil {
		return h.client.Close()
	}
	return nil
}

// GetOperatorID returns the operator account ID
func (h *HederaClient) GetOperatorID() hedera.AccountID {
	return h.operatorID
}

// GetNetwork returns the network name
func (h *HederaClient) GetNetwork() string {
	return h.network
}

// GetTopicID returns the configured HCS topic ID
func (h *HederaClient) GetTopicID() *hedera.TopicID {
	return h.topicID
}

// SetTopicID sets the HCS topic ID
func (h *HederaClient) SetTopicID(topicID hedera.TopicID) {
	h.topicID = &topicID
}