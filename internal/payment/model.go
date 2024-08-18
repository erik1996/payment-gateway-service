package payment

import (
	"payment-gateway-service/internal/provider"
	"payment-gateway-service/internal/utils"
	"time"
)

type Payment struct {
	ID           string              `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Amount       float64             `gorm:"type:numeric(12,2);not null" json:"amount"`
	PaymentType  utils.PaymentType   `gorm:"type:payment_type;not null" json:"payment_type"`
	Status       utils.PaymentStatus `gorm:"type:payment_status;default:INITIALIZED" json:"status"`
	CurrencyCode string              `gorm:"type:varchar(3);not null" json:"currency_code"`
	UserID       int                 `gorm:"not null" json:"user_id"`
	ProviderID   uint                `gorm:"not null;foreignKey:ProviderID;constraint:OnDelete:SET NULL" json:"provider_id"`
	Provider     provider.Provider   `gorm:"foreignKey:ProviderID" json:"provider"`
	ExternalID   string              `gorm:"type:varchar(255)" json:"external_id"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}
