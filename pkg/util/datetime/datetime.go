package datetime

import "time"

// FormatDate returns time formatted as RFC3339 (2020-11-10T13:28:12+09:00)
func FormatDate(t time.Time) string {
	return t.Format(time.RFC3339)
}
