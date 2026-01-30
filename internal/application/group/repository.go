package group

import (
	"context"
	domainGroup "reading-cats-api/internal/domain/group"
)

type Repository interface {
	Insert(ctx context.Context, g *domainGroup.Group) error
	AddMember(ctx context.Context, groupID string, userID string, role string) error
}
