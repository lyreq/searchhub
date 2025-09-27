package contents

import (
	"context"
	"lytemp/internal/domain/models"
	"lytemp/pkg/database"

	"gorm.io/gorm/clause"
)

type Repository struct {
	db *database.Client
}

func NewRepository(db *database.Client) *Repository { return &Repository{db: db} }

func (r *Repository) Search(ctx context.Context, Query string, Type string, Sort string, Page int, PerPage int) ([]models.Content, int64, error) {
	db := r.db.Get().Model(&models.Content{}).WithContext(ctx)

	if Query != "" {
		like := "%" + Query + "%"
		db = db.Where("title ILIKE ? OR CAST(tags AS TEXT) ILIKE ?", like, like)
	}
	if Type != "" {
		db = db.Where("normalized_type = ?", Type)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	switch Sort {
	case "recent":
		db = db.Order("published_at DESC").Order("final_score DESC")
	case "popularity":
		db = db.Order("base_score DESC").Order("final_score DESC")
	case "relevance":
		if Query != "" {
			db = db.
				Order(clause.Expr{
					SQL:  "CASE WHEN title ILIKE ? THEN 1 ELSE 0 END DESC",
					Vars: []interface{}{"%" + Query + "%"},
				}).
				Order("final_score DESC").
				Order("published_at DESC")
		} else {
			db = db.Order("final_score DESC").Order("published_at DESC")
		}
	default: // "score"
		db = db.Order("final_score DESC").Order("published_at DESC")
	}

	page := Page
	per := PerPage
	offset := (page - 1) * per

	var rows []models.Content
	if err := db.Limit(per).Offset(offset).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *Repository) GetByID(ctx context.Context, id uint) (models.Content, error) {
	var row models.Content
	err := r.db.Get().Model(&models.Content{}).
		WithContext(ctx).
		Where("id = ?", id).
		First(&row).Error
	return row, err
}
