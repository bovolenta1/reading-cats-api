package reading

type ReadingProgress struct {
	Day    DayProgress       `json:"day"`
	Streak StreakProgress    `json:"streak"`
	Week   []WeekDayProgress `json:"week"`
}

type DayProgress struct {
	Date      string `json:"date"`
	Pages     int    `json:"pages"`
	GoalPages int    `json:"goal_pages"`
}

type StreakProgress struct {
	CurrentDays int `json:"current_days"`
}

type WeekDayProgress struct {
	Date    string `json:"date"`
	Pages   int    `json:"pages"`
	Checked bool   `json:"checked"`
}
