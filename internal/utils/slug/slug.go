package slug

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
)

var (
	reNonAlnum = regexp.MustCompile(`[^a-z0-9]+`)
	reTrimDash = regexp.MustCompile(`(^-+|-+$)`)
)

func FromTitle(title string) string {
	s := strings.ToLower(strings.TrimSpace(title))
	s = reNonAlnum.ReplaceAllString(s, "-")
	s = reTrimDash.ReplaceAllString(s, "")
	if s == "" {
		s = "post"
	}
	if len(s) > 200 {
		s = s[:200]
		s = reTrimDash.ReplaceAllString(s, "")
		if s == "" {
			s = "post"
		}
	}
	return s
}

func RandSuffix(nbytes int) string {
	b := make([]byte, nbytes)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
