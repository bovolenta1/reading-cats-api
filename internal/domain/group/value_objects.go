package group

import "strings"

type GroupName string

func NewGroupName(v string) (GroupName, error) {
	v = strings.TrimSpace(v)
	if len(v) < 1 || len(v) > 30 {
		return "", ErrInvalidGroupName
	}
	return GroupName(v), nil
}

type IconID string

func NewIconID(v string) (IconID, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return "", ErrInvalidIconID
	}
	return IconID(v), nil
}

type Visibility string

const (
	VisibilityInviteOnly Visibility = "INVITE_ONLY"
	VisibilityPublic     Visibility = "PUBLIC_SOON"
	VisibilityFounders   Visibility = "FOUNDERS"
)

func (v Visibility) String() string {
	return string(v)
}
