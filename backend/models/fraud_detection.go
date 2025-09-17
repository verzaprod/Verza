package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// FraudDetectionResult represents a fraud analysis result stored in database
type FraudDetectionResult struct {
	ID               uint                   `json:"id" gorm:"primaryKey"`
	RequestID        string                 `json:"request_id" gorm:"uniqueIndex;not null"`
	UserID           string                 `json:"user_id" gorm:"index;not null"`
	OverallRiskScore float64                `json:"overall_risk_score" gorm:"not null"`
	RiskLevel        string                 `json:"risk_level" gorm:"not null"`
	FraudFlags       FraudFlagsJSON         `json:"fraud_flags" gorm:"type:json"`
	ComponentScores  ComponentScoresJSON    `json:"component_scores" gorm:"type:json"`
	DocumentData     DocumentMetadataJSON   `json:"document_data" gorm:"type:json"`
	BiometricData    BiometricDataJSON      `json:"biometric_data" gorm:"type:json"`
	BehavioralData   BehavioralDataJSON     `json:"behavioral_data" gorm:"type:json"`
	Recommendation   string                 `json:"recommendation" gorm:"not null"`
	Confidence       float64                `json:"confidence" gorm:"not null"`
	ProcessingTime   int64                  `json:"processing_time_ms" gorm:"not null"`
	Status           string                 `json:"status" gorm:"default:'pending'"` // pending, approved, rejected, manual_review
	ReviewedBy       *string                `json:"reviewed_by" gorm:"index"`
	ReviewedAt       *time.Time             `json:"reviewed_at"`
	ReviewNotes      *string                `json:"review_notes"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

// UserBehaviorProfile represents user behavioral patterns
type UserBehaviorProfile struct {
	ID                     uint                    `json:"id" gorm:"primaryKey"`
	UserID                 string                  `json:"user_id" gorm:"uniqueIndex;not null"`
	RequestHistory         RequestHistoryJSON      `json:"request_history" gorm:"type:json"`
	DeviceHistory          DeviceHistoryJSON       `json:"device_history" gorm:"type:json"`
	LocationHistory        LocationHistoryJSON     `json:"location_history" gorm:"type:json"`
	BehaviorMetrics        BehaviorMetricsJSON     `json:"behavior_metrics" gorm:"type:json"`
	RiskScore              float64                 `json:"risk_score" gorm:"default:0"`
	TotalRequests          int                     `json:"total_requests" gorm:"default:0"`
	SuccessfulRequests     int                     `json:"successful_requests" gorm:"default:0"`
	FraudulentRequests     int                     `json:"fraudulent_requests" gorm:"default:0"`
	LastRequestAt          *time.Time              `json:"last_request_at"`
	CreatedAt              time.Time               `json:"created_at"`
	UpdatedAt              time.Time               `json:"updated_at"`
}

// BiometricHash represents stored biometric hashes for duplicate detection
type BiometricHash struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Hash          string    `json:"hash" gorm:"uniqueIndex;not null"`
	UserID        string    `json:"user_id" gorm:"index;not null"`
	UsageCount    int       `json:"usage_count" gorm:"default:1"`
	FirstUsedAt   time.Time `json:"first_used_at"`
	LastUsedAt    time.Time `json:"last_used_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// FraudPattern represents detected fraud patterns for ML training
type FraudPattern struct {
	ID              uint                  `json:"id" gorm:"primaryKey"`
	PatternType     string                `json:"pattern_type" gorm:"index;not null"`
	PatternData     PatternDataJSON       `json:"pattern_data" gorm:"type:json"`
	Severity        string                `json:"severity" gorm:"not null"`
	Confidence      float64               `json:"confidence" gorm:"not null"`
	OccurrenceCount int                   `json:"occurrence_count" gorm:"default:1"`
	LastSeen        time.Time             `json:"last_seen"`
	IsActive        bool                  `json:"is_active" gorm:"default:true"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
}

// MLModelMetrics represents ML model performance metrics
type MLModelMetrics struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	ModelName         string    `json:"model_name" gorm:"uniqueIndex;not null"`
	ModelVersion      string    `json:"model_version" gorm:"not null"`
	Accuracy          float64   `json:"accuracy"`
	Precision         float64   `json:"precision"`
	Recall            float64   `json:"recall"`
	F1Score           float64   `json:"f1_score"`
	FalsePositiveRate float64   `json:"false_positive_rate"`
	FalseNegativeRate float64   `json:"false_negative_rate"`
	TotalPredictions  int       `json:"total_predictions" gorm:"default:0"`
	CorrectPredictions int      `json:"correct_predictions" gorm:"default:0"`
	LastUpdated       time.Time `json:"last_updated"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// FraudAlert represents fraud alerts for monitoring
type FraudAlert struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	AlertType   string    `json:"alert_type" gorm:"index;not null"`
	Severity    string    `json:"severity" gorm:"not null"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description" gorm:"not null"`
	UserID      *string   `json:"user_id" gorm:"index"`
	RequestID   *string   `json:"request_id" gorm:"index"`
	Metadata    AlertMetadataJSON `json:"metadata" gorm:"type:json"`
	Status      string    `json:"status" gorm:"default:'active'"` // active, resolved, dismissed
	ResolvedBy  *string   `json:"resolved_by"`
	ResolvedAt  *time.Time `json:"resolved_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// JSON field types for GORM

// FraudFlagsJSON handles JSON marshaling for fraud flags
type FraudFlagsJSON []FraudFlag

type FraudFlag struct {
	Type        string  `json:"type"`
	Severity    string  `json:"severity"`
	Description string  `json:"description"`
	Score       float64 `json:"score"`
	Evidence    string  `json:"evidence"`
}

func (f FraudFlagsJSON) Value() (driver.Value, error) {
	return json.Marshal(f)
}

func (f *FraudFlagsJSON) Scan(value interface{}) error {
	if value == nil {
		*f = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into FraudFlagsJSON", value)
	}
	return json.Unmarshal(bytes, f)
}

// ComponentScoresJSON handles JSON marshaling for component scores
type ComponentScoresJSON struct {
	DocumentRisk   float64 `json:"document_risk"`
	BiometricRisk  float64 `json:"biometric_risk"`
	BehavioralRisk float64 `json:"behavioral_risk"`
	AnomalyScore   float64 `json:"anomaly_score"`
}

func (c ComponentScoresJSON) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *ComponentScoresJSON) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into ComponentScoresJSON", value)
	}
	return json.Unmarshal(bytes, c)
}

// DocumentMetadataJSON handles JSON marshaling for document metadata
type DocumentMetadataJSON struct {
	DocumentType    string    `json:"document_type"`
	DocumentHash    string    `json:"document_hash"`
	IssuanceDate    time.Time `json:"issuance_date"`
	ExpirationDate  time.Time `json:"expiration_date"`
	CountryOfOrigin string    `json:"country_of_origin"`
	DocumentNumber  string    `json:"document_number"`
	ImageQuality    float64   `json:"image_quality"`
	OCRConfidence   float64   `json:"ocr_confidence"`
}

func (d DocumentMetadataJSON) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (d *DocumentMetadataJSON) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into DocumentMetadataJSON", value)
	}
	return json.Unmarshal(bytes, d)
}

// BiometricDataJSON handles JSON marshaling for biometric data
type BiometricDataJSON struct {
	FaceMatchScore   float64 `json:"face_match_score"`
	LivenessScore    float64 `json:"liveness_score"`
	FaceQuality      float64 `json:"face_quality"`
	BiometricHash    string  `json:"biometric_hash"`
	VerificationTime int64   `json:"verification_time_ms"`
}

func (b BiometricDataJSON) Value() (driver.Value, error) {
	return json.Marshal(b)
}

func (b *BiometricDataJSON) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into BiometricDataJSON", value)
	}
	return json.Unmarshal(bytes, b)
}

// BehavioralDataJSON handles JSON marshaling for behavioral data
type BehavioralDataJSON struct {
	IPAddress        string           `json:"ip_address"`
	UserAgent        string           `json:"user_agent"`
	DeviceFingerprint string          `json:"device_fingerprint"`
	GeoLocation      GeoLocationJSON  `json:"geo_location"`
	RequestFrequency RequestFreqJSON  `json:"request_frequency"`
	SessionData      SessionDataJSON  `json:"session_data"`
}

type GeoLocationJSON struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Country   string  `json:"country"`
	City      string  `json:"city"`
	ISP       string  `json:"isp"`
	VPNRisk   float64 `json:"vpn_risk"`
}

type RequestFreqJSON struct {
	RequestsLast24h int     `json:"requests_last_24h"`
	RequestsLastWeek int    `json:"requests_last_week"`
	AverageInterval float64 `json:"average_interval_minutes"`
	BurstPattern    bool    `json:"burst_pattern"`
}

type SessionDataJSON struct {
	SessionDuration int64   `json:"session_duration_ms"`
	PageViews       int     `json:"page_views"`
	MouseMovements  int     `json:"mouse_movements"`
	Keystrokes      int     `json:"keystrokes"`
	BotScore        float64 `json:"bot_score"`
}

func (b BehavioralDataJSON) Value() (driver.Value, error) {
	return json.Marshal(b)
}

func (b *BehavioralDataJSON) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into BehavioralDataJSON", value)
	}
	return json.Unmarshal(bytes, b)
}

// RequestHistoryJSON handles JSON marshaling for request history
type RequestHistoryJSON []RequestPattern

type RequestPattern struct {
	Timestamp       time.Time `json:"timestamp"`
	RequestType     string    `json:"request_type"`
	SessionDuration int64     `json:"session_duration"`
	Success         bool      `json:"success"`
}

func (r RequestHistoryJSON) Value() (driver.Value, error) {
	return json.Marshal(r)
}

func (r *RequestHistoryJSON) Scan(value interface{}) error {
	if value == nil {
		*r = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into RequestHistoryJSON", value)
	}
	return json.Unmarshal(bytes, r)
}

// DeviceHistoryJSON handles JSON marshaling for device history
type DeviceHistoryJSON []string

func (d DeviceHistoryJSON) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (d *DeviceHistoryJSON) Scan(value interface{}) error {
	if value == nil {
		*d = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into DeviceHistoryJSON", value)
	}
	return json.Unmarshal(bytes, d)
}

// LocationHistoryJSON handles JSON marshaling for location history
type LocationHistoryJSON []GeoLocationJSON

func (l LocationHistoryJSON) Value() (driver.Value, error) {
	return json.Marshal(l)
}

func (l *LocationHistoryJSON) Scan(value interface{}) error {
	if value == nil {
		*l = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into LocationHistoryJSON", value)
	}
	return json.Unmarshal(bytes, l)
}

// BehaviorMetricsJSON handles JSON marshaling for behavior metrics
type BehaviorMetricsJSON struct {
	AverageSessionDuration float64 `json:"average_session_duration"`
	RequestFrequency       float64 `json:"request_frequency"`
	SuccessRate           float64 `json:"success_rate"`
	DeviceConsistency     float64 `json:"device_consistency"`
	LocationConsistency   float64 `json:"location_consistency"`
}

func (b BehaviorMetricsJSON) Value() (driver.Value, error) {
	return json.Marshal(b)
}

func (b *BehaviorMetricsJSON) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into BehaviorMetricsJSON", value)
	}
	return json.Unmarshal(bytes, b)
}

// PatternDataJSON handles JSON marshaling for pattern data
type PatternDataJSON map[string]interface{}

func (p PatternDataJSON) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *PatternDataJSON) Scan(value interface{}) error {
	if value == nil {
		*p = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into PatternDataJSON", value)
	}
	return json.Unmarshal(bytes, p)
}

// AlertMetadataJSON handles JSON marshaling for alert metadata
type AlertMetadataJSON map[string]interface{}

func (a AlertMetadataJSON) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *AlertMetadataJSON) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into AlertMetadataJSON", value)
	}
	return json.Unmarshal(bytes, a)
}

// Table names
func (FraudDetectionResult) TableName() string {
	return "fraud_detection_results"
}

func (UserBehaviorProfile) TableName() string {
	return "user_behavior_profiles"
}

func (BiometricHash) TableName() string {
	return "biometric_hashes"
}

func (FraudPattern) TableName() string {
	return "fraud_patterns"
}

func (MLModelMetrics) TableName() string {
	return "ml_model_metrics"
}

func (FraudAlert) TableName() string {
	return "fraud_alerts"
}

// Model hooks

// BeforeCreate hook for FraudDetectionResult
func (fdr *FraudDetectionResult) BeforeCreate(tx *gorm.DB) error {
	if fdr.RequestID == "" {
		return fmt.Errorf("request_id is required")
	}
	if fdr.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	return nil
}

// BeforeCreate hook for UserBehaviorProfile
func (ubp *UserBehaviorProfile) BeforeCreate(tx *gorm.DB) error {
	if ubp.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	return nil
}

// BeforeCreate hook for BiometricHash
func (bh *BiometricHash) BeforeCreate(tx *gorm.DB) error {
	if bh.Hash == "" {
		return fmt.Errorf("hash is required")
	}
	if bh.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	bh.FirstUsedAt = time.Now()
	bh.LastUsedAt = time.Now()
	return nil
}

// BeforeUpdate hook for BiometricHash
func (bh *BiometricHash) BeforeUpdate(tx *gorm.DB) error {
	bh.LastUsedAt = time.Now()
	return nil
}