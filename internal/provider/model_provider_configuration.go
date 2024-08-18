package provider

import (
	"payment-gateway-service/internal/country"
	"payment-gateway-service/internal/currency"
	"time"
)

// ProviderConfiguration represents a configuration for a payment provider in a specific country and currency.
type ProviderConfiguration struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CountryID    uint      `gorm:"not null" json:"country_id"`
	CurrencyID   uint      `gorm:"not null" json:"currency_id"`
	ProviderID   uint      `gorm:"not null" json:"provider_id"`
	BaseURL      string    `gorm:"not null" json:"base_url"`
	Priority     int       `gorm:"not null;check:priority >= 1" json:"priority"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ProviderName string    `gorm:"column:provider_name" json:"provider_name"`

	// Relationships
	Country  country.Country   `gorm:"foreignKey:CountryID"`
	Currency currency.Currency `gorm:"foreignKey:CurrencyID"`
	Provider Provider          `gorm:"foreignKey:ProviderID"`
}

func (ProviderConfiguration) TableName() string {
	return "provider_configurations"
}
