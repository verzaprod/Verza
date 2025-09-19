package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"verza/backend/models"
	"verza/backend/pkg/blockchain"
)

// EscrowService handles escrow-related business logic
type EscrowService struct {
	db           *gorm.DB
	redisClient  *redis.Client
	hederaClient *blockchain.HederaClient
	emailService *EmailService
}

// DisputeEvidence represents evidence for escrow disputes
type DisputeEvidence struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	FileURL     string `json:"file_url,omitempty"`
	Hash        string `json:"hash,omitempty"`
}

// NewEscrowService creates a new escrow service
func NewEscrowService(db *gorm.DB, redisClient *redis.Client, hederaClient *blockchain.HederaClient, emailService *EmailService) *EscrowService {
	return &EscrowService{
		db:           db,
		redisClient:  redisClient,
		hederaClient: hederaClient,
		emailService: emailService,
	}
}

// InitiateEscrow creates a new escrow transaction and locks funds
func (es *EscrowService) InitiateEscrow(
	escrowID, verificationRequestID, payerID string,
	amount float64, currency, payerAccountID, payerPrivateKey, description string,
	autoReleaseHours *int,
) (string, error) {
	// Validate verification request exists
	var verificationRequest models.VerificationRequest
	if err := es.db.First(&verificationRequest, "id = ?", verificationRequestID).Error; err != nil {
		return "", fmt.Errorf("verification request not found: %w", err)
	}

	// Check if escrow already exists
	var existingEscrow models.EscrowTransaction
	if err := es.db.Where("verification_request_id = ?", verificationRequestID).First(&existingEscrow).Error; err == nil {
		return "", fmt.Errorf("escrow already exists for verification request %s", verificationRequestID)
	}

	// Convert amount to tinybars (1 HBAR = 100,000,000 tinybars)
	tinybars := int64(amount * 100000000)

	// Transfer funds to escrow contract
	transactionID, err := es.hederaClient.TransferHBAR(payerAccountID, payerPrivateKey, es.getEscrowContractAccountID(), tinybars)
	if err != nil {
		return "", fmt.Errorf("failed to transfer funds to escrow: %w", err)
	}

	// Calculate auto-release time
	var autoReleaseAt *time.Time
	if autoReleaseHours != nil && *autoReleaseHours > 0 {
		releaseTime := time.Now().Add(time.Duration(*autoReleaseHours) * time.Hour)
		autoReleaseAt = &releaseTime
	}

	// Create escrow record
	escrow := models.EscrowTransaction{
		ID:                    escrowID,
		VerificationRequestID: verificationRequestID,
		PayerID:               payerID,
		PayerAccountID:        payerAccountID,
		Amount:                fmt.Sprintf("%.8f", amount),
		Currency:              currency,
		Status:                "locked",
		Description:           description,
		AutoReleaseAt:         autoReleaseAt,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	if err := es.db.Create(&escrow).Error; err != nil {
		return "", fmt.Errorf("failed to create escrow record: %w", err)
	}

	// Create initial transaction record
	transactionRecord := models.EscrowTransactionRecord{
		ID:            uuid.New().String(),
		EscrowID:      escrowID,
		Type:          "lock",
		Amount:        fmt.Sprintf("%.8f", amount),
		HederaTxID:    transactionID,
		Status:        "completed",
		CreatedAt:     time.Now(),
	}

	if err := es.db.Create(&transactionRecord).Error; err != nil {
		return "", fmt.Errorf("failed to create transaction record: %w", err)
	}

	// Update verification request status
	if err := es.db.Model(&verificationRequest).Update("status", "escrow_locked").Error; err != nil {
		return "", fmt.Errorf("failed to update verification request status: %w", err)
	}

	// Cache escrow info
	es.cacheEscrowInfo(escrowID, &escrow)

	// Send notification emails if email service is available
	if es.emailService != nil {
		go es.sendEscrowNotification(escrowID, "initiated")
	}

	return transactionID, nil
}

// ReleaseEscrow releases funds to the verifier
func (es *EscrowService) ReleaseEscrow(escrowID, verifierAccountID, releaseReason string, partialAmount *float64) (string, float64, error) {
	// Get escrow
	var escrow models.EscrowTransaction
	if err := es.db.First(&escrow, "id = ?", escrowID).Error; err != nil {
		return "", 0, fmt.Errorf("escrow not found: %w", err)
	}

	// Validate escrow status
	if escrow.Status != "locked" {
		return "", 0, fmt.Errorf("escrow is not in locked status, current status: %s", escrow.Status)
	}

	// Calculate release amount
	totalAmount, err := strconv.ParseFloat(escrow.Amount, 64)
	if err != nil {
		return "", 0, fmt.Errorf("invalid escrow amount: %w", err)
	}

	releaseAmount := totalAmount
	if partialAmount != nil {
		if *partialAmount <= 0 || *partialAmount > totalAmount {
			return "", 0, fmt.Errorf("invalid partial amount: %.8f", *partialAmount)
		}
		releaseAmount = *partialAmount
	}

	// Convert to tinybars
	tinybars := int64(releaseAmount * 100000000)

	// Transfer funds to verifier
	transactionID, err := es.hederaClient.TransferHBAR(es.getEscrowContractAccountID(), "", verifierAccountID, tinybars)
	if err != nil {
		return "", 0, fmt.Errorf("failed to transfer funds to verifier: %w", err)
	}

	// Update escrow status
	newStatus := "released"
	if partialAmount != nil && *partialAmount < totalAmount {
		newStatus = "partially_released"
		// Update remaining amount
		remainingAmount := totalAmount - releaseAmount
		escrow.Amount = fmt.Sprintf("%.8f", remainingAmount)
	}

	escrow.Status = newStatus
	escrow.VerifierAccountID = &verifierAccountID
	escrow.UpdatedAt = time.Now()

	if err := es.db.Save(&escrow).Error; err != nil {
		return "", 0, fmt.Errorf("failed to update escrow: %w", err)
	}

	// Create transaction record
	transactionRecord := models.EscrowTransactionRecord{
		ID:            uuid.New().String(),
		EscrowID:      escrowID,
		Type:          "release",
		Amount:        fmt.Sprintf("%.8f", releaseAmount),
		HederaTxID:    transactionID,
		Status:        "completed",
		Metadata:      map[string]interface{}{"reason": releaseReason, "verifier_account": verifierAccountID},
		CreatedAt:     time.Now(),
	}

	if err := es.db.Create(&transactionRecord).Error; err != nil {
		return "", 0, fmt.Errorf("failed to create transaction record: %w", err)
	}

	// Update verification request status
	if err := es.db.Model(&models.VerificationRequest{}).Where("id = ?", escrow.VerificationRequestID).Update("status", "completed").Error; err != nil {
		return "", 0, fmt.Errorf("failed to update verification request status: %w", err)
	}

	// Clear cache and update
	es.clearEscrowCache(escrowID)
	es.cacheEscrowInfo(escrowID, &escrow)

	// Send notification
	if es.emailService != nil {
		go es.sendEscrowNotification(escrowID, "released")
	}

	return transactionID, releaseAmount, nil
}

// RefundEscrow refunds funds back to the payer
func (es *EscrowService) RefundEscrow(escrowID, refundReason string, partialAmount *float64) (string, float64, error) {
	// Get escrow
	var escrow models.EscrowTransaction
	if err := es.db.First(&escrow, "id = ?", escrowID).Error; err != nil {
		return "", 0, fmt.Errorf("escrow not found: %w", err)
	}

	// Validate escrow status
	if escrow.Status != "locked" && escrow.Status != "disputed" {
		return "", 0, fmt.Errorf("escrow cannot be refunded, current status: %s", escrow.Status)
	}

	// Calculate refund amount
	totalAmount, err := strconv.ParseFloat(escrow.Amount, 64)
	if err != nil {
		return "", 0, fmt.Errorf("invalid escrow amount: %w", err)
	}

	refundAmount := totalAmount
	if partialAmount != nil {
		if *partialAmount <= 0 || *partialAmount > totalAmount {
			return "", 0, fmt.Errorf("invalid partial amount: %.8f", *partialAmount)
		}
		refundAmount = *partialAmount
	}

	// Convert to tinybars
	tinybars := int64(refundAmount * 100000000)

	// Transfer funds back to payer
	transactionID, err := es.hederaClient.TransferHBAR(es.getEscrowContractAccountID(), "", escrow.PayerAccountID, tinybars)
	if err != nil {
		return "", 0, fmt.Errorf("failed to transfer refund to payer: %w", err)
	}

	// Update escrow status
	newStatus := "refunded"
	if partialAmount != nil && *partialAmount < totalAmount {
		newStatus = "partially_refunded"
		// Update remaining amount
		remainingAmount := totalAmount - refundAmount
		escrow.Amount = fmt.Sprintf("%.8f", remainingAmount)
	}

	escrow.Status = newStatus
	escrow.UpdatedAt = time.Now()

	if err := es.db.Save(&escrow).Error; err != nil {
		return "", 0, fmt.Errorf("failed to update escrow: %w", err)
	}

	// Create transaction record
	transactionRecord := models.EscrowTransactionRecord{
		ID:            uuid.New().String(),
		EscrowID:      escrowID,
		Type:          "refund",
		Amount:        fmt.Sprintf("%.8f", refundAmount),
		HederaTxID:    transactionID,
		Status:        "completed",
		Metadata:      map[string]interface{}{"reason": refundReason},
		CreatedAt:     time.Now(),
	}

	if err := es.db.Create(&transactionRecord).Error; err != nil {
		return "", 0, fmt.Errorf("failed to create transaction record: %w", err)
	}

	// Update verification request status
	if err := es.db.Model(&models.VerificationRequest{}).Where("id = ?", escrow.VerificationRequestID).Update("status", "refunded").Error; err != nil {
		return "", 0, fmt.Errorf("failed to update verification request status: %w", err)
	}

	// Clear cache and update
	es.clearEscrowCache(escrowID)
	es.cacheEscrowInfo(escrowID, &escrow)

	// Send notification
	if es.emailService != nil {
		go es.sendEscrowNotification(escrowID, "refunded")
	}

	return transactionID, refundAmount, nil
}

// CreateDispute creates a dispute for an escrow transaction
func (es *EscrowService) CreateDispute(
	escrowID, initiatedBy, reason, description string,
	evidence []DisputeEvidence,
	metadata map[string]interface{},
) (string, error) {
	// Get escrow
	var escrow models.EscrowTransaction
	if err := es.db.First(&escrow, "id = ?", escrowID).Error; err != nil {
		return "", fmt.Errorf("escrow not found: %w", err)
	}

	// Check if dispute already exists
	var existingDispute models.EscrowDispute
	if err := es.db.Where("escrow_id = ?", escrowID).First(&existingDispute).Error; err == nil {
		return "", fmt.Errorf("dispute already exists for escrow %s", escrowID)
	}

	// Validate escrow can be disputed
	if escrow.Status != "locked" && escrow.Status != "released" {
		return "", fmt.Errorf("escrow cannot be disputed, current status: %s", escrow.Status)
	}

	// Create dispute
	disputeID := uuid.New().String()
	dispute := models.EscrowDispute{
		ID:          disputeID,
		EscrowID:    escrowID,
		InitiatedBy: initiatedBy,
		Reason:      reason,
		Description: description,
		Status:      "open",
		Evidence:    evidence,
		Metadata:    metadata,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := es.db.Create(&dispute).Error; err != nil {
		return "", fmt.Errorf("failed to create dispute: %w", err)
	}

	// Update escrow status
	escrow.Status = "disputed"
	escrow.UpdatedAt = time.Now()

	if err := es.db.Save(&escrow).Error; err != nil {
		return "", fmt.Errorf("failed to update escrow status: %w", err)
	}

	// Clear cache
	es.clearEscrowCache(escrowID)

	// Send notification
	if es.emailService != nil {
		go es.sendEscrowNotification(escrowID, "disputed")
	}

	return disputeID, nil
}

// GetEscrowByID retrieves an escrow by ID
func (es *EscrowService) GetEscrowByID(escrowID string) (*models.EscrowTransaction, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("escrow:%s", escrowID)
	cachedEscrow, err := es.redisClient.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var escrow models.EscrowTransaction
		if json.Unmarshal([]byte(cachedEscrow), &escrow) == nil {
			return &escrow, nil
		}
	}

	// Get from database
	var escrow models.EscrowTransaction
	if err := es.db.First(&escrow, "id = ?", escrowID).Error; err != nil {
		return nil, fmt.Errorf("escrow not found: %w", err)
	}

	// Cache the result
	es.cacheEscrowInfo(escrowID, &escrow)

	return &escrow, nil
}

// GetEscrowWithDetails retrieves escrow with related transactions and dispute
func (es *EscrowService) GetEscrowWithDetails(escrowID string) (*models.EscrowTransaction, error) {
	var escrow models.EscrowTransaction
	if err := es.db.Preload("Transactions").Preload("Dispute").First(&escrow, "id = ?", escrowID).Error; err != nil {
		return nil, fmt.Errorf("escrow not found: %w", err)
	}

	return &escrow, nil
}

// GetEscrowByVerificationRequest retrieves escrow by verification request ID
func (es *EscrowService) GetEscrowByVerificationRequest(verificationRequestID string) (*models.EscrowTransaction, error) {
	var escrow models.EscrowTransaction
	if err := es.db.Where("verification_request_id = ?", verificationRequestID).First(&escrow).Error; err != nil {
		return nil, fmt.Errorf("escrow not found for verification request: %w", err)
	}

	return &escrow, nil
}

// ListEscrows retrieves a paginated list of escrows with filtering
func (es *EscrowService) ListEscrows(page, limit int, status, currency, payerID, verifierID string) ([]models.EscrowTransaction, int64, error) {
	offset := (page - 1) * limit

	// Build query
	query := es.db.Model(&models.EscrowTransaction{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if currency != "" {
		query = query.Where("currency = ?", currency)
	}

	if payerID != "" {
		query = query.Where("payer_id = ?", payerID)
	}

	if verifierID != "" {
		query = query.Where("verifier_account_id = ?", verifierID)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count escrows: %w", err)
	}

	// Get escrows
	var escrows []models.EscrowTransaction
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&escrows).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve escrows: %w", err)
	}

	return escrows, total, nil
}

// ValidateVerificationRequest validates that a verification request exists and belongs to the user
func (es *EscrowService) ValidateVerificationRequest(verificationRequestID, userID string) (*models.VerificationRequest, error) {
	var verificationRequest models.VerificationRequest
	if err := es.db.Where("id = ? AND user_id = ?", verificationRequestID, userID).First(&verificationRequest).Error; err != nil {
		return nil, fmt.Errorf("verification request not found or not owned by user: %w", err)
	}

	return &verificationRequest, nil
}

// GetVerificationRequestByID retrieves a verification request by ID
func (es *EscrowService) GetVerificationRequestByID(verificationRequestID string) (*models.VerificationRequest, error) {
	var verificationRequest models.VerificationRequest
	if err := es.db.First(&verificationRequest, "id = ?", verificationRequestID).Error; err != nil {
		return nil, fmt.Errorf("verification request not found: %w", err)
	}

	return &verificationRequest, nil
}

// ProcessAutoRelease processes escrows that are eligible for auto-release
func (es *EscrowService) ProcessAutoRelease() error {
	// Find escrows eligible for auto-release
	var escrows []models.EscrowTransaction
	if err := es.db.Where("status = ? AND auto_release_at IS NOT NULL AND auto_release_at <= ?", "locked", time.Now()).Find(&escrows).Error; err != nil {
		return fmt.Errorf("failed to find auto-release escrows: %w", err)
	}

	for _, escrow := range escrows {
		// Get verification request to find verifier
		verificationRequest, err := es.GetVerificationRequestByID(escrow.VerificationRequestID)
		if err != nil {
			continue
		}

		if verificationRequest.VerifierID == nil {
			// No verifier assigned, refund to payer
			_, _, err = es.RefundEscrow(escrow.ID, "auto_refund_no_verifier", nil)
		} else {
			// Release to verifier
			_, _, err = es.ReleaseEscrow(escrow.ID, *verificationRequest.VerifierAccountID, "auto_release", nil)
		}

		if err != nil {
			// Log error but continue processing other escrows
			fmt.Printf("Failed to auto-process escrow %s: %v\n", escrow.ID, err)
		}
	}

	return nil
}

// Helper methods

func (es *EscrowService) getEscrowContractAccountID() string {
	// Return the escrow contract account ID
	return "0.0.123458" // This should be configured
}

func (es *EscrowService) cacheEscrowInfo(escrowID string, escrow *models.EscrowTransaction) {
	cacheKey := fmt.Sprintf("escrow:%s", escrowID)
	escrowJSON, _ := json.Marshal(escrow)
	es.redisClient.Set(context.Background(), cacheKey, escrowJSON, 10*time.Minute)
}

func (es *EscrowService) clearEscrowCache(escrowID string) {
	cacheKey := fmt.Sprintf("escrow:%s", escrowID)
	es.redisClient.Del(context.Background(), cacheKey)
}

func (es *EscrowService) sendEscrowNotification(escrowID, eventType string) {
	if es.emailService == nil {
		return
	}

	// Get escrow details
	escrow, err := es.GetEscrowByID(escrowID)
	if err != nil {
		return
	}

	// Get verification request to get user email
	verificationRequest, err := es.GetVerificationRequestByID(escrow.VerificationRequestID)
	if err != nil {
		return
	}

	// Send appropriate notification based on event type
	switch eventType {
	case "initiated":
		es.emailService.SendTemplateEmail(
			verificationRequest.UserEmail,
			"escrow_initiated",
			"Escrow Initiated - Verification Request",
			map[string]interface{}{
				"escrow_id": escrowID,
				"amount":    escrow.Amount,
				"currency":  escrow.Currency,
			},
		)
	case "released":
		es.emailService.SendTemplateEmail(
			verificationRequest.UserEmail,
			"escrow_released",
			"Verification Complete - Funds Released",
			map[string]interface{}{
				"escrow_id": escrowID,
				"amount":    escrow.Amount,
				"currency":  escrow.Currency,
			},
		)
	case "refunded":
		es.emailService.SendTemplateEmail(
			verificationRequest.UserEmail,
			"escrow_refunded",
			"Verification Cancelled - Funds Refunded",
			map[string]interface{}{
				"escrow_id": escrowID,
				"amount":    escrow.Amount,
				"currency":  escrow.Currency,
			},
		)
	case "disputed":
		es.emailService.SendTemplateEmail(
			verificationRequest.UserEmail,
			"escrow_disputed",
			"Verification Disputed - Under Review",
			map[string]interface{}{
				"escrow_id": escrowID,
				"amount":    escrow.Amount,
				"currency":  escrow.Currency,
			},
		)
	}
}