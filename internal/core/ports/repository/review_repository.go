package portrepository

import (
	"context"

	"github.com/netojso/elephrases-api/internal/core/domain"
)

type ReviewRepository interface {
	Save(ctx context.Context, review *domain.Review) error
}
