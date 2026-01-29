package user

import "errors"

var (
	ErrInvalidCognitoSub = errors.New("invalid cognito sub")
	ErrInvalidEmail      = errors.New("invalid email")
	ErrInvalidName       = errors.New("invalid display name")
	ErrInvalidAvatarURL  = errors.New("invalid avatar url")
	ErrUserNotFound      = errors.New("user not found")
)
