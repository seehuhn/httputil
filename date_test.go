package httputil

import (
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {
	dates := []string{
		"Sun, 06 Nov 1994 08:49:37 GMT",
		"Sunday, 06-Nov-94 08:49:37 GMT",
		"Sun Nov  6 08:49:37 1994",
	}
	expected := time.Date(1994, time.November, 6, 8, 49, 37, 0, time.UTC)
	for _, date := range dates {
		res := ParseDate(date)
		if !res.Equal(expected) {
			t.Errorf("wrong time: expected %s, got %s", expected, res)
		}
	}
}
