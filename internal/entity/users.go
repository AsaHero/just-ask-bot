package entity

import "time"

// User представляет сущность пользователя в системе
type Users struct {
	GUID       string `gorm:"type=uuid;primaryKey;default=uuid_generate_v4()"`
	ExternalID int64  `gorm:"uniqueIndex"`
	Username   string
	FirstName  string
	LastName   string
	IsBlocked  bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

