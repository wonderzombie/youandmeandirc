package youandmeandirc

import "strings"

func has(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func last(ss []string) string {
	return ss[len(ss)-1]
}

func sortaContains(a, b string) bool {
	return strings.Contains(strings.ToLower(a), strings.ToLower(b))
}
