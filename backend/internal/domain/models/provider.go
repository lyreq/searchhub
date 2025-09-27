package models

import (
	"time"

	"gorm.io/gorm"
)

type ProviderFormat string
type ProviderStatus string

const (
	ProviderFormatJSON ProviderFormat = "json"
	ProviderFormatXML  ProviderFormat = "xml"

	ProviderStatusOK       ProviderStatus = "ok"
	ProviderStatusDegraded ProviderStatus = "degraded"
	ProviderStatusDown     ProviderStatus = "down"
)

type ProviderSync struct {
	gorm.Model
	Name               string         `gorm:"type:text;uniqueIndex;not null" json:"name"`
	Format             ProviderFormat `gorm:"type:text;not null" json:"format"`
	RateLimitPerMinute int            `gorm:"type:int;not null;default:60" json:"rate_limit_per_minute"`
	LastSyncAt         *time.Time     `gorm:"type:timestamptz" json:"last_sync_at,omitempty"`
	Status             ProviderStatus `gorm:"type:text;not null;default:'ok'" json:"status"`
	LastError          *string        `gorm:"type:text" json:"last_error,omitempty"`
}
