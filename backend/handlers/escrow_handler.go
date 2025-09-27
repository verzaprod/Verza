package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/verza/models"
	"github.com/verza/services"
)

// EscrowHandler handles escrow-related HTTP requests
type EscrowHandler struct {
	escrowService *services.EscrowService
	fraudService  *services.FraudDetectionService
}

// InitiateEscrowRequest represents the request to initiate escrow
type InitiateEscrowRequest struct {
	VerificationRequestID string  `json:"verification_request_id" binding:"required"`
	Amount                float64 `json:"amount" binding:"required,gt=0"`
	Currency              string  `json:"currency" binding:"required"`
	PayerAccountID        string  `json:"payer_account_id" binding:"required"`
	PayerPrivateKey       string  `json:"payer_private_key" binding:"required"`
	Description           string  `json:"description"`
	AutoReleaseHours      *int    `json:"auto_release_hours,omitempty"`
}

// InitiateEscrowResponse represents the response from initiating escrow
type InitiateEscrowResponse struct {
	EscrowID      string `json:"escrow_id"`
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
	CreatedAt     string `json:"created_at"`
	Message       string `json:"message"`
}

// ReleaseEscrowRequest represents the request to release escrow funds
type ReleaseEscrowRequest struct {
	EscrowID           string `json:"escrow_id" binding:"required"`
	VerifierAccountID  string `json:"verifier_account_id" binding:"required"`
	ReleaseReason      string `json:"release_reason"`
	SkipFraudCheck     bool   `json:"skip_fraud_check,omitempty"`
	PartialAmount      *float64 `json:"partial_amount,omitempty"`
}

// ReleaseEscrowResponse represents the response from releasing escrow
type ReleaseEscrowResponse struct {
	EscrowID      string `json:"escrow_id"`
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	ReleasedAmount string `json:"released_amount"`
	FraudScore    *float64 `json:"fraud_score,omitempty"`
	Message       string `json:"message"`
}

// RefundEscrowRequest represents the request to refund escrow funds
type RefundEscrowRequest struct {
	EscrowID      string  `json:"escrow_id" binding:"required"`
	RefundReason  string  `json:"refund_reason" binding:"required"`
	PartialAmount *float64 `json:"partial_amount,omitempty"`
	AdminOverride bool    `json:"admin_override,omitempty"`
}

// RefundEscrowResponse represents the response from refunding escrow
type RefundEscrowResponse struct {
	EscrowID       string `json:"escrow_id"`
	TransactionID  string `json:"transaction_id"`
	Status         string `json:"status"`
	RefundedAmount string `json:"refunded_amount"`
	Message        string `json:"message"`
}

// EscrowStatusResponse represents escrow status information
type EscrowStatusResponse struct {
	EscrowID              string                 `json:"escrow_id"`
	VerificationRequestID string                 `json:"verification_request_id"`
	Status                string                 `json:"status"`
	Amount                string                 `json:"amount"`
	Currency              string                 `json:"currency"`
	PayerAccountID        string                 `json:"payer_account_id"`
	VerifierAccountID     *string                `json:"verifier_account_id"`
	Description           string                 `json:"description"`
	CreatedAt             string                 `json:"created_at"`
	UpdatedAt             string                 `json:"updated_at"`
	AutoReleaseAt         *string                `json:"auto_release_at"`
	Transactions          []EscrowTransactionInfo `json:"transactions"`
	Dispute               *EscrowDisputeInfo     `json:"dispute,omitempty"`
}

// EscrowTransactionInfo represents transaction information
type EscrowTransactionInfo struct {
	ID            string `json:"id"`
	Type          string `json:"type"`
	Amount        string `json:"amount"`
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
}

// EscrowDisputeInfo represents dispute information
type EscrowDisputeInfo struct {
	ID          string `json:"id"`
	Reason      string `json:"reason"`
	Status      string `json:"status"`
	InitiatedBy string `json:"initiated_by"`
	CreatedAt   string `json:"created_at"`
}

// DisputeEscrowRequest represents the request to dispute escrow
type DisputeEscrowRequest struct {
	EscrowID     string                 `json:"escrow_id" binding:"required"`
	Reason       string                 `json:"reason" binding:"required"`
	Description  string                 `json:"description"`
	Evidence     []DisputeEvidence      `json:"evidence,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// DisputeEvidence represents evidence for a dispute
type DisputeEvidence struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	FileURL     string `json:"file_url,omitempty"`
	Hash        string `json:"hash,omitempty"`
}

// DisputeEscrowResponse represents the response from disputing escrow
type DisputeEscrowResponse struct {
	DisputeID string `json:"dispute_id"`
	EscrowID  string `json:"escrow_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

// EscrowListResponse represents a list of escrows
type EscrowListResponse struct {
	Escrows    []EscrowSummary `json:"escrows"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"total_pages"`
}

// EscrowSummary represents a summary of escrow information
type EscrowSummary struct {
	EscrowID              string  `json:"escrow_id"`
	VerificationRequestID string  `json:"verification_request_id"`
	Status                string  `json:"status"`
	Amount                string  `json:"amount"`
	Currency              string  `json:"currency"`
	PayerAccountID        string  `json:"payer_account_id"`
	VerifierAccountID     *string `json:"verifier_account_id"`
	CreatedAt             string  `json:"created_at"`
	AutoReleaseAt         *string `json:"auto_release_at"`
	HasDispute            bool    `json:"has_dispute"`
}

// NewEscrowHandler creates a new escrow handler
func NewEscrowHandler(escrowService *services.EscrowService, fraudService *services.FraudDetectionService) *EscrowHandler {
	return &EscrowHandler{
		escrowService: escrowService,
		fraudService:  fraudService,
	}
}

// Local response helpers to replace deprecated utils.SuccessResponse/ErrorResponse
func successResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func errorResponse(c *gin.Context, status int, message string, err interface{}) {
	resp := gin.H{
		"success": false,
		"message": message,
	}
	if err != nil {
		resp["error"] = err
	}
	c.JSON(status, resp)
}

// InitiateEscrow handles POST /escrow/initiate
func (eh *EscrowHandler) InitiateEscrow(c *gin.Context) {
	var req InitiateEscrowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		errorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	// Validate verification request exists and belongs to user
	verificationRequest, err := eh.escrowService.ValidateVerificationRequest(req.VerificationRequestID, userID.(string))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid verification request", err)
		return
	}

	// Check if escrow already exists for this verification request
	existingEscrow, err := eh.escrowService.GetEscrowByVerificationRequest(req.VerificationRequestID)
	if err == nil && existingEscrow != nil {
		errorResponse(c, http.StatusConflict, "Escrow already exists for this verification request", nil)
		return
	}

	// Create escrow
	escrowID := uuid.New().String()
	transactionID, err := eh.escrowService.InitiateEscrow(
		escrowID,
		req.VerificationRequestID,
		userID.(string),
		req.Amount,
		req.Currency,
		req.PayerAccountID,
		req.PayerPrivateKey,
		req.Description,
		req.AutoReleaseHours,
	)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to initiate escrow", err)
		return
	}

	// Get created escrow for response
	escrow, err := eh.escrowService.GetEscrowByID(escrowID)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to retrieve created escrow", err)
		return
	}

	response := InitiateEscrowResponse{
		EscrowID:      escrow.ID,
		TransactionID: transactionID,
		Status:        escrow.Status,
		Amount:        escrow.Amount,
		Currency:      escrow.Currency,
		CreatedAt:     escrow.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Message:       "Escrow initiated successfully",
	}

	successResponse(c, "Escrow initiated successfully", response)
}

// ReleaseEscrow handles POST /escrow/release
func (eh *EscrowHandler) ReleaseEscrow(c *gin.Context) {
	var req ReleaseEscrowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	// Get user role
	userRole, _ := c.Get("user_role")

	// Get escrow
	escrow, err := eh.escrowService.GetEscrowByID(req.EscrowID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Escrow not found", err)
		return
	}

	// Check authorization (admin, verifier, or system)
	if userRole != "admin" && userRole != "system" {
		// Check if user is the assigned verifier
		verificationRequest, err := eh.escrowService.GetVerificationRequestByID(escrow.VerificationRequestID)
		if err != nil || verificationRequest.VerifierID == nil || *verificationRequest.VerifierID != userID.(string) {
			utils.ErrorResponse(c, http.StatusForbidden, "Not authorized to release this escrow", nil)
			return
		}
	}

	// Run fraud detection if not skipped
	var fraudScore *float64
	if !req.SkipFraudCheck && eh.fraudService != nil {
		fraudResult, err := eh.fraudService.AnalyzeVerificationRequest(escrow.VerificationRequestID)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to run fraud detection", err)
			return
		}
		fraudScore = &fraudResult.RiskScore

		// Block release if high fraud risk (score > 80)
		if fraudResult.RiskScore > 80 {
			utils.ErrorResponse(c, http.StatusForbidden, "Release blocked due to high fraud risk", map[string]interface{}{
				"fraud_score": fraudResult.RiskScore,
				"fraud_flags": fraudResult.Flags,
			})
			return
		}
	}

	// Release escrow
	transactionID, releasedAmount, err := eh.escrowService.ReleaseEscrow(
		req.EscrowID,
		req.VerifierAccountID,
		req.ReleaseReason,
		req.PartialAmount,
	)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to release escrow", err)
		return
	}

	// Get updated escrow
	updatedEscrow, _ := eh.escrowService.GetEscrowByID(req.EscrowID)

	response := ReleaseEscrowResponse{
		EscrowID:       req.EscrowID,
		TransactionID:  transactionID,
		Status:         updatedEscrow.Status,
		ReleasedAmount: fmt.Sprintf("%.8f", releasedAmount),
		FraudScore:     fraudScore,
		Message:        "Escrow released successfully",
	}

	utils.SuccessResponse(c, "Escrow released successfully", response)
}

// RefundEscrow handles POST /escrow/refund
func (eh *EscrowHandler) RefundEscrow(c *gin.Context) {
	var req RefundEscrowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	// Get user ID and role from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	userRole, _ := c.Get("user_role")

	// Get escrow
	escrow, err := eh.escrowService.GetEscrowByID(req.EscrowID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Escrow not found", err)
		return
	}

	// Check authorization (admin, system, or payer)
	if userRole != "admin" && userRole != "system" && !req.AdminOverride {
		if escrow.PayerID != userID.(string) {
			utils.ErrorResponse(c, http.StatusForbidden, "Not authorized to refund this escrow", nil)
			return
		}
	}

	// Refund escrow
	transactionID, refundedAmount, err := eh.escrowService.RefundEscrow(
		req.EscrowID,
		req.RefundReason,
		req.PartialAmount,
	)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to refund escrow", err)
		return
	}

	// Get updated escrow
	updatedEscrow, _ := eh.escrowService.GetEscrowByID(req.EscrowID)

	response := RefundEscrowResponse{
		EscrowID:       req.EscrowID,
		TransactionID:  transactionID,
		Status:         updatedEscrow.Status,
		RefundedAmount: fmt.Sprintf("%.8f", refundedAmount),
		Message:        "Escrow refunded successfully",
	}

	utils.SuccessResponse(c, "Escrow refunded successfully", response)
}

// GetEscrowStatus handles GET /escrow/status/:id
func (eh *EscrowHandler) GetEscrowStatus(c *gin.Context) {
	escrowID := c.Param("id")
	if escrowID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Escrow ID is required", nil)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	userRole, _ := c.Get("user_role")

	// Get escrow with related data
	escrow, err := eh.escrowService.GetEscrowWithDetails(escrowID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Escrow not found", err)
		return
	}

	// Check authorization
	if userRole != "admin" && userRole != "system" {
		// Check if user is payer or verifier
		verificationRequest, err := eh.escrowService.GetVerificationRequestByID(escrow.VerificationRequestID)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get verification request", err)
			return
		}

		if escrow.PayerID != userID.(string) && 
		   (verificationRequest.VerifierID == nil || *verificationRequest.VerifierID != userID.(string)) {
			utils.ErrorResponse(c, http.StatusForbidden, "Not authorized to view this escrow", nil)
			return
		}
	}

	// Build response
	response := eh.buildEscrowStatusResponse(escrow)

	utils.SuccessResponse(c, "Escrow status retrieved successfully", response)
}

// ListEscrows handles GET /escrow/list
func (eh *EscrowHandler) ListEscrows(c *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit > 100 {
		limit = 100
	}

	// Get filter parameters
	status := c.Query("status")
	currency := c.Query("currency")
	payerID := c.Query("payer_id")
	verifierID := c.Query("verifier_id")

	// Get user ID and role from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	userRole, _ := c.Get("user_role")

	// Apply user-based filtering for non-admin users
	if userRole != "admin" && userRole != "system" {
		// Users can only see their own escrows (as payer or verifier)
		if payerID == "" && verifierID == "" {
			payerID = userID.(string)
		} else if payerID != userID.(string) && verifierID != userID.(string) {
			utils.ErrorResponse(c, http.StatusForbidden, "Not authorized to view these escrows", nil)
			return
		}
	}

	// Get escrows
	escrows, total, err := eh.escrowService.ListEscrows(page, limit, status, currency, payerID, verifierID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve escrows", err)
		return
	}

	// Convert to response format
	escrowSummaries := make([]EscrowSummary, len(escrows))
	for i, escrow := range escrows {
		escrowSummaries[i] = eh.buildEscrowSummary(&escrow)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	response := EscrowListResponse{
		Escrows:    escrowSummaries,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	utils.SuccessResponse(c, "Escrows retrieved successfully", response)
}

// DisputeEscrow handles POST /escrow/dispute
func (eh *EscrowHandler) DisputeEscrow(c *gin.Context) {
	var req DisputeEscrowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	// Get escrow
	escrow, err := eh.escrowService.GetEscrowByID(req.EscrowID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Escrow not found", err)
		return
	}

	// Check if user is authorized to dispute (payer or verifier)
	verificationRequest, err := eh.escrowService.GetVerificationRequestByID(escrow.VerificationRequestID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get verification request", err)
		return
	}

	if escrow.PayerID != userID.(string) && 
	   (verificationRequest.VerifierID == nil || *verificationRequest.VerifierID != userID.(string)) {
		utils.ErrorResponse(c, http.StatusForbidden, "Not authorized to dispute this escrow", nil)
		return
	}

	// Create dispute
	disputeID, err := eh.escrowService.CreateDispute(
		req.EscrowID,
		userID.(string),
		req.Reason,
		req.Description,
		req.Evidence,
		req.Metadata,
	)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create dispute", err)
		return
	}

	// Get updated escrow
	updatedEscrow, _ := eh.escrowService.GetEscrowByID(req.EscrowID)

	response := DisputeEscrowResponse{
		DisputeID: disputeID,
		EscrowID:  req.EscrowID,
		Status:    updatedEscrow.Status,
		Message:   "Dispute created successfully",
	}

	utils.SuccessResponse(c, "Dispute created successfully", response)
}

// Helper methods

func (eh *EscrowHandler) buildEscrowStatusResponse(escrow *models.EscrowTransaction) *EscrowStatusResponse {
	response := &EscrowStatusResponse{
		EscrowID:              escrow.ID,
		VerificationRequestID: escrow.VerificationRequestID,
		Status:                escrow.Status,
		Amount:                escrow.Amount,
		Currency:              escrow.Currency,
		PayerAccountID:        escrow.PayerAccountID,
		VerifierAccountID:     escrow.VerifierAccountID,
		Description:           escrow.Description,
		CreatedAt:             escrow.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:             escrow.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Transactions:          []EscrowTransactionInfo{},
	}

	if escrow.AutoReleaseAt != nil {
		autoReleaseAt := escrow.AutoReleaseAt.Format("2006-01-02T15:04:05Z07:00")
		response.AutoReleaseAt = &autoReleaseAt
	}

	// Add transaction history
	for _, tx := range escrow.Transactions {
		txInfo := EscrowTransactionInfo{
			ID:            tx.ID,
			Type:          tx.Type,
			Amount:        tx.Amount,
			TransactionID: tx.HederaTxID,
			Status:        tx.Status,
			CreatedAt:     tx.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		response.Transactions = append(response.Transactions, txInfo)
	}

	// Add dispute information if exists
	if escrow.Dispute != nil {
		response.Dispute = &EscrowDisputeInfo{
			ID:          escrow.Dispute.ID,
			Reason:      escrow.Dispute.Reason,
			Status:      escrow.Dispute.Status,
			InitiatedBy: escrow.Dispute.InitiatedBy,
			CreatedAt:   escrow.Dispute.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return response
}

func (eh *EscrowHandler) buildEscrowSummary(escrow *models.EscrowTransaction) EscrowSummary {
	summary := EscrowSummary{
		EscrowID:              escrow.ID,
		VerificationRequestID: escrow.VerificationRequestID,
		Status:                escrow.Status,
		Amount:                escrow.Amount,
		Currency:              escrow.Currency,
		PayerAccountID:        escrow.PayerAccountID,
		VerifierAccountID:     escrow.VerifierAccountID,
		CreatedAt:             escrow.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		HasDispute:            escrow.Dispute != nil,
	}

	if escrow.AutoReleaseAt != nil {
		autoReleaseAt := escrow.AutoReleaseAt.Format("2006-01-02T15:04:05Z07:00")
		summary.AutoReleaseAt = &autoReleaseAt
	}

	return summary
}