package provider

import (
	"time"
)

// Provider represents a payment provider in the system.
type Provider struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"unique;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Provider) TableName() string {
	return "payment_providers"
}
