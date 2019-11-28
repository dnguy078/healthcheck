package utils

import (
	"crypto/rand"
	"fmt"
	"regexp"
)

var (
	uuidRegex = regexp.MustCompile(`[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`)
)

func UUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

func ContainsUUID(input string) bool {
	return uuidRegex.MatchString(input)
}

func ExtractUUID(url string) string {
	return uuidRegex.FindString(url)
}
