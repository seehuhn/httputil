// headers_test.go - unit tests for headers.go
// Copyright (C) 2014  Jochen Voss <voss@seehuhn.de>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package httputil

import (
	"reflect"
	"testing"
)

func TestNormalizeHeader(t *testing.T) {
	type testCase struct {
		in, out string
	}
	table := []testCase{
		{"", ""},
		{" ", ""},
		{"a", "a"},
		{"*", "*"},
		{"a  b", "a b"},
		{"a,b", "a,b"},
		{"a, b", "a,b"},
		{"a  ,  b", "a,b"},
		{"a=b", "a=b"},
		{"a= b", "a=b"},
		{"a  =  b", "a=b"},
		{`a = "b  c"`, `a="b  c"`},
	}
	for _, test := range table {
		res := NormalizeHeader(test.in)
		if res != test.out {
			t.Errorf("normalizing %q: expected %q, got %q",
				test.in, test.out, res)
		}
	}
}

func TestTokenizeHeader(t *testing.T) {
	type testCase struct {
		in  string
		out []string
		err error
	}

	table := []testCase{
		{"", []string{}, nil},
		{" ", []string{}, nil},
		{"a", []string{"a"}, nil},
		{"aaaa", []string{"aaaa"}, nil},
		{"a bb ccc", []string{"a", "bb", "ccc"}, nil},
		{"\"\\", nil, ErrUnterminatedEscape},
		{"\"", nil, ErrUnterminatedString},
		{"\r", nil, ErrUnexpectedControlCharacter},
		{"\t", []string{}, nil},
		{" \t a  \t\t  ", []string{"a"}, nil},
		{`"a  \" b" c`, []string{`"a  \" b"`, "c"}, nil},
		{"<>", []string{"<", ">"}, nil},
	}

	for _, test := range table {
		res, err := tokenizeHeader(test.in)
		if err != test.err {
			t.Errorf("tokenizing %q returned error '%v' instead of '%v'",
				test.in, err, test.err)
		} else if !reflect.DeepEqual(res, test.out) {
			t.Errorf("tokenizing %q returned result '%#v' instead of '%#v'",
				test.in, res, test.out)
		}

		in2 := NormalizeHeader(test.in)
		if len(in2) > len(test.in) {
			t.Errorf("%q normalized to longer %q", test.in, in2)
		}
		res, err = tokenizeHeader(in2)
		if err != test.err {
			t.Errorf("tokenizing %q returned error '%v' instead of '%v'",
				in2, err, test.err)
		} else if !reflect.DeepEqual(res, test.out) {
			t.Errorf("tokenizing %q returned result '%#v' instead of '%#v'",
				in2, res, test.out)
		}
	}
}

func TestHeaderParts(t *testing.T) {
	in := HeaderParts{}
	out := in.String()
	expected := ""
	if out != expected {
		t.Errorf("wrong formatting, expected %q, got %q", expected, out)
	}

	in = HeaderParts{
		{"a", ""},
		{"b", "c"},
	}
	out = in.String()
	expected = "a, b=c"
	if out != expected {
		t.Errorf("wrong formatting, expected %q, got %q", expected, out)
	}
}

func TestParseHeader(t *testing.T) {
	type parseTest struct {
		in  string
		out HeaderParts
		err error
	}

	table := []parseTest{
		{"", HeaderParts{}, nil},
		{"\"", nil, ErrUnterminatedString},
		{"\"hello\"", nil, ErrUnexpectedQuotedString},
		{"a", HeaderParts{{"a", ""}}, nil},
		{"a=1", HeaderParts{{"a", "1"}}, nil},
		{"a=\"1\"", HeaderParts{{"a", "\"1\""}}, nil},
		{" ,,,, , ,, , , ", HeaderParts{}, nil},
		{"a,b", HeaderParts{{"a", ""}, {"b", ""}}, nil},
		{"a=1,b,c=2", HeaderParts{{"a", "1"}, {"b", ""}, {"c", "2"}}, nil},
		{"a,b c,d", nil, ErrMissingComma},
		{"a,b=,d", nil, ErrMissingComma},
		{"a,b=c d,e", nil, ErrMissingComma},
	}

	for _, test := range table {
		res, err := ParseHeader(test.in)
		if err != test.err {
			t.Errorf("parsing %q returned error '%v' instead of '%v'",
				test.in, err, test.err)
		} else if !reflect.DeepEqual(res, test.out) {
			t.Errorf("parsing %q returned result '%#v' instead of '%#v'",
				test.in, res, test.out)
		}
	}
}
