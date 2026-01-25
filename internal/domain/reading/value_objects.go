package reading

import (
	"fmt"
	"time"
)

type LocalDate string // "YYYY-MM-DD"
type Pages int
type TargetDatePolicy struct {
	GraceHour int // 2
}
type StreakDays int
type StreakPolicy struct{}

func NewPages(v int) (Pages, error) {
	if v <= 0 || v > 500 {
		return 0, fmt.Errorf("pages out of range")
	}
	return Pages(v), nil
}

func DateOf(t time.Time, loc *time.Location) LocalDate {
	return LocalDate(t.In(loc).Format("2006-01-02"))
}

func (d LocalDate) AddDays(n int) LocalDate {
	tt, _ := time.Parse("2006-01-02", string(d))
	return LocalDate(tt.AddDate(0, 0, n).Format("2006-01-02"))
}

func (d LocalDate) String() string { return string(d) }

func (p TargetDatePolicy) Resolve(now time.Time, loc *time.Location, hasYesterday bool) LocalDate {
	realDate := DateOf(now, loc)
	if now.In(loc).Hour() < p.GraceHour {
		if hasYesterday {
			return realDate
		}
		return realDate.AddDays(-1)
	}
	return realDate
}

func (StreakPolicy) Next(targetDate LocalDate, lastDate LocalDate, lastStreak StreakDays, found bool) StreakDays {
	if !found {
		return 1
	}
	if lastDate == targetDate.AddDays(-1) {
		return lastStreak + 1
	}
	return 1
}
