package httputil

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
		if c < '\x21' || c == '"' || c > '\x7e' {
			return false
		}
	}
	return true
}

func EtagsEqualStrong(a, b string) bool {
	return a[0] == '"' && a == b
}

func EtagsEqualWeak(a, b string) bool {
	if len(a) > 2 && a[0] == 'W' && a[1] == '/' {
		a = a[2:]
	}
	if len(b) > 2 && b[0] == 'W' && b[1] == '/' {
		b = b[2:]
	}
	return a == b
}
