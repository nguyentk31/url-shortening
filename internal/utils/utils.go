package utils

import (
	"fmt"
	"regexp"
	"strings"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func FormatUrl(url *string) error {
	// If scheme missing, add http://
	if !strings.HasPrefix(*url, "http://") && !strings.HasPrefix(*url, "https://") {
		*url = "http://" + *url
	}

	re := regexp.MustCompile(`^https?:\/\/(?:www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b(?:[-a-zA-Z0-9()@:%_\+.~#?&\/=]*)$`)
	if !re.MatchString(*url) {
		return fmt.Errorf("invalid URL: %s", *url)
	}

	return nil
}

func ConvertBase10ToBase62(base10 int64) string {
	base62 := ""
	for base10 > 0 {
		remainder := base10 % 62
		base62 = string(base62Chars[remainder]) + base62
		base10 /= 62
	}
	return base62
}

// func ConvertBase62ToBase10(base62 string) (int64, error) {
// 	var base10 int64
// 	base10l := len(base62)
// 	for i := 0; i < base10l; i++ {
// 		char := base62[i]
// 		index := strings.IndexByte(base62Chars, char)
// 		if index == -1 {
// 			return 0, fmt.Errorf("invalid base62 character: %c", char)
// 		}
// 		base10 = base10*62 + int64(index)
// 	}
// 	return base10, nil
// }
