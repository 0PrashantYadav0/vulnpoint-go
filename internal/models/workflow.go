package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Workflow struct {
	ID                uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID            uuid.UUID       `gorm:"type:uuid;not null" json:"user_id"`
	Name              string          `gorm:"not null" json:"name"`
	Nodes             JSONArray       `gorm:"type:jsonb;default:'[]'" json:"nodes"`
	Edges             JSONArray       `gorm:"type:jsonb;default:'[]'" json:"edges"`
	IsActive          bool            `gorm:"default:false" json:"is_active"`
	ScheduleFrequency string          `json:"schedule_frequency,omitempty"`
	ScheduleEnabled   bool            `gorm:"default:false" json:"schedule_enabled"`
	NextRun           *time.Time      `json:"next_run,omitempty"`
	LastExecution     json.RawMessage `gorm:"type:jsonb" json:"last_execution,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

// JSONArray custom type for handling JSONB arrays
type JSONArray []interface{}

func (j JSONArray) Value() (driver.Value, error) {
	if j == nil {
		return "[]", nil
	}
	return json.Marshal(j)
}

func (j *JSONArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

func (Workflow) TableName() string {
	return "workflows"
}

func (w *Workflow) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}