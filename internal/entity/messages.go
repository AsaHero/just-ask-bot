package entity

import "time"

type Messages struct {
	ID          int64 `gorm:"type=uuid;primaryKey"`
	ExternalID  int64 `gorm:"uniqueIndex"`
	ChatID      string
	UserID      string
	Role        string
	Row         string
	Messages    string
	MexTokens   int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeliveredAt time.Time
}
