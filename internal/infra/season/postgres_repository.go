package season

import (
	"context"
	"fmt"

	domainSeason "reading-cats-api/internal/domain/season"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Insert(ctx context.Context, s *domainSeason.Season) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO group_seasons (id, group_id, status, started_at, ends_at, timezone, metric, created_by_user_id, created_at, updated_at)
		 VALUES ($1, $2, $3::group_season_status, $4, $5, $6, $7::group_metric, $8, $9, $10)`,
		s.ID,
		s.GroupID,
		s.Status.String(),
		s.StartedAt,
		s.EndsAt,
		string(s.Timezone),
		s.Metric.String(),
		s.CreatedByUserID,
		s.CreatedAt,
		s.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert season: %w", err)
	}
	return nil
}
