package fetch_provider

import (
	"context"
	"errors"
	"lytemp/internal/domain/models"
	"lytemp/pkg/database"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *database.Client
}

func NewRepository(db *database.Client) *Repository {
	return &Repository{db: db}
}


func (r *Repository) UpsertContent(ctx context.Context, c *models.Content) (bool, error) {
	db := r.db.Get().WithContext(ctx)

	var existing models.Content
	err := db.
		Where("provider = ? AND provider_item_id = ?", c.Provider, c.ProviderItemID).
		First(&existing).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		if err := db.Create(c).Error; err != nil {
			return false, err
		}
		return true, nil

	case err != nil:
		return false, err

	default:
		updates := map[string]interface{}{
			"title":                c.Title,
			"normalized_type":      c.NormalizedType,
			"views":                c.Views,
			"likes":                c.Likes,
			"duration_seconds":     c.DurationSeconds,
			"reading_time_minutes": c.ReadingTimeMinutes,
			"reactions":            c.Reactions,
			"comments":             c.Comments,
			"published_at":         c.PublishedAt,
			"tags":                 c.Tags,
			"base_score":           c.BaseScore,
			"freshness_score":      c.FreshnessScore,
			"engagement_score":     c.EngagementScore,
			"final_score":          c.FinalScore,
			"raw_payload":          c.RawPayload,
			"indexed_at":           c.IndexedAt,
		}

		if err := db.Model(&models.Content{}).
			Where("id = ?", existing.ID).
			Updates(updates).Error; err != nil {
			return false, err
		}
		return false, nil
	}
}

func (r *Repository) UpsertProviderSync(
	ctx context.Context,
	name string,
	format models.ProviderFormat,
	rateLimit int,
	status models.ProviderStatus,
	lastSync time.Time,
	lastErr *string,
) error {
	db := r.db.Get().WithContext(ctx)

	var ps models.ProviderSync
	err := db.Where("name = ?", name).First(&ps).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		row := models.ProviderSync{
			Name:               name,
			Format:             format,
			RateLimitPerMinute: rateLimit,
			Status:             status,
			LastSyncAt:         &lastSync,
			LastError:          lastErr,
		}
		return db.Create(&row).Error

	case err != nil:
		return err

	default:
		updates := map[string]interface{}{
			"format":                format,
			"rate_limit_per_minute": rateLimit,
			"status":                status,
			"last_sync_at":          lastSync,
			"last_error":            lastErr,
			"updated_at":            gorm.Expr("NOW()"),
		}

		return db.Model(&models.ProviderSync{}).
			Where("id = ?", ps.ID).
			Updates(updates).Error
	}
}
