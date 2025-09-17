package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"verza/backend/models"
	"verza/backend/pkg/blockchain"
	"verza/backend/services"
	"verza/backend/utils"
)

// VerifierHandler handles verifier-related API endpoints
type VerifierHandler struct {
	db            *gorm.DB
	hederaClient  *blockchain.HederaClient
	verifierService *services.VerifierService
	jwtService    *utils.JWTService
}

// NewVerifierHandler creates a new verifier handler
func NewVerifierHandler(db *gorm.DB, hederaClient *blockchain.HederaClient, verifierService *services.VerifierService, jwtService *utils.JWTService) *VerifierHandler {
	return &VerifierHandler{
		db:              db,
		hederaClient:    hederaClient,
		verifierService: verifierService,
		jwtService:      jwtService,
	}
}

// RegisterVerifierRequest represents the request to register a new verifier
type RegisterVerifierRequest struct {
	Name                string                 `json:"name" binding:"required,min=2,max=100"`
	Email               string                 `json:"email" binding:"required,email"`
	CompanyName         *string                `json:"company_name"`
	Website             *string                `json:"website"`
	Description         string                 `json:"description" binding:"required,min=10,max=500"`
	Specializations     []string               `json:"specializations" binding:"required,min=1"`
	Credentials         []VerifierCredential   `json:"credentials" binding:"required,min=1"`
	StakeAmount         string                 `json:"stake_amount" binding:"required"` // In HBAR
	HederaAccountID     string                 `json:"hedera_account_id" binding:"required"`
	HederaPrivateKey    string                 `json:"hedera_private_key" binding:"required"`
	BusinessLicense     *string                `json:"business_license"`
	InsuranceInfo       *InsuranceInfo         `json:"insurance_info"`
	OperatingHours      *OperatingHours        `json:"operating_hours"`
	SupportedCountries  []string               `json:"supported_countries" binding:"required,min=1"`
	PricingTiers        []PricingTier          `json:"pricing_tiers" binding:"required,min=1"`
}

// VerifierCredential represents a verifier's credential
type VerifierCredential struct {
	Type        string `json:"type" binding:"required"`
	Issuer      string `json:"issuer" binding:"required"`
	Number      string `json:"number" binding:"required"`
	IssuedDate  string `json:"issued_date" binding:"required"`
	ExpiryDate  string `json:"expiry_date" binding:"required"`
	DocumentURL string `json:"document_url" binding:"required,url"`
}

// InsuranceInfo represents verifier's insurance information
type InsuranceInfo struct {
	Provider     string `json:"provider" binding:"required"`
	PolicyNumber string `json:"policy_number" binding:"required"`
	Coverage     string `json:"coverage" binding:"required"`
	ExpiryDate   string `json:"expiry_date" binding:"required"`
}

// OperatingHours represents verifier's operating hours
type OperatingHours struct {
	Timezone string                    `json:"timezone" binding:"required"`
	Schedule map[string]DaySchedule   `json:"schedule" binding:"required"`
}

// DaySchedule represents schedule for a specific day
type DaySchedule struct {
	Open  bool   `json:"open"`
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`
}

// PricingTier represents a verifier's pricing tier
type PricingTier struct {
	ServiceType   string  `json:"service_type" binding:"required"`
	BasePrice     float64 `json:"base_price" binding:"required,min=0"`
	Currency      string  `json:"currency" binding:"required"`
	TurnaroundTime string `json:"turnaround_time" binding:"required"`
	Description   string  `json:"description"`
}

// RegisterVerifierResponse represents the response after verifier registration
type RegisterVerifierResponse struct {
	VerifierID      string `json:"verifier_id"`
	Status          string `json:"status"`
	StakeTransactionID string `json:"stake_transaction_id"`
	Message         string `json:"message"`
	NextSteps       []string `json:"next_steps"`
}

// UpdateStakeRequest represents the request to update verifier stake
type UpdateStakeRequest struct {
	Action string `json:"action" binding:"required,oneof=increase decrease"`
	Amount string `json:"amount" binding:"required"`
}

// VerifierProfileResponse represents verifier profile information
type VerifierProfileResponse struct {
	ID                  string                 `json:"id"`
	Name                string                 `json:"name"`
	Email               string                 `json:"email"`
	CompanyName         *string                `json:"company_name"`
	Website             *string                `json:"website"`
	Description         string                 `json:"description"`
	Specializations     []string               `json:"specializations"`
	Credentials         []VerifierCredential   `json:"credentials"`
	Status              string                 `json:"status"`
	ReputationScore     float64                `json:"reputation_score"`
	TotalVerifications  int                    `json:"total_verifications"`
	SuccessfulVerifications int                `json:"successful_verifications"`
	StakeAmount         string                 `json:"stake_amount"`
	AvailableStake      string                 `json:"available_stake"`
	LockedStake         string                 `json:"locked_stake"`
	Earnings            string                 `json:"earnings"`
	JoinedAt            time.Time              `json:"joined_at"`
	LastActiveAt        *time.Time             `json:"last_active_at"`
	OperatingHours      *OperatingHours        `json:"operating_hours"`
	SupportedCountries  []string               `json:"supported_countries"`
	PricingTiers        []PricingTier          `json:"pricing_tiers"`
	VerificationStats   VerificationStats      `json:"verification_stats"`
}

// VerificationStats represents verifier's verification statistics
type VerificationStats struct {
	Last30Days    StatsTimeframe `json:"last_30_days"`
	Last90Days    StatsTimeframe `json:"last_90_days"`
	AllTime       StatsTimeframe `json:"all_time"`
	AverageRating float64        `json:"average_rating"`
	ResponseTime  string         `json:"average_response_time"`
}

// StatsTimeframe represents statistics for a specific timeframe
type StatsTimeframe struct {
	TotalRequests     int     `json:"total_requests"`
	CompletedRequests int     `json:"completed_requests"`
	SuccessRate       float64 `json:"success_rate"`
	Earnings          string  `json:"earnings"`
}

// RegisterVerifier handles verifier registration
func (vh *VerifierHandler) RegisterVerifier(c *gin.Context) {
	var req RegisterVerifierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Validate stake amount
	stakeAmount, err := strconv.ParseFloat(req.StakeAmount, 64)
	if err != nil || stakeAmount < 100 { // Minimum 100 HBAR stake
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid stake amount. Minimum stake is 100 HBAR",
		})
		return
	}

	// Validate Hedera account
	if !vh.hederaClient.ValidateAccountID(req.HederaAccountID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Hedera account ID",
		})
		return
	}

	// Check if verifier already exists
	var existingVerifier models.Verifier
	result := vh.db.Where("email = ? OR hedera_account_id = ?", req.Email, req.HederaAccountID).First(&existingVerifier)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Verifier already exists with this email or Hedera account",
		})
		return
	}

	// Create verifier record
	verifierID := uuid.New().String()
	verifier := models.Verifier{
		ID:                 verifierID,
		Name:               req.Name,
		Email:              req.Email,
		CompanyName:        req.CompanyName,
		Website:            req.Website,
		Description:        req.Description,
		Specializations:    models.StringArrayJSON(req.Specializations),
		Credentials:        vh.convertCredentials(req.Credentials),
		Status:             "pending_verification",
		ReputationScore:    0.0,
		StakeAmount:        req.StakeAmount,
		HederaAccountID:    req.HederaAccountID,
		BusinessLicense:    req.BusinessLicense,
		InsuranceInfo:      vh.convertInsuranceInfo(req.InsuranceInfo),
		OperatingHours:     vh.convertOperatingHours(req.OperatingHours),
		SupportedCountries: models.StringArrayJSON(req.SupportedCountries),
		PricingTiers:       vh.convertPricingTiers(req.PricingTiers),
	}

	// Start database transaction
	tx := vh.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Save verifier to database
	if err := tx.Create(&verifier).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create verifier record",
			"details": err.Error(),
		})
		return
	}

	// Process stake transaction on Hedera
	stakeTransactionID, err := vh.verifierService.ProcessStakeTransaction(verifierID, req.HederaAccountID, req.HederaPrivateKey, stakeAmount)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process stake transaction",
			"details": err.Error(),
		})
		return
	}

	// Update verifier with stake transaction ID
	verifier.StakeTransactionID = &stakeTransactionID
	if err := tx.Save(&verifier).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update verifier with stake transaction",
		})
		return
	}

	// Commit transaction
	tx.Commit()

	// Send verification email (async)
	go vh.verifierService.SendVerificationEmail(verifier.Email, verifier.Name, verifierID)

	response := RegisterVerifierResponse{
		VerifierID:         verifierID,
		Status:             "pending_verification",
		StakeTransactionID: stakeTransactionID,
		Message:            "Verifier registration submitted successfully",
		NextSteps: []string{
			"Check your email for verification instructions",
			"Complete KYC verification process",
			"Wait for admin approval",
			"Start accepting verification requests",
		},
	}

	c.JSON(http.StatusCreated, response)
}

// GetVerifierProfile retrieves verifier profile information
func (vh *VerifierHandler) GetVerifierProfile(c *gin.Context) {
	verifierID := c.Param("id")
	if verifierID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Verifier ID is required",
		})
		return
	}

	// Get verifier from database
	var verifier models.Verifier
	if err := vh.db.First(&verifier, "id = ?", verifierID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Verifier not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve verifier",
			})
		}
		return
	}

	// Get verification statistics
	stats, err := vh.verifierService.GetVerificationStats(verifierID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve verification statistics",
		})
		return
	}

	// Get stake information from blockchain
	stakeInfo, err := vh.verifierService.GetStakeInfo(verifier.HederaAccountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve stake information",
		})
		return
	}

	response := VerifierProfileResponse{
		ID:                      verifier.ID,
		Name:                    verifier.Name,
		Email:                   verifier.Email,
		CompanyName:             verifier.CompanyName,
		Website:                 verifier.Website,
		Description:             verifier.Description,
		Specializations:         []string(verifier.Specializations),
		Credentials:             vh.convertCredentialsFromDB(verifier.Credentials),
		Status:                  verifier.Status,
		ReputationScore:         verifier.ReputationScore,
		TotalVerifications:      verifier.TotalVerifications,
		SuccessfulVerifications: verifier.SuccessfulVerifications,
		StakeAmount:             verifier.StakeAmount,
		AvailableStake:          stakeInfo.Available,
		LockedStake:             stakeInfo.Locked,
		Earnings:                verifier.TotalEarnings,
		JoinedAt:                verifier.CreatedAt,
		LastActiveAt:            verifier.LastActiveAt,
		OperatingHours:          vh.convertOperatingHoursFromDB(verifier.OperatingHours),
		SupportedCountries:      []string(verifier.SupportedCountries),
		PricingTiers:            vh.convertPricingTiersFromDB(verifier.PricingTiers),
		VerificationStats:       *stats,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateStake handles stake amount updates
func (vh *VerifierHandler) UpdateStake(c *gin.Context) {
	verifierID := c.Param("id")
	if verifierID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Verifier ID is required",
		})
		return
	}

	var req UpdateStakeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Validate amount
	amount, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil || amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid amount",
		})
		return
	}

	// Get verifier
	var verifier models.Verifier
	if err := vh.db.First(&verifier, "id = ?", verifierID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Verifier not found",
		})
		return
	}

	// Process stake update
	transactionID, newStakeAmount, err := vh.verifierService.UpdateStake(verifierID, verifier.HederaAccountID, req.Action, amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update stake",
			"details": err.Error(),
		})
		return
	}

	// Update verifier record
	verifier.StakeAmount = fmt.Sprintf("%.2f", newStakeAmount)
	if err := vh.db.Save(&verifier).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update verifier record",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":            "Stake updated successfully",
		"transaction_id":     transactionID,
		"new_stake_amount":   fmt.Sprintf("%.2f", newStakeAmount),
		"action":             req.Action,
		"amount":             req.Amount,
	})
}

// GetVerifierReputation retrieves verifier reputation information
func (vh *VerifierHandler) GetVerifierReputation(c *gin.Context) {
	verifierID := c.Param("id")
	if verifierID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Verifier ID is required",
		})
		return
	}

	reputation, err := vh.verifierService.GetReputationDetails(verifierID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve reputation details",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, reputation)
}

// ListVerifiers retrieves a list of verifiers with filtering and pagination
func (vh *VerifierHandler) ListVerifiers(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")
	specialization := c.Query("specialization")
	country := c.Query("country")
	minReputation, _ := strconv.ParseFloat(c.DefaultQuery("min_reputation", "0"), 64)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	verifiers, total, err := vh.verifierService.ListVerifiers(page, limit, status, specialization, country, minReputation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve verifiers",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"verifiers": verifiers,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// UpdateVerifierProfile updates verifier profile information
func (vh *VerifierHandler) UpdateVerifierProfile(c *gin.Context) {
	verifierID := c.Param("id")
	if verifierID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Verifier ID is required",
		})
		return
	}

	// Get current user from JWT token
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	userClaims := claims.(*utils.JWTClaims)
	if userClaims.UserID != verifierID && userClaims.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Access denied",
		})
		return
	}

	var updateReq map[string]interface{}
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Update verifier profile
	err := vh.verifierService.UpdateProfile(verifierID, updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update profile",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
	})
}

// Helper methods for data conversion

func (vh *VerifierHandler) convertCredentials(creds []VerifierCredential) models.CredentialsJSON {
	result := make(models.CredentialsJSON, len(creds))
	for i, cred := range creds {
		result[i] = models.CredentialInfo{
			Type:        cred.Type,
			Issuer:      cred.Issuer,
			Number:      cred.Number,
			IssuedDate:  cred.IssuedDate,
			ExpiryDate:  cred.ExpiryDate,
			DocumentURL: cred.DocumentURL,
		}
	}
	return result
}

func (vh *VerifierHandler) convertCredentialsFromDB(creds models.CredentialsJSON) []VerifierCredential {
	result := make([]VerifierCredential, len(creds))
	for i, cred := range creds {
		result[i] = VerifierCredential{
			Type:        cred.Type,
			Issuer:      cred.Issuer,
			Number:      cred.Number,
			IssuedDate:  cred.IssuedDate,
			ExpiryDate:  cred.ExpiryDate,
			DocumentURL: cred.DocumentURL,
		}
	}
	return result
}

func (vh *VerifierHandler) convertInsuranceInfo(info *InsuranceInfo) *models.InsuranceInfoJSON {
	if info == nil {
		return nil
	}
	return &models.InsuranceInfoJSON{
		Provider:     info.Provider,
		PolicyNumber: info.PolicyNumber,
		Coverage:     info.Coverage,
		ExpiryDate:   info.ExpiryDate,
	}
}

func (vh *VerifierHandler) convertOperatingHours(hours *OperatingHours) *models.OperatingHoursJSON {
	if hours == nil {
		return nil
	}
	schedule := make(map[string]models.DayScheduleJSON)
	for day, daySchedule := range hours.Schedule {
		schedule[day] = models.DayScheduleJSON{
			Open:  daySchedule.Open,
			Start: daySchedule.Start,
			End:   daySchedule.End,
		}
	}
	return &models.OperatingHoursJSON{
		Timezone: hours.Timezone,
		Schedule: schedule,
	}
}

func (vh *VerifierHandler) convertOperatingHoursFromDB(hours *models.OperatingHoursJSON) *OperatingHours {
	if hours == nil {
		return nil
	}
	schedule := make(map[string]DaySchedule)
	for day, daySchedule := range hours.Schedule {
		schedule[day] = DaySchedule{
			Open:  daySchedule.Open,
			Start: daySchedule.Start,
			End:   daySchedule.End,
		}
	}
	return &OperatingHours{
		Timezone: hours.Timezone,
		Schedule: schedule,
	}
}

func (vh *VerifierHandler) convertPricingTiers(tiers []PricingTier) models.PricingTiersJSON {
	result := make(models.PricingTiersJSON, len(tiers))
	for i, tier := range tiers {
		result[i] = models.PricingTierInfo{
			ServiceType:    tier.ServiceType,
			BasePrice:      tier.BasePrice,
			Currency:       tier.Currency,
			TurnaroundTime: tier.TurnaroundTime,
			Description:    tier.Description,
		}
	}
	return result
}

func (vh *VerifierHandler) convertPricingTiersFromDB(tiers models.PricingTiersJSON) []PricingTier {
	result := make([]PricingTier, len(tiers))
	for i, tier := range tiers {
		result[i] = PricingTier{
			ServiceType:    tier.ServiceType,
			BasePrice:      tier.BasePrice,
			Currency:       tier.Currency,
			TurnaroundTime: tier.TurnaroundTime,
			Description:    tier.Description,
		}
	}
	return result
}