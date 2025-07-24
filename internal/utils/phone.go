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
	phone = strings.TrimSpace(phone)
	if strings.HasPrefix(phone, "0") {
		return "+62" + phone[1:]
	}
	if strings.HasPrefix(phone, "62") {
		return "+" + phone
	}
	return phone // diasumsikan sudah +62
}
