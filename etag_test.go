package httputil

import (
	"testing"
)

func TestIsValidEtag(t *testing.T) {
	goodCases := []string{
		`"xyzzy"`,
		`W/"xyzzy"`,
		`""`,
		`"\"`,
	}
	for _, etag := range goodCases {
		if !EtagIsValid(etag) {
			t.Errorf("failed to accept valid etag %s", etag)
		}
	}

	badCases := []string{
		``,
		`a`,
		`W\`,
		`"""`,
		`w/"xyzzy"`,
		"\"abc\ndef\"",
	}
	for _, etag := range badCases {
		if EtagIsValid(etag) {
			t.Errorf("wrongly accepted invalid etag %s", etag)
		}
	}
}

func TestEtagsEqual(t *testing.T) {
	type testCase struct {
		a, b                   string
		strongEqual, weakEqual bool
	}
	table := []testCase{
		{`"a"`, `"b"`, false, false},
		{`W/"1"`, `W/"1"`, false, true},
		{`W/"1"`, `W/"2"`, false, false},
		{`W/"1"`, `"1"`, false, true},
		{`"1"`, `"1"`, true, true},
	}
	for _, test := range table {
		strongEqual := EtagsEqualStrong(test.a, test.b)
		if strongEqual != test.strongEqual {
			t.Errorf("strong comp. of %s and %s failed: expected %t, got %t",
				test.a, test.b, test.strongEqual, strongEqual)
		}
		weakEqual := EtagsEqualWeak(test.a, test.b)
		if weakEqual != test.weakEqual {
			t.Errorf("weak comp. of %s and %s failed: expected %t, got %t",
				test.a, test.b, test.weakEqual, weakEqual)
		}
	}
}
