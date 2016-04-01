package httputil

// EtagIsValid checks whether `etag` has the correct format for a
// (weak or strong) entity tag, as defined in section 2.3 of RFC 7232.
func EtagIsValid(etag string) bool {
	n := len(etag)
	if n < 2 {
		return false
	}
	if etag[0] == 'W' && etag[1] == '/' {
		etag = etag[2:]
		n -= 2
		if n < 2 {
			return false
		}
	}
	if etag[0] != '"' || etag[n-1] != '"' {
		return false
	}
	for i := 1; i < n-1; i++ {
		c := etag[i]
		if c < '\x21' || c == '"' || c == '\x7f' {
			return false
		}
	}
	return true
}

// EtagsEqualStrong checks whether two etags are equivalent under
// strong comparison, as defined in section 2.3.2 of RFC 7232.
func EtagsEqualStrong(a, b string) bool {
	return len(a) > 0 && a[0] == '"' && a == b
}

// EtagsEqualWeak checks whether two etags are equivalent under weak
// comparison, as defined in section 2.3.2 of RFC 7232.
func EtagsEqualWeak(a, b string) bool {
	if len(a) > 2 && a[0] == 'W' && a[1] == '/' {
		a = a[2:]
	}
	if len(b) > 2 && b[0] == 'W' && b[1] == '/' {
		b = b[2:]
	}
	return a == b
}

// EtagsSplit splits a list of comma-separated entity tags.
func EtagsSplit(list string) []string {
	var res []string
	i := 0
	n := len(list)
	for i < n {
		for ; i < n; i++ {
			if list[i] != ' ' && list[i] != '\t' {
				break
			}
		}
		if i == n {
			break
		}
		start := i

		quoteCount := 0
		for ; i < n; i++ {
			if list[i] == '"' {
				quoteCount++
			}
			if quoteCount == 2 {
				i++
				break
			}
		}
		end := i

		res = append(res, list[start:end])

		for ; i < n; i++ {
			if list[i] != ' ' && list[i] != '\t' {
				break
			}
		}
		if i < n && list[i] == ',' {
			i++
		}
	}
	return res
}
