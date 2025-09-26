package empty

import (
	"context"
	"lytemp/internal/domain/models"
	"lytemp/pkg/database"

	"gorm.io/gorm"
)

type Repository struct {
	db *database.Client
}

func NewRepository(db *database.Client) *Repository {
	return &Repository{db: db}
}

// Example create
func (r *Repository) Create(ctx context.Context, city *models.User, tx ...*gorm.DB) error {
	return nil
	// return r.Get(tx...).WithContext(ctx).Create(city).Error
}

func (r *Repository) Begin() *gorm.DB {
	return r.db.Get().Begin()
}

func (r *Repository) Get(tx ...*gorm.DB) *gorm.DB {
	if len(tx) > 0 && tx[0] != nil {
		return tx[0]
	}

	return r.db.Get()
}
