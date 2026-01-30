package season

import "strings"

type Status string

const (
	StatusDraft  Status = "DRAFT"
	StatusActive Status = "ACTIVE"
	StatusEnded  Status = "ENDED"
)

func (s Status) String() string {
	return string(s)
}

type Metric string

const (
	MetricCheckinsPerDay Metric = "CHECKINS_PER_DAY"
	MetricPagesSoon      Metric = "PAGES_SOON"
	MetricMinutesSoon    Metric = "MINUTES_SOON"
)

func (m Metric) String() string {
	return string(m)
}

type Timezone string

func NewTimezone(v string) (Timezone, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return "", ErrInvalidTimezone
	}
	return Timezone(v), nil
}
