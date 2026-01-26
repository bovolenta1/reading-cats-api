package reading

import readingDomain "reading-cats-api/internal/domain/reading"

type ChangeGoalInput struct {
	CognitoSub string
	Pages      readingDomain.Pages
}

type ChangeGoalOutput struct {
	CurrentGoal *GoalRecord `json:"current_goal"`
	NextGoal    *GoalRecord `json:"next_goal"`
}

type GoalRecord struct {
	DailyPages int    `json:"daily_pages"`
	ValidFrom  string `json:"valid_from"`
}
