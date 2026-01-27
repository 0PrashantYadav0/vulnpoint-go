package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScanResult struct {
	ID           uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	WorkflowID   *uuid.UUID      `gorm:"type:uuid" json:"workflow_id,omitempty"`
	UserID       uuid.UUID       `gorm:"type:uuid;not null" json:"user_id"`
	ScanType     string          `gorm:"not null" json:"scan_type"`
	TargetURL    string          `gorm:"not null" json:"target_url"`
	Status       string          `gorm:"default:'pending'" json:"status"`
	Results      json.RawMessage `gorm:"type:jsonb" json:"results,omitempty"`
	ErrorMessage string          `json:"error_message,omitempty"`
	StartedAt    *time.Time      `json:"started_at,omitempty"`
	CompletedAt  *time.Time      `json:"completed_at,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
}

func (ScanResult) TableName() string {
	return "scan_results"
}

func (s *ScanResult) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}