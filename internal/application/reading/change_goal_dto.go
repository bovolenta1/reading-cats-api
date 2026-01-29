package reading

import (
	readingDomain "reading-cats-api/internal/domain/reading"
	userDomain "reading-cats-api/internal/domain/user"
)

type ChangeGoalInput struct {
	Claims userDomain.IDPClaims
	Pages  readingDomain.Pages
}

type ChangeGoalOutput struct {
	CurrentGoal *GoalRecord `json:"current_goal"`
	NextGoal    *GoalRecord `json:"next_goal"`
}

type GoalRecord struct {
	DailyPages int    `json:"daily_pages"`
	ValidFrom  string `json:"valid_from"`
}
