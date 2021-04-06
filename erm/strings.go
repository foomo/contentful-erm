package erm

func sliceIncludes(slice []string, key string) bool {
	for _, val := range slice {
		if val == key {
			return true
		}
	}
	return false
}

// latinLettersOnly returns a string with non-latin-letter characters removed
func latinLettersOnly(s string) string {
	i, b := 0, make([]byte, len(s))
	for _, c := range s {
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
			// Is a letter a-z or A-Z
			b[i] = byte(c)
			i++
		}
	}
	return string(b[:i])
}
