package user

type MeDTO struct {
	ID          string `json:"id"`
	CognitoSub  string `json:"cognitoSub"`
	Email       string `json:"email,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	AvatarURL   string `json:"avatarUrl,omitempty"`
	Source      string `json:"profileSource,omitempty"`
}
