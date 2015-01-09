package httputil

import (
	"errors"
	"unicode/utf8"
)

var (
	ErrMissingComma               = errors.New("missing comma")
	ErrUnexpectedControlCharacter = errors.New("unexpected control character")
	ErrUnexpectedQuotedString     = errors.New("unexpected quoted string")
	ErrUnterminatedEscape         = errors.New("unterminated escape")
	ErrUnterminatedString         = errors.New("unterminated string")
)

func isTSpecial(c rune) bool {
	for _, d := range "()<>@,;:\\\"/[]?={} \t" {
		if c == d {
			return true
		}
	}
	return false
}

func isCtl(c rune) bool {
	return c < 32 || c == 127
}

func tokenizeHeader(header string) ([]string, error) {
	res := []string{}

	start := 0
	quoted := false
	escaped := false

	runes := []rune(header)
	for pos, c := range runes {
		if quoted {
			if escaped {
				escaped = false
			} else if c == '\\' {
				escaped = true
			} else if c == '"' {
				quoted = false
				res = append(res, string(runes[start:pos+1]))
				start = pos + 1
			}
		} else if isTSpecial(c) {
			if pos > start {
				res = append(res, string(runes[start:pos]))
			}

			start = pos + 1
			if c == '"' {
				quoted = true
				start = pos
			} else if c != ' ' && c != '\t' {
				res = append(res, string(runes[pos:pos+1]))
			}
		} else if isCtl(c) {
			return nil, ErrUnexpectedControlCharacter
		}
	}

	if escaped {
		return nil, ErrUnterminatedEscape
	} else if quoted {
		return nil, ErrUnterminatedString
	} else if start < len(runes) {
		res = append(res, string(runes[start:]))
	}

	return res, nil
}

// NormalizeHeader converts a HTTP header value into a standard form
// by removing all optional white-space.
func NormalizeHeader(value string) string {
	res := ""
	tokens, err := tokenizeHeader(value)
	if err != nil {
		return value
	}
	for _, token := range tokens {
		t, _ := utf8.DecodeRuneInString(token)
		if t == utf8.RuneError {
			return value
		}
		if res == "" || isTSpecial(t) {
			res += token
		} else {
			t, _ = utf8.DecodeLastRuneInString(res)
			if t == utf8.RuneError {
				return value
			}
			if isTSpecial(t) {
				res += token
			} else {
				res += " " + token
			}
		}
	}
	return res
}

type HeaderPart struct {
	Key, Value string
}

type HeaderParts []HeaderPart

func (parts HeaderParts) String() string {
	if len(parts) == 0 {
		return ""
	}

	n := 2 * (len(parts) - 1) // ", " between entries
	for _, part := range parts {
		n += len(part.Key)
		if part.Value != "" {
			n += 1 + len(part.Value) // value with leading "=" sign
		}
	}
	res := make([]byte, n)

	p := 0
	for i, part := range parts {
		if i > 0 {
			p += copy(res[p:], ", ")
		}
		p += copy(res[p:], part.Key)
		if part.Value != "" {
			p += copy(res[p:], "=")
			p += copy(res[p:], part.Value)
		}
	}
	return string(res)
}

// ParseHeader parses HTTP an header value.  This function only works
// for headers which are defined to be comma-separated lists of tokens
// and key-value pairs, e.g. for the Vary and Cache-Control headers.
func ParseHeader(value string) (HeaderParts, error) {
	tokens, err := tokenizeHeader(value)
	if err != nil {
		return nil, err
	}

	res := []HeaderPart{}

	requireComma := false

	i := 0
	n := len(tokens)
	for i < n {
		if tokens[i] == "," {
			i++
			requireComma = false
			continue
		}
		if requireComma {
			return nil, ErrMissingComma
		}
		if tokens[i][0] == '"' {
			return nil, ErrUnexpectedQuotedString
		}

		part := HeaderPart{Key: tokens[i]}
		i++
		if i+1 < n && tokens[i] == "=" {
			part.Value = tokens[i+1]
			i += 2
		}
		res = append(res, part)
		requireComma = true
	}

	return res, nil
}
