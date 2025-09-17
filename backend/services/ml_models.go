package services

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

// DocumentClassifier analyzes document authenticity
type DocumentClassifier struct {
	model *DocumentMLModel
}

// DocumentMLModel represents the ML model for document classification
type DocumentMLModel struct {
	weights map[string]float64
	thresholds map[string]float64
}

// AnomalyDetector detects anomalous patterns in verification requests
type AnomalyDetector struct {
	model *AnomalyMLModel
}

// AnomalyMLModel represents the anomaly detection model
type AnomalyMLModel struct {
	featureStats map[string]FeatureStats
	threshold    float64
}

// FeatureStats contains statistical information for features
type FeatureStats struct {
	Mean   float64 `json:"mean"`
	StdDev float64 `json:"std_dev"`
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
}

// BehavioralAnalyzer analyzes user behavioral patterns
type BehavioralAnalyzer struct {
	model *BehavioralMLModel
}

// BehavioralMLModel represents the behavioral analysis model
type BehavioralMLModel struct {
	userProfiles map[string]*UserBehaviorProfile
	globalStats  *GlobalBehaviorStats
}

// UserBehaviorProfile contains user-specific behavioral patterns
type UserBehaviorProfile struct {
	UserID           string                 `json:"user_id"`
	RequestHistory   []RequestPattern       `json:"request_history"`
	DeviceHistory    []string               `json:"device_history"`
	LocationHistory  []GeoLocation          `json:"location_history"`
	BehaviorMetrics  BehaviorMetrics        `json:"behavior_metrics"`
	LastUpdated      time.Time              `json:"last_updated"`
}

// RequestPattern represents a user's request pattern
type RequestPattern struct {
	Timestamp       time.Time `json:"timestamp"`
	RequestType     string    `json:"request_type"`
	SessionDuration int64     `json:"session_duration"`
	Success         bool      `json:"success"`
}

// BehaviorMetrics contains behavioral analysis metrics
type BehaviorMetrics struct {
	AverageSessionDuration float64 `json:"average_session_duration"`
	RequestFrequency       float64 `json:"request_frequency"`
	SuccessRate           float64 `json:"success_rate"`
	DeviceConsistency     float64 `json:"device_consistency"`
	LocationConsistency   float64 `json:"location_consistency"`
}

// GlobalBehaviorStats contains global behavioral statistics
type GlobalBehaviorStats struct {
	AverageSessionDuration float64 `json:"average_session_duration"`
	AverageRequestFreq     float64 `json:"average_request_frequency"`
	CommonDevicePatterns   []string `json:"common_device_patterns"`
	CommonLocationPatterns []string `json:"common_location_patterns"`
}

// Anomaly represents a detected anomaly
type Anomaly struct {
	Type        string  `json:"type"`
	Severity    string  `json:"severity"`
	Description string  `json:"description"`
	Score       float64 `json:"score"`
	Evidence    string  `json:"evidence"`
}

// NewDocumentClassifier creates a new document classifier
func NewDocumentClassifier() *DocumentClassifier {
	return &DocumentClassifier{
		model: &DocumentMLModel{
			weights: map[string]float64{
				"image_quality":    0.25,
				"ocr_confidence":   0.30,
				"document_age":     0.20,
				"format_validity":  0.15,
				"metadata_consistency": 0.10,
			},
			thresholds: map[string]float64{
				"synthetic_threshold": 0.7,
				"tampered_threshold":  0.6,
				"quality_threshold":   0.5,
			},
		},
	}
}

// Classify analyzes document authenticity and returns fraud probability
func (dc *DocumentClassifier) Classify(doc *DocumentMetadata) (float64, error) {
	if doc == nil {
		return 0, fmt.Errorf("document metadata is nil")
	}

	// Extract features
	features := dc.extractFeatures(doc)

	// Calculate weighted score
	fraudScore := 0.0
	for feature, value := range features {
		if weight, exists := dc.model.weights[feature]; exists {
			fraudScore += value * weight
		}
	}

	// Apply sigmoid function to normalize score
	normalizedScore := 1.0 / (1.0 + math.Exp(-fraudScore))

	return normalizedScore, nil
}

// extractFeatures extracts relevant features from document metadata
func (dc *DocumentClassifier) extractFeatures(doc *DocumentMetadata) map[string]float64 {
	features := make(map[string]float64)

	// Image quality feature (inverted - lower quality = higher fraud risk)
	features["image_quality"] = 1.0 - doc.ImageQuality

	// OCR confidence feature (inverted)
	features["ocr_confidence"] = 1.0 - doc.OCRConfidence

	// Document age feature (very new or very old documents are suspicious)
	docAge := time.Since(doc.IssuanceDate).Hours() / 24 // age in days
	if docAge < 1 {
		features["document_age"] = 0.8 // Very new
	} else if docAge > 3650 { // > 10 years
		features["document_age"] = 0.6 // Very old
	} else {
		features["document_age"] = 0.1 // Normal age
	}

	// Format validity (check document number format)
	features["format_validity"] = dc.validateDocumentFormat(doc)

	// Metadata consistency
	features["metadata_consistency"] = dc.checkMetadataConsistency(doc)

	return features
}

// validateDocumentFormat checks if document format is valid
func (dc *DocumentClassifier) validateDocumentFormat(doc *DocumentMetadata) float64 {
	// Simple format validation based on document type
	switch strings.ToLower(doc.DocumentType) {
	case "passport":
		// Passport numbers are typically 6-9 characters
		if len(doc.DocumentNumber) < 6 || len(doc.DocumentNumber) > 9 {
			return 0.7 // Invalid format
		}
	case "driver_license", "id_card":
		// Driver license/ID numbers vary but typically 8-15 characters
		if len(doc.DocumentNumber) < 8 || len(doc.DocumentNumber) > 15 {
			return 0.6 // Suspicious format
		}
	default:
		return 0.3 // Unknown document type
	}

	return 0.1 // Valid format
}

// checkMetadataConsistency checks consistency of metadata fields
func (dc *DocumentClassifier) checkMetadataConsistency(doc *DocumentMetadata) float64 {
	inconsistencyScore := 0.0

	// Check if expiration date is after issuance date
	if doc.ExpirationDate.Before(doc.IssuanceDate) {
		inconsistencyScore += 0.8
	}

	// Check if document is expired but image quality is very high (suspicious)
	if doc.ExpirationDate.Before(time.Now()) && doc.ImageQuality > 0.95 {
		inconsistencyScore += 0.5
	}

	// Check country consistency (basic validation)
	if doc.CountryOfOrigin == "" {
		inconsistencyScore += 0.3
	}

	return math.Min(inconsistencyScore, 1.0)
}

// NewAnomalyDetector creates a new anomaly detector
func NewAnomalyDetector() *AnomalyDetector {
	return &AnomalyDetector{
		model: &AnomalyMLModel{
			featureStats: map[string]FeatureStats{
				"face_match_score": {Mean: 0.85, StdDev: 0.12, Min: 0.0, Max: 1.0},
				"liveness_score":   {Mean: 0.82, StdDev: 0.15, Min: 0.0, Max: 1.0},
				"image_quality":    {Mean: 0.78, StdDev: 0.18, Min: 0.0, Max: 1.0},
				"ocr_confidence":   {Mean: 0.88, StdDev: 0.10, Min: 0.0, Max: 1.0},
				"session_duration": {Mean: 180000, StdDev: 60000, Min: 0, Max: 600000}, // milliseconds
				"requests_24h":     {Mean: 1.2, StdDev: 0.8, Min: 0, Max: 50},
			},
			threshold: 2.5, // Z-score threshold for anomaly detection
		},
	}
}

// Detect identifies anomalies in the fraud detection request
func (ad *AnomalyDetector) Detect(request *FraudDetectionRequest) (float64, []Anomaly, error) {
	if request == nil {
		return 0, nil, fmt.Errorf("fraud detection request is nil")
	}

	anomalies := make([]Anomaly, 0)
	totalAnomalyScore := 0.0

	// Extract features for anomaly detection
	features := ad.extractFeatures(request)

	// Check each feature for anomalies
	for featureName, value := range features {
		if stats, exists := ad.model.featureStats[featureName]; exists {
			// Calculate Z-score
			zScore := math.Abs((value - stats.Mean) / stats.StdDev)

			if zScore > ad.model.threshold {
				// Anomaly detected
				anomalyScore := math.Min(zScore/ad.model.threshold*50, 100) // Scale to 0-100
				totalAnomalyScore += anomalyScore

				severity := "MEDIUM"
				if zScore > ad.model.threshold*2 {
					severity = "HIGH"
				}

				anomalies = append(anomalies, Anomaly{
					Type:        featureName,
					Severity:    severity,
					Description: fmt.Sprintf("Anomalous %s detected", strings.ReplaceAll(featureName, "_", " ")),
					Score:       anomalyScore,
					Evidence:    fmt.Sprintf("Z-score: %.2f, Value: %.2f, Expected: %.2f±%.2f", zScore, value, stats.Mean, stats.StdDev),
				})
			}
		}
	}

	// Check for correlation anomalies
	correlationAnomalies := ad.detectCorrelationAnomalies(request)
	anomalies = append(anomalies, correlationAnomalies...)

	// Calculate overall anomaly score
	overallScore := math.Min(totalAnomalyScore/float64(len(ad.model.featureStats))*100, 100)

	return overallScore, anomalies, nil
}

// extractFeatures extracts numerical features for anomaly detection
func (ad *AnomalyDetector) extractFeatures(request *FraudDetectionRequest) map[string]float64 {
	features := make(map[string]float64)

	// Biometric features
	features["face_match_score"] = request.BiometricData.FaceMatchScore
	features["liveness_score"] = request.BiometricData.LivenessScore

	// Document features
	features["image_quality"] = request.DocumentData.ImageQuality
	features["ocr_confidence"] = request.DocumentData.OCRConfidence

	// Behavioral features
	features["session_duration"] = float64(request.BehavioralData.SessionData.SessionDuration)
	features["requests_24h"] = float64(request.BehavioralData.RequestFrequency.RequestsLast24h)

	return features
}

// detectCorrelationAnomalies detects anomalies based on feature correlations
func (ad *AnomalyDetector) detectCorrelationAnomalies(request *FraudDetectionRequest) []Anomaly {
	anomalies := make([]Anomaly, 0)

	// High image quality but low OCR confidence (possible tampering)
	if request.DocumentData.ImageQuality > 0.9 && request.DocumentData.OCRConfidence < 0.6 {
		anomalies = append(anomalies, Anomaly{
			Type:        "quality_ocr_mismatch",
			Severity:    "HIGH",
			Description: "High image quality but low OCR confidence suggests tampering",
			Score:       75.0,
			Evidence:    fmt.Sprintf("Image quality: %.2f, OCR confidence: %.2f", request.DocumentData.ImageQuality, request.DocumentData.OCRConfidence),
		})
	}

	// High face match but low liveness (possible photo attack)
	if request.BiometricData.FaceMatchScore > 0.9 && request.BiometricData.LivenessScore < 0.5 {
		anomalies = append(anomalies, Anomaly{
			Type:        "face_liveness_mismatch",
			Severity:    "HIGH",
			Description: "High face match but low liveness suggests photo attack",
			Score:       80.0,
			Evidence:    fmt.Sprintf("Face match: %.2f, Liveness: %.2f", request.BiometricData.FaceMatchScore, request.BiometricData.LivenessScore),
		})
	}

	// Very short session with high bot score
	if request.BehavioralData.SessionData.SessionDuration < 30000 && request.BehavioralData.SessionData.BotScore > 0.8 {
		anomalies = append(anomalies, Anomaly{
			Type:        "bot_short_session",
			Severity:    "MEDIUM",
			Description: "Very short session with high bot score indicates automation",
			Score:       60.0,
			Evidence:    fmt.Sprintf("Session duration: %dms, Bot score: %.2f", request.BehavioralData.SessionData.SessionDuration, request.BehavioralData.SessionData.BotScore),
		})
	}

	return anomalies
}

// NewBehavioralAnalyzer creates a new behavioral analyzer
func NewBehavioralAnalyzer() *BehavioralAnalyzer {
	return &BehavioralAnalyzer{
		model: &BehavioralMLModel{
			userProfiles: make(map[string]*UserBehaviorProfile),
			globalStats: &GlobalBehaviorStats{
				AverageSessionDuration: 180000, // 3 minutes in milliseconds
				AverageRequestFreq:     1.2,    // requests per day
				CommonDevicePatterns:   []string{"Chrome", "Safari", "Firefox", "Edge"},
				CommonLocationPatterns: []string{"US", "CA", "GB", "AU", "DE"},
			},
		},
	}
}

// Analyze performs behavioral analysis for a user
func (ba *BehavioralAnalyzer) Analyze(userID string, behavior *BehavioralSignals) (float64, error) {
	if behavior == nil {
		return 0, fmt.Errorf("behavioral signals are nil")
	}

	// Get or create user profile
	profile := ba.getUserProfile(userID)

	// Update profile with current behavior
	ba.updateUserProfile(profile, behavior)

	// Calculate behavioral risk score
	riskScore := ba.calculateBehavioralRisk(profile, behavior)

	return riskScore, nil
}

// getUserProfile retrieves or creates a user behavior profile
func (ba *BehavioralAnalyzer) getUserProfile(userID string) *UserBehaviorProfile {
	if profile, exists := ba.model.userProfiles[userID]; exists {
		return profile
	}

	// Create new profile
	profile := &UserBehaviorProfile{
		UserID:          userID,
		RequestHistory:  make([]RequestPattern, 0),
		DeviceHistory:   make([]string, 0),
		LocationHistory: make([]GeoLocation, 0),
		BehaviorMetrics: BehaviorMetrics{},
		LastUpdated:     time.Now(),
	}

	ba.model.userProfiles[userID] = profile
	return profile
}

// updateUserProfile updates the user profile with new behavioral data
func (ba *BehavioralAnalyzer) updateUserProfile(profile *UserBehaviorProfile, behavior *BehavioralSignals) {
	// Add current request to history
	pattern := RequestPattern{
		Timestamp:       time.Now(),
		RequestType:     "verification",
		SessionDuration: behavior.SessionData.SessionDuration,
		Success:         true, // Will be updated based on final result
	}
	profile.RequestHistory = append(profile.RequestHistory, pattern)

	// Limit history size
	if len(profile.RequestHistory) > 100 {
		profile.RequestHistory = profile.RequestHistory[1:]
	}

	// Update device history
	deviceFingerprint := behavior.DeviceFingerprint
	if !ba.contains(profile.DeviceHistory, deviceFingerprint) {
		profile.DeviceHistory = append(profile.DeviceHistory, deviceFingerprint)
		if len(profile.DeviceHistory) > 10 {
			profile.DeviceHistory = profile.DeviceHistory[1:]
		}
	}

	// Update location history
	profile.LocationHistory = append(profile.LocationHistory, behavior.GeoLocation)
	if len(profile.LocationHistory) > 50 {
		profile.LocationHistory = profile.LocationHistory[1:]
	}

	// Recalculate behavior metrics
	ba.calculateBehaviorMetrics(profile)

	profile.LastUpdated = time.Now()
}

// calculateBehaviorMetrics calculates behavioral metrics for a user profile
func (ba *BehavioralAnalyzer) calculateBehaviorMetrics(profile *UserBehaviorProfile) {
	if len(profile.RequestHistory) == 0 {
		return
	}

	// Calculate average session duration
	totalDuration := int64(0)
	successCount := 0
	for _, req := range profile.RequestHistory {
		totalDuration += req.SessionDuration
		if req.Success {
			successCount++
		}
	}
	profile.BehaviorMetrics.AverageSessionDuration = float64(totalDuration) / float64(len(profile.RequestHistory))

	// Calculate success rate
	profile.BehaviorMetrics.SuccessRate = float64(successCount) / float64(len(profile.RequestHistory))

	// Calculate request frequency (requests per day)
	if len(profile.RequestHistory) > 1 {
		firstRequest := profile.RequestHistory[0].Timestamp
		lastRequest := profile.RequestHistory[len(profile.RequestHistory)-1].Timestamp
		daysDiff := lastRequest.Sub(firstRequest).Hours() / 24
		if daysDiff > 0 {
			profile.BehaviorMetrics.RequestFrequency = float64(len(profile.RequestHistory)) / daysDiff
		}
	}

	// Calculate device consistency
	uniqueDevices := len(profile.DeviceHistory)
	totalRequests := len(profile.RequestHistory)
	profile.BehaviorMetrics.DeviceConsistency = 1.0 - (float64(uniqueDevices)/float64(totalRequests))

	// Calculate location consistency
	uniqueCountries := ba.countUniqueCountries(profile.LocationHistory)
	profile.BehaviorMetrics.LocationConsistency = 1.0 - (float64(uniqueCountries)/float64(len(profile.LocationHistory)))
}

// calculateBehavioralRisk calculates the behavioral risk score
func (ba *BehavioralAnalyzer) calculateBehavioralRisk(profile *UserBehaviorProfile, current *BehavioralSignals) float64 {
	riskScore := 0.0

	// Check session duration deviation
	if profile.BehaviorMetrics.AverageSessionDuration > 0 {
		currentDuration := float64(current.SessionData.SessionDuration)
		deviationRatio := math.Abs(currentDuration-profile.BehaviorMetrics.AverageSessionDuration) / profile.BehaviorMetrics.AverageSessionDuration
		if deviationRatio > 2.0 { // More than 200% deviation
			riskScore += 30.0
		} else if deviationRatio > 1.0 { // More than 100% deviation
			riskScore += 15.0
		}
	}

	// Check device consistency
	if profile.BehaviorMetrics.DeviceConsistency < 0.5 { // Low device consistency
		riskScore += 25.0
	}

	// Check location consistency
	if profile.BehaviorMetrics.LocationConsistency < 0.3 { // Very low location consistency
		riskScore += 35.0
	}

	// Check request frequency
	if profile.BehaviorMetrics.RequestFrequency > ba.model.globalStats.AverageRequestFreq*3 {
		riskScore += 40.0
	}

	// Check for new device
	if !ba.contains(profile.DeviceHistory, current.DeviceFingerprint) {
		riskScore += 20.0
	}

	// Add some randomness to simulate ML model uncertainty
	riskScore += rand.Float64() * 10.0

	return math.Min(riskScore, 100.0)
}

// Helper functions

// contains checks if a slice contains a string
func (ba *BehavioralAnalyzer) contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// countUniqueCountries counts unique countries in location history
func (ba *BehavioralAnalyzer) countUniqueCountries(locations []GeoLocation) int {
	countries := make(map[string]bool)
	for _, loc := range locations {
		countries[loc.Country] = true
	}
	return len(countries)
}