package util

func Trim(s string) string {
	if s == "" || len(s) < 2 {
		return s
	}

	return s[1 : len(s)-1]
}
