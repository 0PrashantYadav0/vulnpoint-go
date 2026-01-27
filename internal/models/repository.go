package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	GitHubID     int64      `gorm:"uniqueIndex;not null" json:"github_id"`
	FullName     string     `gorm:"not null" json:"full_name"`
	Name         string     `gorm:"not null" json:"name"`
	Description  string     `json:"description"`
	HTMLURL      string     `json:"html_url"`
	Language     string     `json:"language"`
	IsPrivate    bool       `gorm:"default:false" json:"is_private"`
	LastAnalyzed *time.Time `json:"last_analyzed,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (Repository) TableName() string {
	return "repositories"
}

func (r *Repository) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}