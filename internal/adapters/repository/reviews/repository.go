package repository

import (
	"context"

	"github.com/netojso/elephrases-api/internal/adapters/repository/gormx"
	"github.com/netojso/elephrases-api/internal/core/domain"
	portrepository "github.com/netojso/elephrases-api/internal/core/ports/repository"
	"gorm.io/gorm"
)

type reviewRepository struct {
	DB *gorm.DB
}

func NewReviewRepository(db *gorm.DB) portrepository.ReviewRepository {
	return &reviewRepository{DB: db}
}

func (r *reviewRepository) Save(ctx context.Context, review *domain.Review) error {
	return gormx.DB(ctx, r.DB).Create(domainToModel(review)).Error
}
