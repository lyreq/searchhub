package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ContentType string

const (
	ContentTypeVideo   ContentType = "video"
	ContentTypeArticle ContentType = "article"
)

type Content struct {
	gorm.Model
	// Kimlik & provider
	Provider       string `gorm:"type:text;not null;index:idx_provider_item,priority:1" json:"provider"`
	ProviderItemID string `gorm:"type:text;not null;index:idx_provider_item,priority:2" json:"provider_item_id"`

	// Temel bilgiler
	Title          string      `gorm:"type:text;not null;index:,ops:gin" json:"title"`
	NormalizedType ContentType `gorm:"type:text;not null;index:idx_type_score,priority:1" json:"type"`

	// Video metrikleri (opsiyonel)
	Views           *int `gorm:"type:int" json:"views,omitempty"`
	Likes           *int `gorm:"type:int" json:"likes,omitempty"`
	DurationSeconds *int `gorm:"type:int" json:"duration_seconds,omitempty"`

	// Article metrikleri (opsiyonel)
	ReadingTimeMinutes *int `gorm:"type:int" json:"reading_time_minutes,omitempty"`
	Reactions          *int `gorm:"type:int" json:"reactions,omitempty"`
	Comments           *int `gorm:"type:int" json:"comments,omitempty"`

	// Zaman & etiketler
	PublishedAt time.Time      `gorm:"type:timestamptz;not null;index:idx_published_at" json:"published_at"`
	Tags        datatypes.JSON `gorm:"type:jsonb" json:"tags,omitempty"`

	// Skorlar (rapor/sıralama için saklanır)
	BaseScore       float64 `gorm:"type:numeric(8,2);not null" json:"base_score"`
	FreshnessScore  float64 `gorm:"type:numeric(8,2);not null" json:"freshness_score"`
	EngagementScore float64 `gorm:"type:numeric(8,2);not null" json:"engagement_score"`
	FinalScore      float64 `gorm:"type:numeric(8,2);not null;index:idx_final_score,sort:desc;index:idx_type_score,priority:2,sort:desc" json:"final_score"`

	// Ham payload (debug/izlenebilirlik)
	RawPayload datatypes.JSON `gorm:"type:jsonb" json:"raw_payload,omitempty"`

	// İzleme
	IndexedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"indexed_at"`
}
