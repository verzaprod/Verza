package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// FraudDetectionService handles AI-based fraud detection
type FraudDetectionService struct {
	db          *gorm.DB
	redisClient *redis.Client
	models      *MLModels
}

// MLModels contains the machine learning models
type MLModels struct {
	DocumentClassifier *DocumentClassifier
	AnomalyDetector   *AnomalyDetector
	BehavioralAnalyzer *BehavioralAnalyzer
}

// FraudDetectionRequest represents input data for fraud analysis
type FraudDetectionRequest struct {
	RequestID        string                 `json:"request_id"`
	UserID          string                 `json:"user_id"`
	DocumentData    DocumentMetadata       `json:"document_data"`
	BiometricData   BiometricVerification  `json:"biometric_data"`
	BehavioralData  BehavioralSignals      `json:"behavioral_data"`
	Timestamp       time.Time              `json:"timestamp"`
}

// DocumentMetadata contains document-related information
type DocumentMetadata struct {
	DocumentType    string    `json:"document_type"`
	DocumentHash    string    `json:"document_hash"`
	IssuanceDate    time.Time `json:"issuance_date"`
	ExpirationDate  time.Time `json:"expiration_date"`
	CountryOfOrigin string    `json:"country_of_origin"`
	DocumentNumber  string    `json:"document_number"`
	ImageQuality    float64   `json:"image_quality"`
	OCRConfidence   float64   `json:"ocr_confidence"`
}

// BiometricVerification contains biometric analysis results
type BiometricVerification struct {
	FaceMatchScore   float64 `json:"face_match_score"`
	LivenessScore    float64 `json:"liveness_score"`
	FaceQuality      float64 `json:"face_quality"`
	BiometricHash    string  `json:"biometric_hash"`
	VerificationTime int64   `json:"verification_time_ms"`
}

// BehavioralSignals contains user behavioral data
type BehavioralSignals struct {
	IPAddress        string            `json:"ip_address"`
	UserAgent        string            `json:"user_agent"`
	DeviceFingerprint string           `json:"device_fingerprint"`
	GeoLocation      GeoLocation       `json:"geo_location"`
	RequestFrequency RequestFrequency  `json:"request_frequency"`
	SessionData      SessionData       `json:"session_data"`
}

// GeoLocation represents geographical information
type GeoLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Country   string  `json:"country"`
	City      string  `json:"city"`
	ISP       string  `json:"isp"`
	VPNRisk   float64 `json:"vpn_risk"`
}

// RequestFrequency tracks user request patterns
type RequestFrequency struct {
	RequestsLast24h int     `json:"requests_last_24h"`
	RequestsLastWeek int    `json:"requests_last_week"`
	AverageInterval float64 `json:"average_interval_minutes"`
	BurstPattern    bool    `json:"burst_pattern"`
}

// SessionData contains session-related information
type SessionData struct {
	SessionDuration int64  `json:"session_duration_ms"`
	PageViews       int    `json:"page_views"`
	MouseMovements  int    `json:"mouse_movements"`
	Keystrokes      int    `json:"keystrokes"`
	BotScore        float64 `json:"bot_score"`
}

// FraudDetectionResult represents the analysis output
type FraudDetectionResult struct {
	RequestID       string            `json:"request_id"`
	OverallRiskScore float64          `json:"overall_risk_score"`
	RiskLevel       string            `json:"risk_level"`
	FraudFlags      []FraudFlag       `json:"fraud_flags"`
	ComponentScores ComponentScores   `json:"component_scores"`
	Recommendation  string            `json:"recommendation"`
	Confidence      float64           `json:"confidence"`
	ProcessingTime  int64             `json:"processing_time_ms"`
	Timestamp       time.Time         `json:"timestamp"`
}

// FraudFlag represents a specific fraud indicator
type FraudFlag struct {
	Type        string  `json:"type"`
	Severity    string  `json:"severity"`
	Description string  `json:"description"`
	Score       float64 `json:"score"`
	Evidence    string  `json:"evidence"`
}

// ComponentScores breaks down risk by component
type ComponentScores struct {
	DocumentRisk   float64 `json:"document_risk"`
	BiometricRisk  float64 `json:"biometric_risk"`
	BehavioralRisk float64 `json:"behavioral_risk"`
	AnomalyScore   float64 `json:"anomaly_score"`
}

// NewFraudDetectionService creates a new fraud detection service
func NewFraudDetectionService(db *gorm.DB, redisClient *redis.Client) *FraudDetectionService {
	return &FraudDetectionService{
		db:          db,
		redisClient: redisClient,
		models: &MLModels{
			DocumentClassifier: NewDocumentClassifier(),
			AnomalyDetector:   NewAnomalyDetector(),
			BehavioralAnalyzer: NewBehavioralAnalyzer(),
		},
	}
}

// AnalyzeFraud performs comprehensive fraud analysis
func (fds *FraudDetectionService) AnalyzeFraud(ctx context.Context, request *FraudDetectionRequest) (*FraudDetectionResult, error) {
	startTime := time.Now()

	// Initialize result
	result := &FraudDetectionResult{
		RequestID:  request.RequestID,
		FraudFlags: make([]FraudFlag, 0),
		Timestamp:  time.Now(),
	}

	// Analyze document authenticity
	documentScore, documentFlags, err := fds.analyzeDocument(ctx, &request.DocumentData)
	if err != nil {
		return nil, fmt.Errorf("document analysis failed: %w", err)
	}

	// Analyze biometric data
	biometricScore, biometricFlags, err := fds.analyzeBiometric(ctx, &request.BiometricData)
	if err != nil {
		return nil, fmt.Errorf("biometric analysis failed: %w", err)
	}

	// Analyze behavioral patterns
	behavioralScore, behavioralFlags, err := fds.analyzeBehavioral(ctx, request.UserID, &request.BehavioralData)
	if err != nil {
		return nil, fmt.Errorf("behavioral analysis failed: %w", err)
	}

	// Detect anomalies
	anomalyScore, anomalyFlags, err := fds.detectAnomalies(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("anomaly detection failed: %w", err)
	}

	// Combine scores
	result.ComponentScores = ComponentScores{
		DocumentRisk:   documentScore,
		BiometricRisk:  biometricScore,
		BehavioralRisk: behavioralScore,
		AnomalyScore:   anomalyScore,
	}

	// Calculate overall risk score (weighted average)
	result.OverallRiskScore = fds.calculateOverallRisk(result.ComponentScores)

	// Combine all fraud flags
	result.FraudFlags = append(result.FraudFlags, documentFlags...)
	result.FraudFlags = append(result.FraudFlags, biometricFlags...)
	result.FraudFlags = append(result.FraudFlags, behavioralFlags...)
	result.FraudFlags = append(result.FraudFlags, anomalyFlags...)

	// Determine risk level and recommendation
	result.RiskLevel = fds.determineRiskLevel(result.OverallRiskScore)
	result.Recommendation = fds.generateRecommendation(result)
	result.Confidence = fds.calculateConfidence(result)

	// Calculate processing time
	result.ProcessingTime = time.Since(startTime).Milliseconds()

	// Store result in cache
	if err := fds.cacheResult(ctx, result); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to cache fraud detection result: %v\n", err)
	}

	return result, nil
}

// analyzeDocument performs document authenticity analysis
func (fds *FraudDetectionService) analyzeDocument(ctx context.Context, doc *DocumentMetadata) (float64, []FraudFlag, error) {
	flags := make([]FraudFlag, 0)
	riskScore := 0.0

	// Check document age
	if time.Since(doc.IssuanceDate).Hours() < 24 {
		flags = append(flags, FraudFlag{
			Type:        "DOCUMENT_TOO_NEW",
			Severity:    "MEDIUM",
			Description: "Document issued very recently",
			Score:       30.0,
			Evidence:    fmt.Sprintf("Issued %v hours ago", time.Since(doc.IssuanceDate).Hours()),
		})
		riskScore += 30.0
	}

	// Check document expiration
	if doc.ExpirationDate.Before(time.Now()) {
		flags = append(flags, FraudFlag{
			Type:        "EXPIRED_DOCUMENT",
			Severity:    "HIGH",
			Description: "Document has expired",
			Score:       80.0,
			Evidence:    fmt.Sprintf("Expired on %v", doc.ExpirationDate.Format("2006-01-02")),
		})
		riskScore += 80.0
	}

	// Check image quality
	if doc.ImageQuality < 0.7 {
		flags = append(flags, FraudFlag{
			Type:        "LOW_IMAGE_QUALITY",
			Severity:    "MEDIUM",
			Description: "Document image quality is poor",
			Score:       25.0,
			Evidence:    fmt.Sprintf("Quality score: %.2f", doc.ImageQuality),
		})
		riskScore += 25.0
	}

	// Check OCR confidence
	if doc.OCRConfidence < 0.8 {
		flags = append(flags, FraudFlag{
			Type:        "LOW_OCR_CONFIDENCE",
			Severity:    "MEDIUM",
			Description: "OCR confidence is low, possible tampering",
			Score:       35.0,
			Evidence:    fmt.Sprintf("OCR confidence: %.2f", doc.OCRConfidence),
		})
		riskScore += 35.0
	}

	// Use ML classifier for document authenticity
	classifierScore, err := fds.models.DocumentClassifier.Classify(doc)
	if err != nil {
		return 0, nil, err
	}

	if classifierScore > 0.7 {
		flags = append(flags, FraudFlag{
			Type:        "SYNTHETIC_DOCUMENT",
			Severity:    "HIGH",
			Description: "Document appears to be synthetic or forged",
			Score:       classifierScore * 100,
			Evidence:    fmt.Sprintf("ML classifier score: %.2f", classifierScore),
		})
		riskScore += classifierScore * 100
	}

	return math.Min(riskScore, 100.0), flags, nil
}

// analyzeBiometric performs biometric verification analysis
func (fds *FraudDetectionService) analyzeBiometric(ctx context.Context, bio *BiometricVerification) (float64, []FraudFlag, error) {
	flags := make([]FraudFlag, 0)
	riskScore := 0.0

	// Check face match score
	if bio.FaceMatchScore < 0.8 {
		flags = append(flags, FraudFlag{
			Type:        "LOW_FACE_MATCH",
			Severity:    "HIGH",
			Description: "Face match score is below threshold",
			Score:       (1.0 - bio.FaceMatchScore) * 100,
			Evidence:    fmt.Sprintf("Face match score: %.2f", bio.FaceMatchScore),
		})
		riskScore += (1.0 - bio.FaceMatchScore) * 100
	}

	// Check liveness score
	if bio.LivenessScore < 0.7 {
		flags = append(flags, FraudFlag{
			Type:        "LIVENESS_FAILURE",
			Severity:    "HIGH",
			Description: "Liveness detection failed, possible spoofing",
			Score:       (1.0 - bio.LivenessScore) * 80,
			Evidence:    fmt.Sprintf("Liveness score: %.2f", bio.LivenessScore),
		})
		riskScore += (1.0 - bio.LivenessScore) * 80
	}

	// Check face quality
	if bio.FaceQuality < 0.6 {
		flags = append(flags, FraudFlag{
			Type:        "LOW_FACE_QUALITY",
			Severity:    "MEDIUM",
			Description: "Face image quality is poor",
			Score:       25.0,
			Evidence:    fmt.Sprintf("Face quality: %.2f", bio.FaceQuality),
		})
		riskScore += 25.0
	}

	// Check for duplicate biometric hash
	duplicateCount, err := fds.checkBiometricDuplicates(ctx, bio.BiometricHash)
	if err != nil {
		return 0, nil, err
	}

	if duplicateCount > 0 {
		flags = append(flags, FraudFlag{
			Type:        "DUPLICATE_BIOMETRIC",
			Severity:    "HIGH",
			Description: "Biometric data has been used before",
			Score:       90.0,
			Evidence:    fmt.Sprintf("Found %d previous uses", duplicateCount),
		})
		riskScore += 90.0
	}

	return math.Min(riskScore, 100.0), flags, nil
}

// analyzeBehavioral performs behavioral pattern analysis
func (fds *FraudDetectionService) analyzeBehavioral(ctx context.Context, userID string, behavior *BehavioralSignals) (float64, []FraudFlag, error) {
	flags := make([]FraudFlag, 0)
	riskScore := 0.0

	// Check VPN/Proxy usage
	if behavior.GeoLocation.VPNRisk > 0.7 {
		flags = append(flags, FraudFlag{
			Type:        "VPN_PROXY_USAGE",
			Severity:    "MEDIUM",
			Description: "High probability of VPN/Proxy usage",
			Score:       behavior.GeoLocation.VPNRisk * 40,
			Evidence:    fmt.Sprintf("VPN risk score: %.2f", behavior.GeoLocation.VPNRisk),
		})
		riskScore += behavior.GeoLocation.VPNRisk * 40
	}

	// Check request frequency patterns
	if behavior.RequestFrequency.RequestsLast24h > 10 {
		flags = append(flags, FraudFlag{
			Type:        "HIGH_FREQUENCY_REQUESTS",
			Severity:    "HIGH",
			Description: "Unusually high request frequency",
			Score:       60.0,
			Evidence:    fmt.Sprintf("%d requests in last 24h", behavior.RequestFrequency.RequestsLast24h),
		})
		riskScore += 60.0
	}

	// Check bot score
	if behavior.SessionData.BotScore > 0.8 {
		flags = append(flags, FraudFlag{
			Type:        "BOT_BEHAVIOR",
			Severity:    "HIGH",
			Description: "Behavior indicates automated bot activity",
			Score:       behavior.SessionData.BotScore * 85,
			Evidence:    fmt.Sprintf("Bot score: %.2f", behavior.SessionData.BotScore),
		})
		riskScore += behavior.SessionData.BotScore * 85
	}

	// Use behavioral analyzer
	behaviorScore, err := fds.models.BehavioralAnalyzer.Analyze(userID, behavior)
	if err != nil {
		return 0, nil, err
	}

	if behaviorScore > 0.6 {
		flags = append(flags, FraudFlag{
			Type:        "ANOMALOUS_BEHAVIOR",
			Severity:    "MEDIUM",
			Description: "Behavioral patterns are anomalous",
			Score:       behaviorScore * 50,
			Evidence:    fmt.Sprintf("Behavioral anomaly score: %.2f", behaviorScore),
		})
		riskScore += behaviorScore * 50
	}

	return math.Min(riskScore, 100.0), flags, nil
}

// detectAnomalies performs anomaly detection across all data
func (fds *FraudDetectionService) detectAnomalies(ctx context.Context, request *FraudDetectionRequest) (float64, []FraudFlag, error) {
	flags := make([]FraudFlag, 0)

	// Use anomaly detector
	anomalyScore, anomalies, err := fds.models.AnomalyDetector.Detect(request)
	if err != nil {
		return 0, nil, err
	}

	// Convert anomalies to fraud flags
	for _, anomaly := range anomalies {
		flags = append(flags, FraudFlag{
			Type:        "ANOMALY_" + strings.ToUpper(anomaly.Type),
			Severity:    anomaly.Severity,
			Description: anomaly.Description,
			Score:       anomaly.Score,
			Evidence:    anomaly.Evidence,
		})
	}

	return anomalyScore, flags, nil
}

// calculateOverallRisk combines component scores with weights
func (fds *FraudDetectionService) calculateOverallRisk(scores ComponentScores) float64 {
	// Weighted combination of risk scores
	weights := map[string]float64{
		"document":   0.3,
		"biometric":  0.35,
		"behavioral": 0.25,
		"anomaly":    0.1,
	}

	overallScore := scores.DocumentRisk*weights["document"] +
		scores.BiometricRisk*weights["biometric"] +
		scores.BehavioralRisk*weights["behavioral"] +
		scores.AnomalyScore*weights["anomaly"]

	return math.Min(overallScore, 100.0)
}

// determineRiskLevel categorizes risk score
func (fds *FraudDetectionService) determineRiskLevel(score float64) string {
	switch {
	case score >= 80:
		return "CRITICAL"
	case score >= 60:
		return "HIGH"
	case score >= 40:
		return "MEDIUM"
	case score >= 20:
		return "LOW"
	default:
		return "MINIMAL"
	}
}

// generateRecommendation provides action recommendation
func (fds *FraudDetectionService) generateRecommendation(result *FraudDetectionResult) string {
	switch result.RiskLevel {
	case "CRITICAL":
		return "REJECT - High fraud risk detected. Manual review required."
	case "HIGH":
		return "MANUAL_REVIEW - Significant risk factors present. Human verification needed."
	case "MEDIUM":
		return "ADDITIONAL_VERIFICATION - Moderate risk. Request additional documentation."
	case "LOW":
		return "PROCEED_WITH_CAUTION - Low risk but monitor for patterns."
	default:
		return "APPROVE - Minimal fraud risk detected."
	}
}

// calculateConfidence determines confidence in the analysis
func (fds *FraudDetectionService) calculateConfidence(result *FraudDetectionResult) float64 {
	// Base confidence on number of data points and consistency of flags
	flagCount := len(result.FraudFlags)
	highSeverityCount := 0

	for _, flag := range result.FraudFlags {
		if flag.Severity == "HIGH" || flag.Severity == "CRITICAL" {
			highSeverityCount++
		}
	}

	// Higher confidence with more consistent high-severity flags
	confidence := 0.6 + (float64(highSeverityCount)/float64(math.Max(1, float64(flagCount))))*0.4
	return math.Min(confidence, 0.95)
}

// checkBiometricDuplicates checks for duplicate biometric hashes
func (fds *FraudDetectionService) checkBiometricDuplicates(ctx context.Context, biometricHash string) (int, error) {
	// Check Redis cache first
	cacheKey := fmt.Sprintf("biometric_hash:%s", biometricHash)
	count, err := fds.redisClient.Get(ctx, cacheKey).Int()
	if err == nil {
		return count, nil
	}

	// Query database if not in cache
	var duplicateCount int64
	err = fds.db.Table("fraud_detection_results").
		Where("JSON_EXTRACT(biometric_data, '$.biometric_hash') = ?", biometricHash).
		Count(&duplicateCount).Error

	if err != nil {
		return 0, err
	}

	// Cache the result
	fds.redisClient.Set(ctx, cacheKey, int(duplicateCount), time.Hour*24)

	return int(duplicateCount), nil
}

// cacheResult stores the fraud detection result in cache
func (fds *FraudDetectionService) cacheResult(ctx context.Context, result *FraudDetectionResult) error {
	cacheKey := fmt.Sprintf("fraud_result:%s", result.RequestID)
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return fds.redisClient.Set(ctx, cacheKey, resultJSON, time.Hour*24).Err()
}

// GetCachedResult retrieves cached fraud detection result
func (fds *FraudDetectionService) GetCachedResult(ctx context.Context, requestID string) (*FraudDetectionResult, error) {
	cacheKey := fmt.Sprintf("fraud_result:%s", requestID)
	resultJSON, err := fds.redisClient.Get(ctx, cacheKey).Result()
	if err != nil {
		return nil, err
	}

	var result FraudDetectionResult
	err = json.Unmarshal([]byte(resultJSON), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}