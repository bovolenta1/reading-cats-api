package reading

import (
	readingDomain "reading-cats-api/internal/domain/reading"
	userDomain "reading-cats-api/internal/domain/user"
)

type RegisterReadingInput struct {
	Claims userDomain.IDPClaims
	Pages  readingDomain.Pages
}

type RegisterReadingOutput struct {
	Progress readingDomain.ReadingProgress
}

type GetReadingProgressInput struct {
	Claims userDomain.IDPClaims
}

type GetReadingProgressOutput struct {
	Progress    readingDomain.ReadingProgress `json:"progress"`
	CurrentGoal *GoalRecord                   `json:"current_goal"`
	NextGoal    *GoalRecord                   `json:"next_goal,omitempty"`
}
