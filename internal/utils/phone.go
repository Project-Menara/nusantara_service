package utils

import (
	"strings"
	"unicode"
)

func IsDigitsOnly(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func NormalizePhone(phone string) string {
	if strings.HasPrefix(phone, "+62") {
		return "0" + phone[3:]
	}
	return phone
}
