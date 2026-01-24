package user

import (
	"context"
	"errors"

	app "reading-cats-api/internal/application/user"
	domain "reading-cats-api/internal/domain/user"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) FindByCognitoSub(ctx context.Context, sub domain.CognitoSub) (*domain.User, error) {
	q := `
SELECT id, cognito_sub, COALESCE(email,''), COALESCE(display_name,''), COALESCE(avatar_url,''), profile_source, created_at, updated_at
FROM users
WHERE cognito_sub = $1
LIMIT 1;
`
	var u domain.User
	var cognitoSub string
	var email, name, avatar string
	var profileSource string

	err := r.pool.QueryRow(ctx, q, string(sub)).
		Scan(&u.ID, &cognitoSub, &email, &name, &avatar, &profileSource, &u.CreatedAt, &u.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	cs, err := domain.NewCognitoSub(cognitoSub)
	if err != nil {
		return nil, err
	}
	em, err := domain.NewEmail(email)
	if err != nil {
		return nil, err
	}
	dn, err := domain.NewDisplayName(name)
	if err != nil {
		return nil, err
	}
	av, err := domain.NewAvatarURL(avatar)
	if err != nil {
		return nil, err
	}

	u.CognitoSub = cs
	u.Email = em
	u.DisplayName = dn
	u.AvatarURL = av
	u.ProfileSource = domain.ProfileSource(profileSource)

	return &u, nil
}

func (r *PostgresRepository) Insert(ctx context.Context, u *domain.User) error {
	q := `
INSERT INTO users (id, cognito_sub, email, display_name, avatar_url, profile_source, created_at, updated_at)
VALUES ($1, $2, NULLIF($3,''), NULLIF($4,''), NULLIF($5,''), $6, now(), now());
`
	_, err := r.pool.Exec(ctx, q,
		u.ID,
		string(u.CognitoSub),
		string(u.Email),
		string(u.DisplayName),
		string(u.AvatarURL),
		string(u.ProfileSource),
	)
	// Se bater race (2 requests simultâneas), você pode tratar conflito depois com retry:
	_ = app.Repository(nil)
	return err
}
