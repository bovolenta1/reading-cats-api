package group

import (
	"context"
	"fmt"

	domainGroup "reading-cats-api/internal/domain/group"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Insert(ctx context.Context, g *domainGroup.Group) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO groups (id, name, icon_id, visibility, max_members, created_by_user_id, created_at, updated_at)
		 VALUES ($1, $2, $3, $4::group_visibility, $5, $6, $7, $8)`,
		g.ID,
		string(g.Name),
		string(g.IconID),
		g.Visibility.String(),
		g.MaxMembers,
		g.CreatedByUserID,
		g.CreatedAt,
		g.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert group: %w", err)
	}
	return nil
}

func (r *PostgresRepository) AddMember(ctx context.Context, groupID string, userID string, role string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO group_members (group_id, user_id, role, is_active)
		 VALUES ($1, $2, $3::group_member_role, true)`,
		groupID,
		userID,
		role,
	)
	if err != nil {
		return fmt.Errorf("failed to add group member: %w", err)
	}
	return nil
}
