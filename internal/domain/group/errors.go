package group

import "errors"

var (
	ErrInvalidGroupName = errors.New("invalid group name: must be between 1 and 30 characters")
	ErrInvalidIconID    = errors.New("invalid icon_id: must not be empty")
)
