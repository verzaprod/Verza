package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/verza/services"
)

// FraudHandler handles fraud detection HTTP requests
type FraudHandler struct {
	fraudService *services.FraudDetectionService
}

// FraudScoreRequest represents the request for fraud scoring
type FraudScoreRequest struct {
	VerificationRequestID string                 `json:"verification_request_id" binding:"required"`
	DocumentData          *DocumentAnalysisData  `json:"document_data,omitempty"`
	BiometricData         *BiometricAnalysisData `json:"biometric_data,omitempty"`
	BehavioralData        *BehavioralAnalysisData `json:"behavioral_data,omitempty"`
	SkipCache             bool                   `json:"skip_cache,omitempty"`
}

// DocumentAnalysisData represents document analysis input
type DocumentAnalysisData struct {
	DocumentType     string            `json:"document_type"`
	DocumentHash     string            `json:"document_hash"`
	IssuanceDate     string            `json:"issuance_date"`
	ExpirationDate   string            `json:"expiration_date"`
	IssuerCountry    string            `json:"issuer_country"`
	DocumentNumber   string            `json:"document_number"`
	OCRConfidence    float64           `json:"ocr_confidence"`
	ImageQuality     float64           `json:"image_quality"`
	SecurityFeatures map[string]bool   `json:"security_features"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// BiometricAnalysisData represents biometric analysis input
type BiometricAnalysisData struct {
	FaceMatchScore   float64                `json:"face_match_score"`
	LivenessScore    float64                `json:"liveness_score"`
	FaceHash         string                 `json:"face_hash"`
	BiometricQuality float64                `json:"biometric_quality"`
	DeviceInfo       map[string]interface{} `json:"device_info,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// BehavioralAnalysisData represents behavioral analysis input
type BehavioralAnalysisData struct {
	UserID           string                 `json:"user_id"`
	IPAddress        string                 `json:"ip_address"`
	UserAgent        string                 `json:"user_agent"`
	DeviceFingerprint string                `json:"device_fingerprint"`
	SessionDuration  int                    `json:"session_duration"`
	RequestFrequency int                    `json:"request_frequency"`
	Geolocation      map[string]interface{} `json:"geolocation,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// FraudScoreResponse represents the fraud scoring response
type FraudScoreResponse struct {
	VerificationRequestID string                   `json:"verification_request_id"`
	RiskScore            float64                  `json:"risk_score"`
	RiskLevel            string                   `json:"risk_level"`
	Flags                []string                 `json:"flags"`
	Recommendation       string                   `json:"recommendation"`
	Confidence           float64                  `json:"confidence"`
	AnalysisBreakdown    FraudAnalysisBreakdown   `json:"analysis_breakdown"`
	ProcessingTime       string                   `json:"processing_time"`
	AnalyzedAt           string                   `json:"analyzed_at"`
	Metadata             map[string]interface{}   `json:"metadata,omitempty"`
}

// FraudAnalysisBreakdown represents detailed analysis breakdown
type FraudAnalysisBreakdown struct {
	DocumentScore   *float64 `json:"document_score,omitempty"`
	BiometricScore  *float64 `json:"biometric_score,omitempty"`
	BehavioralScore *float64 `json:"behavioral_score,omitempty"`
	AnomalyScore    *float64 `json:"anomaly_score,omitempty"`
	PatternScore    *float64 `json:"pattern_score,omitempty"`
}

// FraudResultResponse represents stored fraud analysis result
type FraudResultResponse struct {
	ID                   string                 `json:"id"`
	VerificationRequestID string                 `json:"verification_request_id"`
	RiskScore            float64                `json:"risk_score"`
	RiskLevel            string                 `json:"risk_level"`
	Flags                []string               `json:"flags"`
	Recommendation       string                 `json:"recommendation"`
	Confidence           float64                `json:"confidence"`
	AnalysisBreakdown    FraudAnalysisBreakdown `json:"analysis_breakdown"`
	ProcessingTime       string                 `json:"processing_time"`
	AnalyzedAt           string                 `json:"analyzed_at"`
	CreatedAt            string                 `json:"created_at"`
	UpdatedAt            string                 `json:"updated_at"`
}

// FraudHistoryResponse represents fraud analysis history
type FraudHistoryResponse struct {
	Results    []FraudResultSummary `json:"results"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	Limit      int                  `json:"limit"`
	TotalPages int                  `json:"total_pages"`
}

// FraudResultSummary represents a summary of fraud analysis result
type FraudResultSummary struct {
	ID                   string  `json:"id"`
	VerificationRequestID string  `json:"verification_request_id"`
	RiskScore            float64 `json:"risk_score"`
	RiskLevel            string  `json:"risk_level"`
	FlagCount            int     `json:"flag_count"`
	Recommendation       string  `json:"recommendation"`
	AnalyzedAt           string  `json:"analyzed_at"`
}

// FraudStatsResponse represents fraud detection statistics
type FraudStatsResponse struct {
	TotalAnalyses     int64                  `json:"total_analyses"`
	HighRiskCount     int64                  `json:"high_risk_count"`
	MediumRiskCount   int64                  `json:"medium_risk_count"`
	LowRiskCount      int64                  `json:"low_risk_count"`
	AverageRiskScore  float64                `json:"average_risk_score"`
	CommonFlags       []FlagStatistic        `json:"common_flags"`
	TrendData         []TrendDataPoint       `json:"trend_data"`
	ProcessingMetrics ProcessingMetrics      `json:"processing_metrics"`
}

// FlagStatistic represents statistics for fraud flags
type FlagStatistic struct {
	Flag        string  `json:"flag"`
	Count       int64   `json:"count"`
	Percentage  float64 `json:"percentage"`
}

// TrendDataPoint represents trend data over time
type TrendDataPoint struct {
	Date             string  `json:"date"`
	AnalysisCount    int64   `json:"analysis_count"`
	AverageRiskScore float64 `json:"average_risk_score"`
	HighRiskCount    int64   `json:"high_risk_count"`
}

// ProcessingMetrics represents processing performance metrics
type ProcessingMetrics struct {
	AverageProcessingTime string  `json:"average_processing_time"`
	MinProcessingTime     string  `json:"min_processing_time"`
	MaxProcessingTime     string  `json:"max_processing_time"`
	SuccessRate           float64 `json:"success_rate"`
}

// NewFraudHandler creates a new fraud detection handler
func NewFraudHandler(fraudService *services.FraudDetectionService) *FraudHandler {
	return &FraudHandler{
		fraudService: fraudService,
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

// AnalyzeFraudRisk handles POST /fraud/score
func (fh *FraudHandler) AnalyzeFraudRisk(c *gin.Context) {
	var req FraudScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	// Get user ID from context for authorization
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	userRole, _ := c.Get("user_role")

	// Validate verification request access
	if userRole != "admin" && userRole != "system" {
		// Check if user owns the verification request or is the assigned verifier
		hasAccess, err := fh.fraudService.ValidateVerificationAccess(req.VerificationRequestID, userID.(string))
		if err != nil || !hasAccess {
			utils.ErrorResponse(c, http.StatusForbidden, "Not authorized to analyze this verification request", nil)
			return
		}
	}

	// Convert request data to service format
	analysisData := fh.convertToAnalysisData(&req)

	// Perform fraud analysis
	result, err := fh.fraudService.AnalyzeVerificationRequestWithData(req.VerificationRequestID, analysisData, req.SkipCache)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to analyze fraud risk", err)
		return
	}

	// Convert result to response format
	response := fh.convertToFraudScoreResponse(result)

	utils.SuccessResponse(c, "Fraud analysis completed successfully", response)
}

// GetFraudResult handles GET /fraud/result/:id
func (fh *FraudHandler) GetFraudResult(c *gin.Context) {
	resultID := c.Param("id")
	if resultID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Result ID is required", nil)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	userRole, _ := c.Get("user_role")

	// Get fraud result
	result, err := fh.fraudService.GetFraudResult(resultID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Fraud result not found", err)
		return
	}

	// Check authorization
	if userRole != "admin" && userRole != "system" {
		hasAccess, err := fh.fraudService.ValidateVerificationAccess(result.VerificationRequestID, userID.(string))
		if err != nil || !hasAccess {
			utils.ErrorResponse(c, http.StatusForbidden, "Not authorized to view this fraud result", nil)
			return
		}
	}

	// Convert to response format
	response := fh.convertToFraudResultResponse(result)

	utils.SuccessResponse(c, "Fraud result retrieved successfully", response)
}

// GetFraudHistory handles GET /fraud/history
func (fh *FraudHandler) GetFraudHistory(c *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit > 100 {
		limit = 100
	}

	// Get filter parameters
	riskLevel := c.Query("risk_level")
	verificationRequestID := c.Query("verification_request_id")
	minRiskScore, _ := strconv.ParseFloat(c.Query("min_risk_score"), 64)
	maxRiskScore, _ := strconv.ParseFloat(c.Query("max_risk_score"), 64)

	// Get user ID and role from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	userRole, _ := c.Get("user_role")

	// Apply user-based filtering for non-admin users
	var userFilter *string
	if userRole != "admin" && userRole != "system" {
		userIDStr := userID.(string)
		userFilter = &userIDStr
	}

	// Get fraud history
	results, total, err := fh.fraudService.GetFraudHistory(
		page, limit, riskLevel, verificationRequestID,
		minRiskScore, maxRiskScore, userFilter,
	)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve fraud history", err)
		return
	}

	// Convert to response format
	summaries := make([]FraudResultSummary, len(results))
	for i, result := range results {
		summaries[i] = fh.convertToFraudResultSummary(&result)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	response := FraudHistoryResponse{
		Results:    summaries,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	utils.SuccessResponse(c, "Fraud history retrieved successfully", response)
}

// GetFraudStats handles GET /fraud/stats
func (fh *FraudHandler) GetFraudStats(c *gin.Context) {
	// Get user role from context
	userRole, _ := c.Get("user_role")

	// Only allow admin and system users to view stats
	if userRole != "admin" && userRole != "system" {
		utils.ErrorResponse(c, http.StatusForbidden, "Not authorized to view fraud statistics", nil)
		return
	}

	// Get time range parameters
	daysBack, _ := strconv.Atoi(c.DefaultQuery("days_back", "30"))
	if daysBack > 365 {
		daysBack = 365
	}

	// Get fraud statistics
	stats, err := fh.fraudService.GetFraudStatistics(daysBack)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve fraud statistics", err)
		return
	}

	// Convert to response format
	response := fh.convertToFraudStatsResponse(stats)

	utils.SuccessResponse(c, "Fraud statistics retrieved successfully", response)
}

// ReanalyzeVerification handles POST /fraud/reanalyze/:id
func (fh *FraudHandler) ReanalyzeVerification(c *gin.Context) {
	verificationRequestID := c.Param("id")
	if verificationRequestID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Verification request ID is required", nil)
		return
	}

	// Get user role from context
	userRole, _ := c.Get("user_role")

	// Only allow admin and system users to reanalyze
	if userRole != "admin" && userRole != "system" {
		utils.ErrorResponse(c, http.StatusForbidden, "Not authorized to reanalyze verification requests", nil)
		return
	}

	// Force reanalysis by skipping cache
	result, err := fh.fraudService.AnalyzeVerificationRequest(verificationRequestID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to reanalyze verification request", err)
		return
	}

	// Convert result to response format
	response := fh.convertToFraudScoreResponse(result)

	utils.SuccessResponse(c, "Verification request reanalyzed successfully", response)
}

// Helper methods

func (fh *FraudHandler) convertToAnalysisData(req *FraudScoreRequest) *services.FraudAnalysisData {
	data := &services.FraudAnalysisData{}

	if req.DocumentData != nil {
		data.DocumentData = &services.DocumentAnalysisData{
			DocumentType:     req.DocumentData.DocumentType,
			DocumentHash:     req.DocumentData.DocumentHash,
			IssuanceDate:     req.DocumentData.IssuanceDate,
			ExpirationDate:   req.DocumentData.ExpirationDate,
			IssuerCountry:    req.DocumentData.IssuerCountry,
			DocumentNumber:   req.DocumentData.DocumentNumber,
			OCRConfidence:    req.DocumentData.OCRConfidence,
			ImageQuality:     req.DocumentData.ImageQuality,
			SecurityFeatures: req.DocumentData.SecurityFeatures,
			Metadata:         req.DocumentData.Metadata,
		}
	}

	if req.BiometricData != nil {
		data.BiometricData = &services.BiometricAnalysisData{
			FaceMatchScore:   req.BiometricData.FaceMatchScore,
			LivenessScore:    req.BiometricData.LivenessScore,
			FaceHash:         req.BiometricData.FaceHash,
			BiometricQuality: req.BiometricData.BiometricQuality,
			DeviceInfo:       req.BiometricData.DeviceInfo,
			Metadata:         req.BiometricData.Metadata,
		}
	}

	if req.BehavioralData != nil {
		data.BehavioralData = &services.BehavioralAnalysisData{
			UserID:           req.BehavioralData.UserID,
			IPAddress:        req.BehavioralData.IPAddress,
			UserAgent:        req.BehavioralData.UserAgent,
			DeviceFingerprint: req.BehavioralData.DeviceFingerprint,
			SessionDuration:  req.BehavioralData.SessionDuration,
			RequestFrequency: req.BehavioralData.RequestFrequency,
			Geolocation:      req.BehavioralData.Geolocation,
			Metadata:         req.BehavioralData.Metadata,
		}
	}

	return data
}

func (fh *FraudHandler) convertToFraudScoreResponse(result *services.FraudDetectionResult) *FraudScoreResponse {
	return &FraudScoreResponse{
		VerificationRequestID: result.VerificationRequestID,
		RiskScore:            result.RiskScore,
		RiskLevel:            result.RiskLevel,
		Flags:                result.Flags,
		Recommendation:       result.Recommendation,
		Confidence:           result.Confidence,
		AnalysisBreakdown: FraudAnalysisBreakdown{
			DocumentScore:   result.DocumentScore,
			BiometricScore:  result.BiometricScore,
			BehavioralScore: result.BehavioralScore,
			AnomalyScore:    result.AnomalyScore,
			PatternScore:    result.PatternScore,
		},
		ProcessingTime: result.ProcessingTime,
		AnalyzedAt:     result.AnalyzedAt.Format("2006-01-02T15:04:05Z07:00"),
		Metadata:       result.Metadata,
	}
}

func (fh *FraudHandler) convertToFraudResultResponse(result *models.FraudDetectionResult) *FraudResultResponse {
	return &FraudResultResponse{
		ID:                   result.ID,
		VerificationRequestID: result.VerificationRequestID,
		RiskScore:            result.RiskScore,
		RiskLevel:            result.RiskLevel,
		Flags:                result.Flags,
		Recommendation:       result.Recommendation,
		Confidence:           result.Confidence,
		AnalysisBreakdown: FraudAnalysisBreakdown{
			DocumentScore:   result.DocumentScore,
			BiometricScore:  result.BiometricScore,
			BehavioralScore: result.BehavioralScore,
			AnomalyScore:    result.AnomalyScore,
			PatternScore:    result.PatternScore,
		},
		ProcessingTime: result.ProcessingTime,
		AnalyzedAt:     result.AnalyzedAt.Format("2006-01-02T15:04:05Z07:00"),
		CreatedAt:      result.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      result.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (fh *FraudHandler) convertToFraudResultSummary(result *models.FraudDetectionResult) FraudResultSummary {
	return FraudResultSummary{
		ID:                   result.ID,
		VerificationRequestID: result.VerificationRequestID,
		RiskScore:            result.RiskScore,
		RiskLevel:            result.RiskLevel,
		FlagCount:            len(result.Flags),
		Recommendation:       result.Recommendation,
		AnalyzedAt:           result.AnalyzedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (fh *FraudHandler) convertToFraudStatsResponse(stats *services.FraudStatistics) *FraudStatsResponse {
	// Convert flag statistics
	flagStats := make([]FlagStatistic, len(stats.CommonFlags))
	for i, flag := range stats.CommonFlags {
		flagStats[i] = FlagStatistic{
			Flag:       flag.Flag,
			Count:      flag.Count,
			Percentage: flag.Percentage,
		}
	}

	// Convert trend data
	trendData := make([]TrendDataPoint, len(stats.TrendData))
	for i, trend := range stats.TrendData {
		trendData[i] = TrendDataPoint{
			Date:             trend.Date,
			AnalysisCount:    trend.AnalysisCount,
			AverageRiskScore: trend.AverageRiskScore,
			HighRiskCount:    trend.HighRiskCount,
		}
	}

	return &FraudStatsResponse{
		TotalAnalyses:   stats.TotalAnalyses,
		HighRiskCount:   stats.HighRiskCount,
		MediumRiskCount: stats.MediumRiskCount,
		LowRiskCount:    stats.LowRiskCount,
		AverageRiskScore: stats.AverageRiskScore,
		CommonFlags:     flagStats,
		TrendData:       trendData,
		ProcessingMetrics: ProcessingMetrics{
			AverageProcessingTime: stats.ProcessingMetrics.AverageProcessingTime,
			MinProcessingTime:     stats.ProcessingMetrics.MinProcessingTime,
			MaxProcessingTime:     stats.ProcessingMetrics.MaxProcessingTime,
			SuccessRate:           stats.ProcessingMetrics.SuccessRate,
		},
	}
}