package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"verza/backend/models"
	"verza/backend/pkg/blockchain"
	"verza/backend/utils"
)

// VerifierService handles verifier-related business logic
type VerifierService struct {
	db           *gorm.DB
	redisClient  *redis.Client
	hederaClient *blockchain.HederaClient
	emailService *EmailService
}

// StakeInfo represents stake information from blockchain
type StakeInfo struct {
	Total     string `json:"total"`
	Available string `json:"available"`
	Locked    string `json:"locked"`
	Pending   string `json:"pending"`
}

// ReputationDetails represents detailed reputation information
type ReputationDetails struct {
	VerifierID      string                 `json:"verifier_id"`
	OverallScore    float64                `json:"overall_score"`
	ScoreBreakdown  ReputationBreakdown    `json:"score_breakdown"`
	RecentActivity  []ReputationActivity   `json:"recent_activity"`
	Trends          ReputationTrends       `json:"trends"`
	Comparisons     ReputationComparisons  `json:"comparisons"`
	LastUpdated     time.Time              `json:"last_updated"`
}

// ReputationBreakdown breaks down reputation score by components
type ReputationBreakdown struct {
	Accuracy        float64 `json:"accuracy"`
	Speed           float64 `json:"speed"`
	Professionalism float64 `json:"professionalism"`
	Reliability     float64 `json:"reliability"`
	CustomerService float64 `json:"customer_service"`
}

// ReputationActivity represents recent reputation-affecting activities
type ReputationActivity struct {
	Date        time.Time `json:"date"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Impact      float64   `json:"impact"`
	RequestID   string    `json:"request_id,omitempty"`
}

// ReputationTrends represents reputation trends over time
type ReputationTrends struct {
	Last7Days   float64 `json:"last_7_days"`
	Last30Days  float64 `json:"last_30_days"`
	Last90Days  float64 `json:"last_90_days"`
	YearToDate  float64 `json:"year_to_date"`
}

// ReputationComparisons compares verifier to peers
type ReputationComparisons struct {
	IndustryAverage    float64 `json:"industry_average"`
	SpecializationAvg  float64 `json:"specialization_average"`
	TopPercentile      float64 `json:"top_percentile"`
	RankInCategory     int     `json:"rank_in_category"`
	TotalInCategory    int     `json:"total_in_category"`
}

// VerificationStats represents verification statistics
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

// VerifierListItem represents a verifier in list responses
type VerifierListItem struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name"`
	CompanyName         *string   `json:"company_name"`
	Specializations     []string  `json:"specializations"`
	Status              string    `json:"status"`
	ReputationScore     float64   `json:"reputation_score"`
	TotalVerifications  int       `json:"total_verifications"`
	SuccessRate         float64   `json:"success_rate"`
	AverageResponseTime string    `json:"average_response_time"`
	SupportedCountries  []string  `json:"supported_countries"`
	JoinedAt            time.Time `json:"joined_at"`
	LastActiveAt        *time.Time `json:"last_active_at"`
	IsOnline            bool      `json:"is_online"`
}

// NewVerifierService creates a new verifier service
func NewVerifierService(db *gorm.DB, redisClient *redis.Client, hederaClient *blockchain.HederaClient, emailService *EmailService) *VerifierService {
	return &VerifierService{
		db:           db,
		redisClient:  redisClient,
		hederaClient: hederaClient,
		emailService: emailService,
	}
}

// ProcessStakeTransaction processes the initial stake transaction for a verifier
func (vs *VerifierService) ProcessStakeTransaction(verifierID, accountID, privateKey string, amount float64) (string, error) {
	// Convert amount to tinybars (1 HBAR = 100,000,000 tinybars)
	tinybars := int64(amount * 100000000)

	// Create stake transaction on Hedera
	transactionID, err := vs.hederaClient.TransferHBAR(accountID, privateKey, vs.getMarketplaceAccountID(), tinybars)
	if err != nil {
		return "", fmt.Errorf("failed to transfer stake to marketplace: %w", err)
	}

	// Record stake transaction in database
	stakeRecord := models.VerifierStake{
		VerifierID:    verifierID,
		TransactionID: transactionID,
		Amount:        fmt.Sprintf("%.8f", amount),
		Type:          "initial_stake",
		Status:        "completed",
		HederaTxID:    transactionID,
	}

	if err := vs.db.Create(&stakeRecord).Error; err != nil {
		return "", fmt.Errorf("failed to record stake transaction: %w", err)
	}

	// Update verifier's stake in smart contract
	err = vs.updateVerifierStakeOnChain(verifierID, accountID, amount)
	if err != nil {
		return "", fmt.Errorf("failed to update stake on smart contract: %w", err)
	}

	return transactionID, nil
}

// UpdateStake updates a verifier's stake amount
func (vs *VerifierService) UpdateStake(verifierID, accountID, action string, amount float64) (string, float64, error) {
	// Get current verifier
	var verifier models.Verifier
	if err := vs.db.First(&verifier, "id = ?", verifierID).Error; err != nil {
		return "", 0, fmt.Errorf("verifier not found: %w", err)
	}

	currentStake, err := strconv.ParseFloat(verifier.StakeAmount, 64)
	if err != nil {
		return "", 0, fmt.Errorf("invalid current stake amount: %w", err)
	}

	var newStakeAmount float64
	var transactionType string
	var transactionID string

	switch action {
	case "increase":
		// Transfer additional HBAR to marketplace
		tinybars := int64(amount * 100000000)
		transactionID, err = vs.hederaClient.TransferHBAR(accountID, "", vs.getMarketplaceAccountID(), tinybars)
		if err != nil {
			return "", 0, fmt.Errorf("failed to transfer additional stake: %w", err)
		}
		newStakeAmount = currentStake + amount
		transactionType = "stake_increase"

	case "decrease":
		// Check if verifier has enough available stake
		stakeInfo, err := vs.GetStakeInfo(accountID)
		if err != nil {
			return "", 0, fmt.Errorf("failed to get stake info: %w", err)
		}

		availableStake, err := strconv.ParseFloat(stakeInfo.Available, 64)
		if err != nil {
			return "", 0, fmt.Errorf("invalid available stake: %w", err)
		}

		if amount > availableStake {
			return "", 0, fmt.Errorf("insufficient available stake for withdrawal")
		}

		// Check minimum stake requirement
		if currentStake-amount < 100 {
			return "", 0, fmt.Errorf("cannot reduce stake below minimum requirement of 100 HBAR")
		}

		// Transfer HBAR back to verifier
		tinybars := int64(amount * 100000000)
		transactionID, err = vs.hederaClient.TransferHBAR(vs.getMarketplaceAccountID(), "", accountID, tinybars)
		if err != nil {
			return "", 0, fmt.Errorf("failed to transfer stake back to verifier: %w", err)
		}
		newStakeAmount = currentStake - amount
		transactionType = "stake_decrease"

	default:
		return "", 0, fmt.Errorf("invalid action: %s", action)
	}

	// Record stake transaction
	stakeRecord := models.VerifierStake{
		VerifierID:    verifierID,
		TransactionID: transactionID,
		Amount:        fmt.Sprintf("%.8f", amount),
		Type:          transactionType,
		Status:        "completed",
		HederaTxID:    transactionID,
	}

	if err := vs.db.Create(&stakeRecord).Error; err != nil {
		return "", 0, fmt.Errorf("failed to record stake transaction: %w", err)
	}

	// Update stake on smart contract
	err = vs.updateVerifierStakeOnChain(verifierID, accountID, newStakeAmount)
	if err != nil {
		return "", 0, fmt.Errorf("failed to update stake on smart contract: %w", err)
	}

	return transactionID, newStakeAmount, nil
}

// GetStakeInfo retrieves stake information from the blockchain
func (vs *VerifierService) GetStakeInfo(accountID string) (*StakeInfo, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("stake_info:%s", accountID)
	cachedInfo, err := vs.redisClient.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var stakeInfo StakeInfo
		if json.Unmarshal([]byte(cachedInfo), &stakeInfo) == nil {
			return &stakeInfo, nil
		}
	}

	// Query smart contract for stake information
	stakeData, err := vs.hederaClient.CallContract(vs.getMarketplaceContractID(), "getVerifierStake", []interface{}{accountID})
	if err != nil {
		return nil, fmt.Errorf("failed to query stake from smart contract: %w", err)
	}

	// Parse stake data (assuming it returns total, available, locked, pending)
	stakeInfo := &StakeInfo{
		Total:     vs.formatHBARAmount(stakeData[0].(int64)),
		Available: vs.formatHBARAmount(stakeData[1].(int64)),
		Locked:    vs.formatHBARAmount(stakeData[2].(int64)),
		Pending:   vs.formatHBARAmount(stakeData[3].(int64)),
	}

	// Cache the result for 5 minutes
	stakeInfoJSON, _ := json.Marshal(stakeInfo)
	vs.redisClient.Set(context.Background(), cacheKey, stakeInfoJSON, 5*time.Minute)

	return stakeInfo, nil
}

// GetReputationDetails retrieves detailed reputation information
func (vs *VerifierService) GetReputationDetails(verifierID string) (*ReputationDetails, error) {
	// Get verifier
	var verifier models.Verifier
	if err := vs.db.First(&verifier, "id = ?", verifierID).Error; err != nil {
		return nil, fmt.Errorf("verifier not found: %w", err)
	}

	// Calculate reputation breakdown
	breakdown, err := vs.calculateReputationBreakdown(verifierID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate reputation breakdown: %w", err)
	}

	// Get recent activities
	activities, err := vs.getRecentReputationActivities(verifierID, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent activities: %w", err)
	}

	// Calculate trends
	trends, err := vs.calculateReputationTrends(verifierID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate reputation trends: %w", err)
	}

	// Get comparisons
	comparisons, err := vs.getReputationComparisons(verifierID, verifier.Specializations)
	if err != nil {
		return nil, fmt.Errorf("failed to get reputation comparisons: %w", err)
	}

	return &ReputationDetails{
		VerifierID:      verifierID,
		OverallScore:    verifier.ReputationScore,
		ScoreBreakdown:  *breakdown,
		RecentActivity:  activities,
		Trends:          *trends,
		Comparisons:     *comparisons,
		LastUpdated:     time.Now(),
	}, nil
}

// GetVerificationStats retrieves verification statistics for a verifier
func (vs *VerifierService) GetVerificationStats(verifierID string) (*VerificationStats, error) {
	now := time.Now()
	last30Days := now.AddDate(0, 0, -30)
	last90Days := now.AddDate(0, 0, -90)

	// Get stats for different timeframes
	stats30, err := vs.getStatsForTimeframe(verifierID, last30Days, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get 30-day stats: %w", err)
	}

	stats90, err := vs.getStatsForTimeframe(verifierID, last90Days, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get 90-day stats: %w", err)
	}

	statsAllTime, err := vs.getStatsForTimeframe(verifierID, time.Time{}, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get all-time stats: %w", err)
	}

	// Calculate average rating
	averageRating, err := vs.getAverageRating(verifierID)
	if err != nil {
		return nil, fmt.Errorf("failed to get average rating: %w", err)
	}

	// Calculate average response time
	responseTime, err := vs.getAverageResponseTime(verifierID)
	if err != nil {
		return nil, fmt.Errorf("failed to get average response time: %w", err)
	}

	return &VerificationStats{
		Last30Days:    *stats30,
		Last90Days:    *stats90,
		AllTime:       *statsAllTime,
		AverageRating: averageRating,
		ResponseTime:  responseTime,
	}, nil
}

// ListVerifiers retrieves a paginated list of verifiers with filtering
func (vs *VerifierService) ListVerifiers(page, limit int, status, specialization, country string, minReputation float64) ([]VerifierListItem, int64, error) {
	offset := (page - 1) * limit

	// Build query
	query := vs.db.Model(&models.Verifier{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if specialization != "" {
		query = query.Where("JSON_CONTAINS(specializations, ?)", fmt.Sprintf(`"%s"`, specialization))
	}

	if country != "" {
		query = query.Where("JSON_CONTAINS(supported_countries, ?)", fmt.Sprintf(`"%s"`, country))
	}

	if minReputation > 0 {
		query = query.Where("reputation_score >= ?", minReputation)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count verifiers: %w", err)
	}

	// Get verifiers
	var verifiers []models.Verifier
	if err := query.Offset(offset).Limit(limit).Order("reputation_score DESC, created_at DESC").Find(&verifiers).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve verifiers: %w", err)
	}

	// Convert to list items
	result := make([]VerifierListItem, len(verifiers))
	for i, verifier := range verifiers {
		successRate := 0.0
		if verifier.TotalVerifications > 0 {
			successRate = float64(verifier.SuccessfulVerifications) / float64(verifier.TotalVerifications) * 100
		}

		// Check if verifier is online (last active within 15 minutes)
		isOnline := false
		if verifier.LastActiveAt != nil {
			isOnline = time.Since(*verifier.LastActiveAt) < 15*time.Minute
		}

		result[i] = VerifierListItem{
			ID:                  verifier.ID,
			Name:                verifier.Name,
			CompanyName:         verifier.CompanyName,
			Specializations:     []string(verifier.Specializations),
			Status:              verifier.Status,
			ReputationScore:     verifier.ReputationScore,
			TotalVerifications:  verifier.TotalVerifications,
			SuccessRate:         successRate,
			AverageResponseTime: vs.formatResponseTime(verifier.AverageResponseTime),
			SupportedCountries:  []string(verifier.SupportedCountries),
			JoinedAt:            verifier.CreatedAt,
			LastActiveAt:        verifier.LastActiveAt,
			IsOnline:            isOnline,
		}
	}

	return result, total, nil
}

// UpdateProfile updates a verifier's profile information
func (vs *VerifierService) UpdateProfile(verifierID string, updates map[string]interface{}) error {
	// Get current verifier
	var verifier models.Verifier
	if err := vs.db.First(&verifier, "id = ?", verifierID).Error; err != nil {
		return fmt.Errorf("verifier not found: %w", err)
	}

	// Validate and apply updates
	allowedFields := map[string]bool{
		"description":         true,
		"website":             true,
		"operating_hours":     true,
		"pricing_tiers":       true,
		"supported_countries": true,
	}

	updateData := make(map[string]interface{})
	for field, value := range updates {
		if allowedFields[field] {
			updateData[field] = value
		}
	}

	if len(updateData) == 0 {
		return fmt.Errorf("no valid fields to update")
	}

	// Update verifier
	if err := vs.db.Model(&verifier).Updates(updateData).Error; err != nil {
		return fmt.Errorf("failed to update verifier profile: %w", err)
	}

	// Clear cache
	vs.clearVerifierCache(verifierID)

	return nil
}

// SendVerificationEmail sends verification email to new verifier
func (vs *VerifierService) SendVerificationEmail(email, name, verifierID string) error {
	if vs.emailService == nil {
		return fmt.Errorf("email service not configured")
	}

	verificationLink := fmt.Sprintf("https://verza.app/verify-verifier/%s", verifierID)

	emailData := map[string]interface{}{
		"name":              name,
		"verification_link": verificationLink,
		"verifier_id":       verifierID,
	}

	return vs.emailService.SendTemplateEmail(email, "verifier_verification", "Welcome to Verza - Verify Your Account", emailData)
}

// Helper methods

func (vs *VerifierService) getMarketplaceAccountID() string {
	// Return the marketplace contract account ID
	return "0.0.123456" // This should be configured
}

func (vs *VerifierService) getMarketplaceContractID() string {
	// Return the marketplace contract ID
	return "0.0.123457" // This should be configured
}

func (vs *VerifierService) updateVerifierStakeOnChain(verifierID, accountID string, amount float64) error {
	// Call smart contract to update verifier stake
	tinybars := int64(amount * 100000000)
	_, err := vs.hederaClient.CallContract(vs.getMarketplaceContractID(), "updateVerifierStake", []interface{}{accountID, tinybars})
	return err
}

func (vs *VerifierService) formatHBARAmount(tinybars int64) string {
	hbar := float64(tinybars) / 100000000
	return fmt.Sprintf("%.8f", hbar)
}

func (vs *VerifierService) calculateReputationBreakdown(verifierID string) (*ReputationBreakdown, error) {
	// This would typically involve complex calculations based on verification history
	// For now, return mock data
	return &ReputationBreakdown{
		Accuracy:        85.5,
		Speed:           92.3,
		Professionalism: 88.7,
		Reliability:     90.1,
		CustomerService: 87.9,
	}, nil
}

func (vs *VerifierService) getRecentReputationActivities(verifierID string, limit int) ([]ReputationActivity, error) {
	// Query recent activities that affected reputation
	activities := []ReputationActivity{
		{
			Date:        time.Now().AddDate(0, 0, -1),
			Type:        "verification_completed",
			Description: "Successfully completed identity verification",
			Impact:      +0.5,
			RequestID:   "req_123",
		},
		{
			Date:        time.Now().AddDate(0, 0, -3),
			Type:        "positive_feedback",
			Description: "Received 5-star rating from client",
			Impact:      +0.3,
			RequestID:   "req_122",
		},
	}
	return activities, nil
}

func (vs *VerifierService) calculateReputationTrends(verifierID string) (*ReputationTrends, error) {
	// Calculate reputation trends over different time periods
	return &ReputationTrends{
		Last7Days:  +0.2,
		Last30Days: +1.5,
		Last90Days: +3.2,
		YearToDate: +8.7,
	}, nil
}

func (vs *VerifierService) getReputationComparisons(verifierID string, specializations models.StringArrayJSON) (*ReputationComparisons, error) {
	// Compare verifier reputation to industry averages
	return &ReputationComparisons{
		IndustryAverage:   82.5,
		SpecializationAvg: 85.3,
		TopPercentile:     95.0,
		RankInCategory:    15,
		TotalInCategory:   150,
	}, nil
}

func (vs *VerifierService) getStatsForTimeframe(verifierID string, startTime, endTime time.Time) (*StatsTimeframe, error) {
	query := vs.db.Model(&models.VerificationRequest{}).Where("verifier_id = ?", verifierID)

	if !startTime.IsZero() {
		query = query.Where("created_at >= ?", startTime)
	}
	query = query.Where("created_at <= ?", endTime)

	var totalRequests int64
	if err := query.Count(&totalRequests).Error; err != nil {
		return nil, err
	}

	var completedRequests int64
	if err := query.Where("status IN ?", []string{"completed", "approved"}).Count(&completedRequests).Error; err != nil {
		return nil, err
	}

	successRate := 0.0
	if totalRequests > 0 {
		successRate = float64(completedRequests) / float64(totalRequests) * 100
	}

	// Calculate earnings (mock for now)
	earnings := "0.00"

	return &StatsTimeframe{
		TotalRequests:     int(totalRequests),
		CompletedRequests: int(completedRequests),
		SuccessRate:       math.Round(successRate*100) / 100,
		Earnings:          earnings,
	}, nil
}

func (vs *VerifierService) getAverageRating(verifierID string) (float64, error) {
	var avgRating float64
	err := vs.db.Model(&models.VerificationRequest{}).
		Where("verifier_id = ? AND rating IS NOT NULL", verifierID).
		Select("AVG(rating)").Scan(&avgRating).Error
	return math.Round(avgRating*100) / 100, err
}

func (vs *VerifierService) getAverageResponseTime(verifierID string) (string, error) {
	// Calculate average response time in hours
	var avgMinutes float64
	err := vs.db.Model(&models.VerificationRequest{}).
		Where("verifier_id = ? AND accepted_at IS NOT NULL", verifierID).
		Select("AVG(TIMESTAMPDIFF(MINUTE, created_at, accepted_at))").Scan(&avgMinutes).Error

	if err != nil {
		return "N/A", err
	}

	if avgMinutes < 60 {
		return fmt.Sprintf("%.0f minutes", avgMinutes), nil
	}
	return fmt.Sprintf("%.1f hours", avgMinutes/60), nil
}

func (vs *VerifierService) formatResponseTime(minutes *int) string {
	if minutes == nil {
		return "N/A"
	}
	if *minutes < 60 {
		return fmt.Sprintf("%d minutes", *minutes)
	}
	return fmt.Sprintf("%.1f hours", float64(*minutes)/60)
}

func (vs *VerifierService) clearVerifierCache(verifierID string) {
	cacheKeys := []string{
		fmt.Sprintf("verifier_profile:%s", verifierID),
		fmt.Sprintf("verifier_stats:%s", verifierID),
		fmt.Sprintf("verifier_reputation:%s", verifierID),
	}

	for _, key := range cacheKeys {
		vs.redisClient.Del(context.Background(), key)
	}
}