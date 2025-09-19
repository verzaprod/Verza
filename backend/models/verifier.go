package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Verifier represents a registered verifier in the marketplace
type Verifier struct {
	ID                    string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID                string         `json:"user_id" gorm:"not null;index"`
	HederaAccountID       string         `json:"hedera_account_id" gorm:"not null;uniqueIndex"`
	BusinessName          string         `json:"business_name" gorm:"not null"`
	BusinessType          string         `json:"business_type" gorm:"not null"`
	RegistrationNumber    string         `json:"registration_number"`
	LicenseNumber         string         `json:"license_number"`
	ContactEmail          string         `json:"contact_email" gorm:"not null"`
	ContactPhone          string         `json:"contact_phone"`
	Website               string         `json:"website"`
	Description           string         `json:"description" gorm:"type:text"`
	Specializations       pq.StringArray `json:"specializations" gorm:"type:text[]"`
	SupportedDocuments    pq.StringArray `json:"supported_documents" gorm:"type:text[]"`
	OperatingCountries    pq.StringArray `json:"operating_countries" gorm:"type:text[]"`
	OperatingHours        string         `json:"operating_hours"`
	TimeZone              string         `json:"time_zone"`
	Languages             pq.StringArray `json:"languages" gorm:"type:text[]"`
	Status                string         `json:"status" gorm:"not null;default:'pending'"` // pending, active, suspended, deactivated
	VerificationLevel     string         `json:"verification_level" gorm:"not null;default:'basic'"` // basic, standard, premium
	KYCStatus             string         `json:"kyc_status" gorm:"not null;default:'pending'"` // pending, verified, rejected
	KYCDocuments          JSONField      `json:"kyc_documents" gorm:"type:jsonb"`
	ComplianceCertificates JSONField     `json:"compliance_certificates" gorm:"type:jsonb"`
	InsuranceInfo         JSONField      `json:"insurance_info" gorm:"type:jsonb"`
	StakeInfo             VerifierStake  `json:"stake_info" gorm:"embedded;embeddedPrefix:stake_"`
	ReputationScore       float64        `json:"reputation_score" gorm:"default:0"`
	TotalVerifications    int64          `json:"total_verifications" gorm:"default:0"`
	SuccessfulVerifications int64        `json:"successful_verifications" gorm:"default:0"`
	SuccessRate           float64        `json:"success_rate" gorm:"default:0"`
	AverageResponseTime   int            `json:"average_response_time" gorm:"default:0"` // in minutes
	LastActiveAt          *time.Time     `json:"last_active_at"`
	RegisteredAt          time.Time      `json:"registered_at" gorm:"not null;default:now()"`
	ApprovedAt            *time.Time     `json:"approved_at"`
	ApprovedBy            *string        `json:"approved_by"`
	SuspendedAt           *time.Time     `json:"suspended_at"`
	SuspendedBy           *string        `json:"suspended_by"`
	SuspensionReason      *string        `json:"suspension_reason"`
	Settings              JSONField      `json:"settings" gorm:"type:jsonb"`
	Metadata              JSONField      `json:"metadata" gorm:"type:jsonb"`
	CreatedAt             time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt             time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt             gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	User                  *User                   `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Verifications         []VerificationRequest   `json:"verifications,omitempty" gorm:"foreignKey:VerifierID"`
	ReputationHistory     []VerifierReputation    `json:"reputation_history,omitempty" gorm:"foreignKey:VerifierID"`
	StakeTransactions     []VerifierStakeTransaction `json:"stake_transactions,omitempty" gorm:"foreignKey:VerifierID"`
	PricingRules          []VerifierPricing       `json:"pricing_rules,omitempty" gorm:"foreignKey:VerifierID"`
}

// VerifierStake represents staking information for a verifier
type VerifierStake struct {
	RequiredAmount    float64   `json:"required_amount" gorm:"not null;default:1000"`
	CurrentAmount     float64   `json:"current_amount" gorm:"default:0"`
	LockedAmount      float64   `json:"locked_amount" gorm:"default:0"`
	AvailableAmount   float64   `json:"available_amount" gorm:"default:0"`
	Currency          string    `json:"currency" gorm:"not null;default:'HBAR'"`
	StakeStatus       string    `json:"stake_status" gorm:"not null;default:'insufficient'"` // insufficient, sufficient, locked, slashed
	LastStakeUpdate   time.Time `json:"last_stake_update" gorm:"default:now()"`
	SlashedAmount     float64   `json:"slashed_amount" gorm:"default:0"`
	TotalSlashed      float64   `json:"total_slashed" gorm:"default:0"`
	SlashingHistory   JSONField `json:"slashing_history" gorm:"type:jsonb"`
}

// VerifierReputation tracks reputation changes over time
type VerifierReputation struct {
	ID                string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	VerifierID        string    `json:"verifier_id" gorm:"not null;index"`
	VerificationID    string    `json:"verification_id" gorm:"not null;index"`
	PreviousScore     float64   `json:"previous_score"`
	NewScore          float64   `json:"new_score"`
	ScoreChange       float64   `json:"score_change"`
	ChangeReason      string    `json:"change_reason" gorm:"not null"`
	ChangeType        string    `json:"change_type" gorm:"not null"` // increase, decrease, penalty, bonus
	ImpactFactor      float64   `json:"impact_factor" gorm:"default:1.0"`
	CalculatedBy      string    `json:"calculated_by" gorm:"not null;default:'system'"`
	Notes             string    `json:"notes" gorm:"type:text"`
	Metadata          JSONField `json:"metadata" gorm:"type:jsonb"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	Verifier          *Verifier           `json:"verifier,omitempty" gorm:"foreignKey:VerifierID"`
	Verification      *VerificationRequest `json:"verification,omitempty" gorm:"foreignKey:VerificationID"`
}

// VerifierStakeTransaction tracks all stake-related transactions
type VerifierStakeTransaction struct {
	ID                string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	VerifierID        string    `json:"verifier_id" gorm:"not null;index"`
	TransactionType   string    `json:"transaction_type" gorm:"not null"` // deposit, withdrawal, lock, unlock, slash, reward
	Amount            float64   `json:"amount" gorm:"not null"`
	Currency          string    `json:"currency" gorm:"not null;default:'HBAR'"`
	HederaTxID        string    `json:"hedera_tx_id"`
	HederaTxHash      string    `json:"hedera_tx_hash"`
	Status            string    `json:"status" gorm:"not null;default:'pending'"` // pending, confirmed, failed, cancelled
	Reason            string    `json:"reason"`
	RelatedEntityID   *string   `json:"related_entity_id"` // verification_id, dispute_id, etc.
	RelatedEntityType *string   `json:"related_entity_type"` // verification, dispute, penalty, etc.
	BalanceBefore     float64   `json:"balance_before"`
	BalanceAfter      float64   `json:"balance_after"`
	ProcessedBy       string    `json:"processed_by" gorm:"default:'system'"`
	ProcessedAt       *time.Time `json:"processed_at"`
	FailureReason     *string   `json:"failure_reason"`
	Metadata          JSONField `json:"metadata" gorm:"type:jsonb"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Verifier          *Verifier `json:"verifier,omitempty" gorm:"foreignKey:VerifierID"`
}

// VerifierPricing defines pricing rules for different verification types
type VerifierPricing struct {
	ID                string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	VerifierID        string    `json:"verifier_id" gorm:"not null;index"`
	DocumentType      string    `json:"document_type" gorm:"not null"`
	VerificationType  string    `json:"verification_type" gorm:"not null"`
	BasePrice         float64   `json:"base_price" gorm:"not null"`
	Currency          string    `json:"currency" gorm:"not null;default:'HBAR'"`
	RushPrice         *float64  `json:"rush_price"` // Additional fee for rush processing
	BulkDiscounts     JSONField `json:"bulk_discounts" gorm:"type:jsonb"` // Volume-based discounts
	GeographicPricing JSONField `json:"geographic_pricing" gorm:"type:jsonb"` // Country/region-specific pricing
	ComplexityMultiplier float64 `json:"complexity_multiplier" gorm:"default:1.0"`
	MinimumPrice      float64   `json:"minimum_price" gorm:"default:0"`
	MaximumPrice      *float64  `json:"maximum_price"`
	IsActive          bool      `json:"is_active" gorm:"default:true"`
	EffectiveFrom     time.Time `json:"effective_from" gorm:"not null;default:now()"`
	EffectiveUntil    *time.Time `json:"effective_until"`
	Notes             string    `json:"notes" gorm:"type:text"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Verifier          *Verifier `json:"verifier,omitempty" gorm:"foreignKey:VerifierID"`
}

// VerifierAvailability tracks verifier availability and capacity
type VerifierAvailability struct {
	ID                string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	VerifierID        string    `json:"verifier_id" gorm:"not null;index"`
	Date              time.Time `json:"date" gorm:"not null;index"`
	MaxCapacity       int       `json:"max_capacity" gorm:"not null;default:10"`
	CurrentLoad       int       `json:"current_load" gorm:"default:0"`
	AvailableSlots    int       `json:"available_slots" gorm:"default:10"`
	IsAvailable       bool      `json:"is_available" gorm:"default:true"`
	UnavailableReason *string   `json:"unavailable_reason"`
	SpecialHours      JSONField `json:"special_hours" gorm:"type:jsonb"` // Override normal operating hours
	HolidaySchedule   JSONField `json:"holiday_schedule" gorm:"type:jsonb"`
	EmergencyContact  JSONField `json:"emergency_contact" gorm:"type:jsonb"`
	AutoAcceptLimit   int       `json:"auto_accept_limit" gorm:"default:5"`
	PriorityBookings  JSONField `json:"priority_bookings" gorm:"type:jsonb"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Verifier          *Verifier `json:"verifier,omitempty" gorm:"foreignKey:VerifierID"`
}

// VerifierPerformanceMetrics tracks detailed performance metrics
type VerifierPerformanceMetrics struct {
	ID                     string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	VerifierID             string    `json:"verifier_id" gorm:"not null;index"`
	PeriodStart            time.Time `json:"period_start" gorm:"not null"`
	PeriodEnd              time.Time `json:"period_end" gorm:"not null"`
	TotalRequests          int64     `json:"total_requests" gorm:"default:0"`
	CompletedRequests      int64     `json:"completed_requests" gorm:"default:0"`
	RejectedRequests       int64     `json:"rejected_requests" gorm:"default:0"`
	CancelledRequests      int64     `json:"cancelled_requests" gorm:"default:0"`
	AverageResponseTime    float64   `json:"average_response_time" gorm:"default:0"` // in hours
	MedianResponseTime     float64   `json:"median_response_time" gorm:"default:0"`
	FastestResponseTime    float64   `json:"fastest_response_time" gorm:"default:0"`
	SlowestResponseTime    float64   `json:"slowest_response_time" gorm:"default:0"`
	AccuracyScore          float64   `json:"accuracy_score" gorm:"default:0"`
	CustomerSatisfaction   float64   `json:"customer_satisfaction" gorm:"default:0"`
	DisputeRate            float64   `json:"dispute_rate" gorm:"default:0"`
	FraudDetectionRate     float64   `json:"fraud_detection_rate" gorm:"default:0"`
	RevenueGenerated       float64   `json:"revenue_generated" gorm:"default:0"`
	PenaltiesIncurred      float64   `json:"penalties_incurred" gorm:"default:0"`
	BonusesEarned          float64   `json:"bonuses_earned" gorm:"default:0"`
	QualityScore           float64   `json:"quality_score" gorm:"default:0"`
	ReliabilityScore       float64   `json:"reliability_score" gorm:"default:0"`
	EfficiencyScore        float64   `json:"efficiency_score" gorm:"default:0"`
	OverallRating          float64   `json:"overall_rating" gorm:"default:0"`
	RankInCategory         int       `json:"rank_in_category" gorm:"default:0"`
	RankOverall            int       `json:"rank_overall" gorm:"default:0"`
	TrendDirection         string    `json:"trend_direction" gorm:"default:'stable'"` // improving, declining, stable
	PerformanceNotes       string    `json:"performance_notes" gorm:"type:text"`
	CalculatedAt           time.Time `json:"calculated_at" gorm:"not null;default:now()"`
	Metadata               JSONField `json:"metadata" gorm:"type:jsonb"`

	// Relationships
	Verifier               *Verifier `json:"verifier,omitempty" gorm:"foreignKey:VerifierID"`
}

// VerifierCertification tracks verifier certifications and credentials
type VerifierCertification struct {
	ID                string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	VerifierID        string     `json:"verifier_id" gorm:"not null;index"`
	CertificationType string     `json:"certification_type" gorm:"not null"`
	CertificationName string     `json:"certification_name" gorm:"not null"`
	IssuingAuthority  string     `json:"issuing_authority" gorm:"not null"`
	CertificateNumber string     `json:"certificate_number"`
	IssuedDate        time.Time  `json:"issued_date" gorm:"not null"`
	ExpiryDate        *time.Time `json:"expiry_date"`
	Status            string     `json:"status" gorm:"not null;default:'active'"` // active, expired, revoked, suspended
	VerificationStatus string    `json:"verification_status" gorm:"not null;default:'pending'"` // pending, verified, rejected
	DocumentURL       string     `json:"document_url"`
	DocumentHash      string     `json:"document_hash"`
	Scope             pq.StringArray `json:"scope" gorm:"type:text[]"` // What this certification covers
	Level             string     `json:"level"` // basic, intermediate, advanced, expert
	VerifiedBy        *string    `json:"verified_by"`
	VerifiedAt        *time.Time `json:"verified_at"`
	RejectionReason   *string    `json:"rejection_reason"`
	RenewalReminders  JSONField  `json:"renewal_reminders" gorm:"type:jsonb"`
	Notes             string     `json:"notes" gorm:"type:text"`
	Metadata          JSONField  `json:"metadata" gorm:"type:jsonb"`
	CreatedAt         time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Verifier          *Verifier  `json:"verifier,omitempty" gorm:"foreignKey:VerifierID"`
}

// BeforeCreate hook for Verifier
func (v *Verifier) BeforeCreate(tx *gorm.DB) error {
	// Set default values
	if v.Status == "" {
		v.Status = "pending"
	}
	if v.VerificationLevel == "" {
		v.VerificationLevel = "basic"
	}
	if v.KYCStatus == "" {
		v.KYCStatus = "pending"
	}
	if v.StakeInfo.RequiredAmount == 0 {
		v.StakeInfo.RequiredAmount = 1000 // Default stake requirement
	}
	if v.StakeInfo.Currency == "" {
		v.StakeInfo.Currency = "HBAR"
	}
	if v.StakeInfo.StakeStatus == "" {
		v.StakeInfo.StakeStatus = "insufficient"
	}
	return nil
}

// BeforeUpdate hook for Verifier
func (v *Verifier) BeforeUpdate(tx *gorm.DB) error {
	// Update calculated fields
	if v.TotalVerifications > 0 {
		v.SuccessRate = float64(v.SuccessfulVerifications) / float64(v.TotalVerifications) * 100
	}

	// Update stake available amount
	v.StakeInfo.AvailableAmount = v.StakeInfo.CurrentAmount - v.StakeInfo.LockedAmount

	// Update stake status based on current amount
	if v.StakeInfo.CurrentAmount >= v.StakeInfo.RequiredAmount {
		if v.StakeInfo.StakeStatus == "insufficient" {
			v.StakeInfo.StakeStatus = "sufficient"
		}
	} else {
		v.StakeInfo.StakeStatus = "insufficient"
	}

	return nil
}

// TableName returns the table name for Verifier
func (Verifier) TableName() string {
	return "verifiers"
}

// TableName returns the table name for VerifierReputation
func (VerifierReputation) TableName() string {
	return "verifier_reputation_history"
}

// TableName returns the table name for VerifierStakeTransaction
func (VerifierStakeTransaction) TableName() string {
	return "verifier_stake_transactions"
}

// TableName returns the table name for VerifierPricing
func (VerifierPricing) TableName() string {
	return "verifier_pricing_rules"
}

// TableName returns the table name for VerifierAvailability
func (VerifierAvailability) TableName() string {
	return "verifier_availability"
}

// TableName returns the table name for VerifierPerformanceMetrics
func (VerifierPerformanceMetrics) TableName() string {
	return "verifier_performance_metrics"
}

// TableName returns the table name for VerifierCertification
func (VerifierCertification) TableName() string {
	return "verifier_certifications"
}