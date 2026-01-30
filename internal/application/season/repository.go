package season

import (
	"context"
	domainSeason "reading-cats-api/internal/domain/season"
)

type Repository interface {
	Insert(ctx context.Context, s *domainSeason.Season) error
}
