package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
	"verza/backend/pkg/database"
)

// JSONField is a type alias for GORM jsonb support
type JSONField map[string]interface{}

// StringArrayJSON is a type alias for pq.StringArray to support JSON marshaling
type StringArrayJSON pq.StringArray

// User is an alias to the database.User struct for consistency
type User = database.User

// VerificationRequest represents a verification request in the system
type VerificationRequest struct {
	ID          string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	VerifierID  *string    `json:"verifier_id" gorm:"index"`
	UserID      string     `json:"user_id" gorm:"not null;index"`
	UserEmail   string     `json:"user_email" gorm:"not null"`
	Status      string     `json:"status" gorm:"not null;default:'pending'"` // pending, accepted, completed, rejected, cancelled
	Rating      *float64   `json:"rating"`
	AcceptedAt  *time.Time `json:"accepted_at"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	User     *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Verifier *Verifier `json:"verifier,omitempty" gorm:"foreignKey:VerifierID"`
}

// EmailService is a placeholder struct for email service dependencies
type EmailService struct {
	// Add fields as needed when implementing email functionality
}

// SendTemplateEmail is a placeholder method for email service
func (es *EmailService) SendTemplateEmail(email, template, subject string, data map[string]interface{}) error {
	// TODO: Implement email sending functionality
	return nil
}

// TableName returns the table name for VerificationRequest
func (VerificationRequest) TableName() string {
	return "verification_requests"
}