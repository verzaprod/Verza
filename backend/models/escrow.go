package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// EscrowTransaction represents an escrow transaction for verification requests
type EscrowTransaction struct {
	ID                    string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	VerificationRequestID string         `json:"verification_request_id" gorm:"not null;uniqueIndex"`
	UserID                string         `json:"user_id" gorm:"not null;index"`
	VerifierID            *string        `json:"verifier_id" gorm:"index"`
	Amount                float64        `json:"amount" gorm:"not null"`
	Currency              string         `json:"currency" gorm:"not null;default:'HBAR'"`
	FeeAmount             float64        `json:"fee_amount" gorm:"default:0"`
	TotalAmount           float64        `json:"total_amount" gorm:"not null"`
	Status                string         `json:"status" gorm:"not null;default:'pending'"` // pending, locked, released, refunded, disputed, cancelled
	EscrowType            string         `json:"escrow_type" gorm:"not null;default:'verification'"` // verification, dispute_resolution, penalty
	HederaAccountID       string         `json:"hedera_account_id" gorm:"not null"`
	HederaTxID            *string        `json:"hedera_tx_id"`
	HederaTxHash          *string        `json:"hedera_tx_hash"`
	LockTxID              *string        `json:"lock_tx_id"`
	LockTxHash            *string        `json:"lock_tx_hash"`
	ReleaseTxID           *string        `json:"release_tx_id"`
	ReleaseTxHash         *string        `json:"release_tx_hash"`
	RefundTxID            *string        `json:"refund_tx_id"`
	RefundTxHash          *string        `json:"refund_tx_hash"`
	AutoReleaseEnabled    bool           `json:"auto_release_enabled" gorm:"default:true"`
	AutoReleaseDelay      int            `json:"auto_release_delay" gorm:"default:72"` // hours
	AutoReleaseAt         *time.Time     `json:"auto_release_at"`
	ManualReleaseRequired bool           `json:"manual_release_required" gorm:"default:false"`
	DisputeDeadline       *time.Time     `json:"dispute_deadline"`
	FraudCheckRequired    bool           `json:"fraud_check_required" gorm:"default:true"`
	FraudCheckStatus      string         `json:"fraud_check_status" gorm:"default:'pending'"` // pending, passed, failed, skipped
	FraudCheckResult      JSONField      `json:"fraud_check_result" gorm:"type:jsonb"`
	ReleaseConditions     JSONField      `json:"release_conditions" gorm:"type:jsonb"`
	RefundConditions      JSONField      `json:"refund_conditions" gorm:"type:jsonb"`
	LockedAt              *time.Time     `json:"locked_at"`
	ReleasedAt            *time.Time     `json:"released_at"`
	RefundedAt            *time.Time     `json:"refunded_at"`
	CancelledAt           *time.Time     `json:"cancelled_at"`
	CancellationReason    *string        `json:"cancellation_reason"`
	ProcessedBy           *string        `json:"processed_by"`
	Notes                 string         `json:"notes" gorm:"type:text"`
	Metadata              JSONField      `json:"metadata" gorm:"type:jsonb"`
	CreatedAt             time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt             time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt             gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	User                  *User                   `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Verifier              *Verifier               `json:"verifier,omitempty" gorm:"foreignKey:VerifierID"`
	VerificationRequest   *VerificationRequest    `json:"verification_request,omitempty" gorm:"foreignKey:VerificationRequestID"`
	StatusHistory         []EscrowStatusHistory   `json:"status_history,omitempty" gorm:"foreignKey:EscrowTransactionID"`
	Disputes              []EscrowDispute         `json:"disputes,omitempty" gorm:"foreignKey:EscrowTransactionID"`
	Payments              []EscrowPayment         `json:"payments,omitempty" gorm:"foreignKey:EscrowTransactionID"`
}

// EscrowStatusHistory tracks status changes for escrow transactions
type EscrowStatusHistory struct {
	ID                   string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	EscrowTransactionID  string    `json:"escrow_transaction_id" gorm:"not null;index"`
	PreviousStatus       *string   `json:"previous_status"`
	NewStatus            string    `json:"new_status" gorm:"not null"`
	ChangeReason         string    `json:"change_reason" gorm:"not null"`
	ChangedBy            string    `json:"changed_by" gorm:"not null"`
	ChangedByType        string    `json:"changed_by_type" gorm:"not null;default:'system'"` // system, user, admin, verifier
	Automated            bool      `json:"automated" gorm:"default:false"`
	TriggerEvent         *string   `json:"trigger_event"`
	RelatedTransactionID *string   `json:"related_transaction_id"`
	Notes                string    `json:"notes" gorm:"type:text"`
	Metadata             JSONField `json:"metadata" gorm:"type:jsonb"`
	CreatedAt            time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	EscrowTransaction    *EscrowTransaction `json:"escrow_transaction,omitempty" gorm:"foreignKey:EscrowTransactionID"`
}

// EscrowDispute represents disputes related to escrow transactions
type EscrowDispute struct {
	ID                  string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	EscrowTransactionID string         `json:"escrow_transaction_id" gorm:"not null;index"`
	DisputeType         string         `json:"dispute_type" gorm:"not null"` // quality, fraud, non_delivery, technical, other
	InitiatedBy         string         `json:"initiated_by" gorm:"not null"`
	InitiatorType       string         `json:"initiator_type" gorm:"not null"` // user, verifier, system
	Title               string         `json:"title" gorm:"not null"`
	Description         string         `json:"description" gorm:"type:text;not null"`
	Evidence            JSONField      `json:"evidence" gorm:"type:jsonb"`
	Severity            string         `json:"severity" gorm:"not null;default:'medium'"` // low, medium, high, critical
	Status              string         `json:"status" gorm:"not null;default:'open'"` // open, investigating, resolved, closed, escalated
	Priority            string         `json:"priority" gorm:"not null;default:'normal'"` // low, normal, high, urgent
	AssignedTo          *string        `json:"assigned_to"`
	AssignedAt          *time.Time     `json:"assigned_at"`
	Resolution          *string        `json:"resolution" gorm:"type:text"`
	ResolutionType      *string        `json:"resolution_type"` // refund_user, pay_verifier, partial_refund, no_action, escalate
	ResolvedBy          *string        `json:"resolved_by"`
	ResolvedAt          *time.Time     `json:"resolved_at"`
	EscalatedTo         *string        `json:"escalated_to"`
	EscalatedAt         *time.Time     `json:"escalated_at"`
	EscalationReason    *string        `json:"escalation_reason"`
	AutoResolution      bool           `json:"auto_resolution" gorm:"default:false"`
	ResolutionDeadline  *time.Time     `json:"resolution_deadline"`
	CommunicationLog    JSONField      `json:"communication_log" gorm:"type:jsonb"`
	Tags                pq.StringArray `json:"tags" gorm:"type:text[]"`
	Metadata            JSONField      `json:"metadata" gorm:"type:jsonb"`
	CreatedAt           time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt           time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt           gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	EscrowTransaction   *EscrowTransaction    `json:"escrow_transaction,omitempty" gorm:"foreignKey:EscrowTransactionID"`
	Evidence_Records    []EscrowDisputeEvidence `json:"evidence_records,omitempty" gorm:"foreignKey:DisputeID"`
	ResolutionHistory   []EscrowDisputeResolution `json:"resolution_history,omitempty" gorm:"foreignKey:DisputeID"`
}

// EscrowDisputeEvidence represents evidence submitted for disputes
type EscrowDisputeEvidence struct {
	ID           string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	DisputeID    string    `json:"dispute_id" gorm:"not null;index"`
	SubmittedBy  string    `json:"submitted_by" gorm:"not null"`
	EvidenceType string    `json:"evidence_type" gorm:"not null"` // document, screenshot, video, audio, text, link
	Title        string    `json:"title" gorm:"not null"`
	Description  string    `json:"description" gorm:"type:text"`
	FileURL      *string   `json:"file_url"`
	FileHash     *string   `json:"file_hash"`
	FileSize     *int64    `json:"file_size"`
	MimeType     *string   `json:"mime_type"`
	TextContent  *string   `json:"text_content" gorm:"type:text"`
	Verified     bool      `json:"verified" gorm:"default:false"`
	VerifiedBy   *string   `json:"verified_by"`
	VerifiedAt   *time.Time `json:"verified_at"`
	Weight       float64   `json:"weight" gorm:"default:1.0"` // Evidence weight in resolution
	Metadata     JSONField `json:"metadata" gorm:"type:jsonb"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Dispute      *EscrowDispute `json:"dispute,omitempty" gorm:"foreignKey:DisputeID"`
}

// EscrowDisputeResolution tracks resolution attempts and decisions
type EscrowDisputeResolution struct {
	ID                string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	DisputeID         string    `json:"dispute_id" gorm:"not null;index"`
	ResolutionAttempt int       `json:"resolution_attempt" gorm:"not null;default:1"`
	ResolverID        string    `json:"resolver_id" gorm:"not null"`
	ResolverType      string    `json:"resolver_type" gorm:"not null"` // admin, ai_system, arbitrator, panel
	Decision          string    `json:"decision" gorm:"not null"`
	Reasoning         string    `json:"reasoning" gorm:"type:text;not null"`
	ConfidenceScore   float64   `json:"confidence_score" gorm:"default:0"`
	EvidenceReviewed  JSONField `json:"evidence_reviewed" gorm:"type:jsonb"`
	ActionsTaken      JSONField `json:"actions_taken" gorm:"type:jsonb"`
	FinancialImpact   JSONField `json:"financial_impact" gorm:"type:jsonb"`
	Appealed          bool      `json:"appealed" gorm:"default:false"`
	AppealDeadline    *time.Time `json:"appeal_deadline"`
	Implemented       bool      `json:"implemented" gorm:"default:false"`
	ImplementedAt     *time.Time `json:"implemented_at"`
	ImplementedBy     *string   `json:"implemented_by"`
	Notes             string    `json:"notes" gorm:"type:text"`
	Metadata          JSONField `json:"metadata" gorm:"type:jsonb"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Dispute           *EscrowDispute `json:"dispute,omitempty" gorm:"foreignKey:DisputeID"`
}

// EscrowPayment tracks individual payments within an escrow transaction
type EscrowPayment struct {
	ID                  string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	EscrowTransactionID string    `json:"escrow_transaction_id" gorm:"not null;index"`
	PaymentType         string    `json:"payment_type" gorm:"not null"` // deposit, release, refund, fee, penalty, bonus
	RecipientID         string    `json:"recipient_id" gorm:"not null"`
	RecipientType       string    `json:"recipient_type" gorm:"not null"` // user, verifier, platform, treasury
	Amount              float64   `json:"amount" gorm:"not null"`
	Currency            string    `json:"currency" gorm:"not null;default:'HBAR'"`
	Status              string    `json:"status" gorm:"not null;default:'pending'"` // pending, processing, completed, failed, cancelled
	HederaTxID          *string   `json:"hedera_tx_id"`
	HederaTxHash        *string   `json:"hedera_tx_hash"`
	ProcessingStarted   *time.Time `json:"processing_started"`
	ProcessingCompleted *time.Time `json:"processing_completed"`
	FailureReason       *string   `json:"failure_reason"`
	RetryCount          int       `json:"retry_count" gorm:"default:0"`
	MaxRetries          int       `json:"max_retries" gorm:"default:3"`
	NextRetryAt         *time.Time `json:"next_retry_at"`
	Description         string    `json:"description"`
	Reference           *string   `json:"reference"`
	Metadata            JSONField `json:"metadata" gorm:"type:jsonb"`
	CreatedAt           time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	EscrowTransaction   *EscrowTransaction `json:"escrow_transaction,omitempty" gorm:"foreignKey:EscrowTransactionID"`
}

// EscrowConfiguration stores system-wide escrow settings
type EscrowConfiguration struct {
	ID                      string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ConfigKey               string    `json:"config_key" gorm:"not null;uniqueIndex"`
	ConfigValue             string    `json:"config_value" gorm:"not null"`
	ConfigType              string    `json:"config_type" gorm:"not null;default:'string'"` // string, number, boolean, json
	Description             string    `json:"description" gorm:"type:text"`
	Category                string    `json:"category" gorm:"not null;default:'general'"`
	IsActive                bool      `json:"is_active" gorm:"default:true"`
	RequiresRestart         bool      `json:"requires_restart" gorm:"default:false"`
	ValidationRules         JSONField `json:"validation_rules" gorm:"type:jsonb"`
	DefaultValue            *string   `json:"default_value"`
	MinValue                *float64  `json:"min_value"`
	MaxValue                *float64  `json:"max_value"`
	AllowedValues           pq.StringArray `json:"allowed_values" gorm:"type:text[]"`
	EnvironmentSpecific     bool      `json:"environment_specific" gorm:"default:false"`
	Sensitive               bool      `json:"sensitive" gorm:"default:false"`
	LastModifiedBy          string    `json:"last_modified_by" gorm:"not null"`
	ModificationReason      *string   `json:"modification_reason"`
	EffectiveFrom           time.Time `json:"effective_from" gorm:"not null;default:now()"`
	EffectiveUntil          *time.Time `json:"effective_until"`
	Metadata                JSONField `json:"metadata" gorm:"type:jsonb"`
	CreatedAt               time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt               time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// BeforeCreate hook for EscrowTransaction
func (et *EscrowTransaction) BeforeCreate(tx *gorm.DB) error {
	// Set default values
	if et.Status == "" {
		et.Status = "pending"
	}
	if et.EscrowType == "" {
		et.EscrowType = "verification"
	}
	if et.Currency == "" {
		et.Currency = "HBAR"
	}
	if et.FraudCheckStatus == "" {
		et.FraudCheckStatus = "pending"
	}

	// Calculate total amount
	et.TotalAmount = et.Amount + et.FeeAmount

	// Set auto-release deadline if enabled
	if et.AutoReleaseEnabled && et.AutoReleaseDelay > 0 {
		autoReleaseAt := time.Now().Add(time.Duration(et.AutoReleaseDelay) * time.Hour)
		et.AutoReleaseAt = &autoReleaseAt
	}

	return nil
}

// BeforeUpdate hook for EscrowTransaction
func (et *EscrowTransaction) BeforeUpdate(tx *gorm.DB) error {
	// Update total amount if amount or fee changed
	et.TotalAmount = et.Amount + et.FeeAmount

	// Set timestamps based on status changes
	switch et.Status {
	case "locked":
		if et.LockedAt == nil {
			now := time.Now()
			et.LockedAt = &now
		}
	case "released":
		if et.ReleasedAt == nil {
			now := time.Now()
			et.ReleasedAt = &now
		}
	case "refunded":
		if et.RefundedAt == nil {
			now := time.Now()
			et.RefundedAt = &now
		}
	case "cancelled":
		if et.CancelledAt == nil {
			now := time.Now()
			et.CancelledAt = &now
		}
	}

	return nil
}

// TableName returns the table name for EscrowTransaction
func (EscrowTransaction) TableName() string {
	return "escrow_transactions"
}

// TableName returns the table name for EscrowStatusHistory
func (EscrowStatusHistory) TableName() string {
	return "escrow_status_history"
}

// TableName returns the table name for EscrowDispute
func (EscrowDispute) TableName() string {
	return "escrow_disputes"
}

// TableName returns the table name for EscrowDisputeEvidence
func (EscrowDisputeEvidence) TableName() string {
	return "escrow_dispute_evidence"
}

// TableName returns the table name for EscrowDisputeResolution
func (EscrowDisputeResolution) TableName() string {
	return "escrow_dispute_resolutions"
}

// TableName returns the table name for EscrowPayment
func (EscrowPayment) TableName() string {
	return "escrow_payments"
}

// TableName returns the table name for EscrowConfiguration
func (EscrowConfiguration) TableName() string {
	return "escrow_configurations"
}