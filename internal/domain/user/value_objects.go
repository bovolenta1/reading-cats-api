package user

import (
	"net/mail"
	"net/url"
	"strings"
)

type CognitoSub string
type Email string
type DisplayName string
type AvatarURL string

func NewCognitoSub(v string) (CognitoSub, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return "", ErrInvalidCognitoSub
	}
	return CognitoSub(v), nil
}

func NewEmail(v string) (Email, error) {
	v = strings.TrimSpace(strings.ToLower(v))
	if v == "" {
		return "", nil // email opcional, pensando no futuro com otp via sms
	}
	if _, err := mail.ParseAddress(v); err != nil {
		return "", ErrInvalidEmail
	}
	return Email(v), nil
}

func NewDisplayName(v string) (DisplayName, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return "", nil // opcional
	}
	// regra simples p/ MVP (ajusta depois)
	if len([]rune(v)) > 80 {
		return "", ErrInvalidName
	}
	return DisplayName(v), nil
}

func NewAvatarURL(v string) (AvatarURL, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return "", nil // opcional
	}
	u, err := url.Parse(v)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "", ErrInvalidAvatarURL
	}
	return AvatarURL(v), nil
}
