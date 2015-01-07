package httputil

import (
	"time"
)

// ParseDate converts date strings from HTTP headers (as decribed in
// section 7.1.1.1 or RFC7231) to time.Time objects.  If the date
// string is malformed, a zero date object is returned instead.
func ParseDate(dateStr string) time.Time {
	var t time.Time
	if dateStr != "" {
		for _, format := range []string{time.RFC1123, time.RFC850, time.ANSIC} {
			t, err := time.Parse(format, dateStr)
			if err == nil {
				return t
			}
		}
	}
	return t
}
