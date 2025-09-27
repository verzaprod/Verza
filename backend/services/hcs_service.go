package services

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hashgraph/hedera-sdk-go/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/verza/models"
	"github.com/verza/pkg/database"
)

// HCSService handles Hedera Consensus Service operations
type HCSService struct {
	db          *gorm.DB
	redisClient *redis.Client
	hederaClient *hedera.Client
	topicID     hedera.TopicID
	operatorID  hedera.AccountID
	operatorKey hedera.PrivateKey
	config      *HCSConfig
}

// HCSConfig contains configuration for HCS service
type HCSConfig struct {
	TopicID                string        `json:"topic_id"`
	OperatorAccountID      string        `json:"operator_account_id"`
	OperatorPrivateKey     string        `json:"operator_private_key"`
	Network                string        `json:"network"` // testnet, mainnet
	MaxRetries             int           `json:"max_retries"`
	RetryDelay             time.Duration `json:"retry_delay"`
	BatchSize              int           `json:"batch_size"`
	BatchTimeout           time.Duration `json:"batch_timeout"`
	EnableCompression      bool          `json:"enable_compression"`
	EnableEncryption       bool          `json:"enable_encryption"`
	SubscriptionBufferSize int           `json:"subscription_buffer_size"`
}

// HCSMessage represents a message structure for HCS
type HCSMessage struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`
	Version       string                 `json:"version"`
	Timestamp     time.Time              `json:"timestamp"`
	Source        string                 `json:"source"`
	EventType     string                 `json:"event_type"`
	EntityID      string                 `json:"entity_id"`
	EntityType    string                 `json:"entity_type"`
	Action        string                 `json:"action"`
	Payload       map[string]interface{} `json:"payload"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	CorrelationID *string                `json:"correlation_id,omitempty"`
	SequenceNum   *int64                 `json:"sequence_num,omitempty"`
	Checksum      *string                `json:"checksum,omitempty"`
}

// HCSEventType defines different types of events
type HCSEventType string

const (
	// Verification Events
	EventVerificationRequested  HCSEventType = "verification.requested"
	EventVerificationAssigned   HCSEventType = "verification.assigned"
	EventVerificationCompleted  HCSEventType = "verification.completed"
	EventVerificationRejected   HCSEventType = "verification.rejected"
	EventVerificationCancelled  HCSEventType = "verification.cancelled"

	// Escrow Events
	EventEscrowCreated          HCSEventType = "escrow.created"
	EventEscrowLocked           HCSEventType = "escrow.locked"
	EventEscrowReleased         HCSEventType = "escrow.released"
	EventEscrowRefunded         HCSEventType = "escrow.refunded"
	EventEscrowDisputed         HCSEventType = "escrow.disputed"

	// Verifier Events
	EventVerifierRegistered     HCSEventType = "verifier.registered"
	EventVerifierApproved       HCSEventType = "verifier.approved"
	EventVerifierSuspended      HCSEventType = "verifier.suspended"
	EventVerifierStakeUpdated   HCSEventType = "verifier.stake_updated"
	EventVerifierReputationChanged HCSEventType = "verifier.reputation_changed"

	// Fraud Detection Events
	EventFraudDetected          HCSEventType = "fraud.detected"
	EventFraudAnalysisCompleted HCSEventType = "fraud.analysis_completed"
	EventFraudPatternIdentified HCSEventType = "fraud.pattern_identified"

	// Credential Events
	EventCredentialIssued       HCSEventType = "credential.issued"
	EventCredentialRevoked      HCSEventType = "credential.revoked"
	EventCredentialSuspended    HCSEventType = "credential.suspended"
	EventCredentialReactivated  HCSEventType = "credential.reactivated"

	// DID Events
	EventDIDRegistered          HCSEventType = "did.registered"
	EventDIDUpdated             HCSEventType = "did.updated"
	EventDIDDeactivated         HCSEventType = "did.deactivated"
)

// HCSSubscription represents an active subscription
type HCSSubscription struct {
	ID          string
	TopicID     hedera.TopicID
	StartTime   *time.Time
	EndTime     *time.Time
	Handler     func(*HCSMessage) error
	ErrorHandler func(error)
	Active      bool
	Cancel      context.CancelFunc
}

// HCSMessageRecord stores HCS messages in database
type HCSMessageRecord struct {
	ID              string                 `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TopicID         string                 `json:"topic_id" gorm:"not null;index"`
	SequenceNumber  int64                  `json:"sequence_number" gorm:"not null;index"`
	ConsensusTime   time.Time              `json:"consensus_time" gorm:"not null;index"`
	MessageType     string                 `json:"message_type" gorm:"not null;index"`
	EventType       string                 `json:"event_type" gorm:"not null;index"`
	EntityID        string                 `json:"entity_id" gorm:"not null;index"`
	EntityType      string                 `json:"entity_type" gorm:"not null;index"`
	Action          string                 `json:"action" gorm:"not null"`
	Payload         models.JSONField       `json:"payload" gorm:"type:jsonb"`
	Metadata        models.JSONField       `json:"metadata" gorm:"type:jsonb"`
	RawMessage      []byte                 `json:"raw_message"`
	MessageHash     string                 `json:"message_hash" gorm:"not null;index"`
	Processed       bool                   `json:"processed" gorm:"default:false;index"`
	ProcessedAt     *time.Time             `json:"processed_at"`
	ProcessingError *string                `json:"processing_error"`
	RetryCount      int                    `json:"retry_count" gorm:"default:0"`
	CorrelationID   *string                `json:"correlation_id" gorm:"index"`
	CreatedAt       time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
}

// NewHCSService creates a new HCS service instance
func NewHCSService(db *gorm.DB, redisClient *redis.Client, config *HCSConfig) (*HCSService, error) {
	// Parse Hedera configuration
	topicID, err := hedera.TopicIDFromString(config.TopicID)
	if err != nil {
		return nil, fmt.Errorf("invalid topic ID: %w", err)
	}

	operatorID, err := hedera.AccountIDFromString(config.OperatorAccountID)
	if err != nil {
		return nil, fmt.Errorf("invalid operator account ID: %w", err)
	}

	operatorKey, err := hedera.PrivateKeyFromString(config.OperatorPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("invalid operator private key: %w", err)
	}

	// Create Hedera client
	var client *hedera.Client
	switch config.Network {
	case "testnet":
		client = hedera.ClientForTestnet()
	case "mainnet":
		client = hedera.ClientForMainnet()
	default:
		return nil, fmt.Errorf("unsupported network: %s", config.Network)
	}

	client.SetOperator(operatorID, operatorKey)

	// Set default configuration values
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.RetryDelay == 0 {
		config.RetryDelay = time.Second * 2
	}
	if config.BatchSize == 0 {
		config.BatchSize = 100
	}
	if config.BatchTimeout == 0 {
		config.BatchTimeout = time.Second * 30
	}
	if config.SubscriptionBufferSize == 0 {
		config.SubscriptionBufferSize = 1000
	}

	return &HCSService{
		db:           db,
		redisClient:  redisClient,
		hederaClient: client,
		topicID:      topicID,
		operatorID:   operatorID,
		operatorKey:  operatorKey,
		config:       config,
	}, nil
}

// PublishEvent publishes an event to HCS
func (hcs *HCSService) PublishEvent(eventType HCSEventType, entityID, entityType, action string, payload map[string]interface{}, metadata map[string]interface{}) error {
	message := &HCSMessage{
		ID:         uuid.New().String(),
		Type:       "verza.event",
		Version:    "1.0",
		Timestamp:  time.Now().UTC(),
		Source:     "verza-backend",
		EventType:  string(eventType),
		EntityID:   entityID,
		EntityType: entityType,
		Action:     action,
		Payload:    payload,
		Metadata:   metadata,
	}

	return hcs.publishMessage(message)
}

// PublishVerificationEvent publishes verification-related events
func (hcs *HCSService) PublishVerificationEvent(eventType HCSEventType, verificationID string, payload map[string]interface{}) error {
	return hcs.PublishEvent(eventType, verificationID, "verification", string(eventType), payload, nil)
}

// PublishEscrowEvent publishes escrow-related events
func (hcs *HCSService) PublishEscrowEvent(eventType HCSEventType, escrowID string, payload map[string]interface{}) error {
	return hcs.PublishEvent(eventType, escrowID, "escrow", string(eventType), payload, nil)
}

// PublishVerifierEvent publishes verifier-related events
func (hcs *HCSService) PublishVerifierEvent(eventType HCSEventType, verifierID string, payload map[string]interface{}) error {
	return hcs.PublishEvent(eventType, verifierID, "verifier", string(eventType), payload, nil)
}

// PublishFraudEvent publishes fraud detection events
func (hcs *HCSService) PublishFraudEvent(eventType HCSEventType, entityID string, payload map[string]interface{}) error {
	return hcs.PublishEvent(eventType, entityID, "fraud_detection", string(eventType), payload, nil)
}

// PublishCredentialEvent publishes credential-related events
func (hcs *HCSService) PublishCredentialEvent(eventType HCSEventType, credentialID string, payload map[string]interface{}) error {
	return hcs.PublishEvent(eventType, credentialID, "credential", string(eventType), payload, nil)
}

// PublishDIDEvent publishes DID-related events
func (hcs *HCSService) PublishDIDEvent(eventType HCSEventType, didID string, payload map[string]interface{}) error {
	return hcs.PublishEvent(eventType, didID, "did", string(eventType), payload, nil)
}

// publishMessage sends a message to HCS topic
func (hcs *HCSService) publishMessage(message *HCSMessage) error {
	// Serialize message
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to serialize message: %w", err)
	}

	// Add checksum
	checksum := database.CalculateChecksum(messageBytes)
	message.Checksum = &checksum

	// Re-serialize with checksum
	messageBytes, err = json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to serialize message with checksum: %w", err)
	}

	// Compress if enabled
	if hcs.config.EnableCompression {
		messageBytes, err = gzipCompress(messageBytes)
		if err != nil {
			return fmt.Errorf("failed to compress message: %w", err)
		}
	}

	// Encrypt if enabled
	if hcs.config.EnableEncryption {
		messageBytes, err = passthroughEncrypt(messageBytes)
		if err != nil {
			return fmt.Errorf("failed to encrypt message: %w", err)
		}
	}

	// Create HCS submit message transaction
	transaction := hedera.NewTopicMessageSubmitTransaction().
		SetTopicID(hcs.topicID).
		SetMessage(messageBytes)

	// Execute transaction with retries
	var txResponse hedera.TransactionResponse
	for attempt := 0; attempt <= hcs.config.MaxRetries; attempt++ {
		txResponse, err = transaction.Execute(hcs.hederaClient)
		if err == nil {
			break
		}
		if attempt < hcs.config.MaxRetries {
			time.Sleep(hcs.config.RetryDelay * time.Duration(attempt+1))
		}
	}

	if err != nil {
		return fmt.Errorf("failed to submit message to HCS after %d attempts: %w", hcs.config.MaxRetries+1, err)
	}

	// Get receipt
	receipt, err := txResponse.GetReceipt(hcs.hederaClient)
	if err != nil {
		return fmt.Errorf("failed to get transaction receipt: %w", err)
	}

	if receipt.Status != hedera.StatusSuccess {
		return fmt.Errorf("transaction failed with status: %s", receipt.Status)
	}

	// Cache message for quick retrieval
	cacheKey := fmt.Sprintf("hcs:message:%s", message.ID)
	cacheData, _ := json.Marshal(message)
	hcs.redisClient.Set(context.Background(), cacheKey, cacheData, time.Hour*24)

	log.Printf("Successfully published HCS message: %s (Type: %s, Entity: %s)", message.ID, message.EventType, message.EntityID)
	return nil
}

// Subscribe creates a subscription to HCS topic
func (hcs *HCSService) Subscribe(startTime *time.Time, handler func(*HCSMessage) error, errorHandler func(error)) (*HCSSubscription, error) {
	ctx, cancel := context.WithCancel(context.Background())

	subscription := &HCSSubscription{
		ID:           uuid.New().String(),
		TopicID:      hcs.topicID,
		StartTime:    startTime,
		Handler:      handler,
		ErrorHandler: errorHandler,
		Active:       true,
		Cancel:       cancel,
	}

	// Create HCS subscription query
	query := hedera.NewTopicMessageQuery().
		SetTopicID(hcs.topicID)

	if startTime != nil {
		query = query.SetStartTime(*startTime)
	}

	// Start subscription in goroutine
	go func() {
		defer func() {
			subscription.Active = false
			if r := recover(); r != nil {
				log.Printf("HCS subscription panic recovered: %v", r)
				if errorHandler != nil {
					errorHandler(fmt.Errorf("subscription panic: %v", r))
				}
			}
		}()

		_, err := query.Subscribe(hcs.hederaClient, func(message hedera.TopicMessage) {
			select {
			case <-ctx.Done():
				return
			default:
			}

			// Process message
			if err := hcs.processTopicMessage(message, handler); err != nil {
				log.Printf("Error processing HCS message: %v", err)
				if errorHandler != nil {
					errorHandler(err)
				}
			}
		})

		if err != nil {
			log.Printf("HCS subscription error: %v", err)
			if errorHandler != nil {
				errorHandler(err)
			}
		}
	}()

	return subscription, nil
}

// processTopicMessage processes incoming HCS messages
func (hcs *HCSService) processTopicMessage(topicMessage hedera.TopicMessage, handler func(*HCSMessage) error) error {
	messageBytes := topicMessage.Contents

	// Decrypt if enabled
	if hcs.config.EnableEncryption {
		decrypted, err := passthroughDecrypt(messageBytes)
		if err != nil {
			return fmt.Errorf("failed to decrypt message: %w", err)
		}
		messageBytes = decrypted
	}

	// Decompress if enabled
	if hcs.config.EnableCompression {
		decompressed, err := gzipDecompress(messageBytes)
		if err != nil {
			return fmt.Errorf("failed to decompress message: %w", err)
		}
		messageBytes = decompressed
	}

	// Parse message
	var message HCSMessage
	if err := json.Unmarshal(messageBytes, &message); err != nil {
		return fmt.Errorf("failed to parse HCS message: %w", err)
	}

	// Verify checksum if present
	if message.Checksum != nil {
		// Remove checksum for verification
		checksum := *message.Checksum
		message.Checksum = nil
		verifyBytes, _ := json.Marshal(message)
		expectedChecksum := database.CalculateChecksum(verifyBytes)
		if checksum != expectedChecksum {
			return fmt.Errorf("message checksum verification failed")
		}
		message.Checksum = &checksum
	}

	// Set sequence number from topic message
	message.SequenceNum = &topicMessage.SequenceNumber

	// Store message in database
	if err := hcs.storeMessage(&message, topicMessage, messageBytes); err != nil {
		log.Printf("Failed to store HCS message: %v", err)
	}

	// Call handler
	if handler != nil {
		return handler(&message)
	}

	return nil
}

// storeMessage stores HCS message in database
func (hcs *HCSService) storeMessage(message *HCSMessage, topicMessage hedera.TopicMessage, rawMessage []byte) error {
	messageHash := database.CalculateChecksum(rawMessage)

	record := &HCSMessageRecord{
		TopicID:        hcs.topicID.String(),
		SequenceNumber: topicMessage.SequenceNumber,
		ConsensusTime:  topicMessage.ConsensusTimestamp,
		MessageType:    message.Type,
		EventType:      message.EventType,
		EntityID:       message.EntityID,
		EntityType:     message.EntityType,
		Action:         message.Action,
		Payload:        models.JSONField(message.Payload),
		Metadata:       models.JSONField(message.Metadata),
		RawMessage:     rawMessage,
		MessageHash:    messageHash,
		CorrelationID:  message.CorrelationID,
	}

	return hcs.db.Create(record).Error
}

// GetMessageHistory retrieves message history from database
func (hcs *HCSService) GetMessageHistory(entityID, entityType string, eventTypes []string, limit, offset int) ([]HCSMessageRecord, error) {
	var messages []HCSMessageRecord

	query := hcs.db.Model(&HCSMessageRecord{}).Where("entity_id = ? AND entity_type = ?", entityID, entityType)

	if len(eventTypes) > 0 {
		query = query.Where("event_type IN ?", eventTypes)
	}

	err := query.Order("consensus_time DESC").Limit(limit).Offset(offset).Find(&messages).Error
	return messages, err
}

// GetMessageBySequence retrieves a specific message by sequence number
func (hcs *HCSService) GetMessageBySequence(sequenceNumber int64) (*HCSMessageRecord, error) {
	var message HCSMessageRecord
	err := hcs.db.Where("sequence_number = ?", sequenceNumber).First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// GetLatestSequenceNumber gets the latest sequence number processed
func (hcs *HCSService) GetLatestSequenceNumber() (int64, error) {
	var result struct {
		MaxSequence int64
	}

	err := hcs.db.Model(&HCSMessageRecord{}).
		Select("COALESCE(MAX(sequence_number), 0) as max_sequence").
		Scan(&result).Error

	return result.MaxSequence, err
}

// MarkMessageProcessed marks a message as processed
func (hcs *HCSService) MarkMessageProcessed(messageID string, processingError *string) error {
	updates := map[string]interface{}{
		"processed":    true,
		"processed_at": time.Now(),
	}

	if processingError != nil {
		updates["processing_error"] = *processingError
		updates["retry_count"] = gorm.Expr("retry_count + 1")
	}

	return hcs.db.Model(&HCSMessageRecord{}).Where("id = ?", messageID).Updates(updates).Error
}

// GetUnprocessedMessages retrieves unprocessed messages
func (hcs *HCSService) GetUnprocessedMessages(limit int) ([]HCSMessageRecord, error) {
	var messages []HCSMessageRecord
	err := hcs.db.Where("processed = false AND retry_count < ?", 5).
		Order("consensus_time ASC").
		Limit(limit).
		Find(&messages).Error
	return messages, err
}

// Close closes the HCS service and cleans up resources
func (hcs *HCSService) Close() error {
	if hcs.hederaClient != nil {
		return hcs.hederaClient.Close()
	}
	return nil
}

// TableName returns the table name for HCSMessageRecord
func (HCSMessageRecord) TableName() string {
	return "hcs_message_records"
}

// gzipCompress compresses data using gzip
func gzipCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	_, err := w.Write(data)
	if err != nil {
		_ = w.Close()
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// gzipDecompress decompresses gzip data
func gzipDecompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	decompressed, err := io.ReadAll(r)
	_ = r.Close()
	if err != nil {
		return nil, err
	}
	return decompressed, nil
}

// passthroughEncrypt is a placeholder for encryption (no-op)
func passthroughEncrypt(data []byte) ([]byte, error) {
	return data, nil
}

// passthroughDecrypt is a placeholder for decryption (no-op)
func passthroughDecrypt(data []byte) ([]byte, error) {
	return data, nil
}