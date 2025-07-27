package utils

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// var digitsOnlyRegex = regexp.MustCompile(`^[0-9]+$`)

func IsDigitsOnly(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func IsDigitsPinOnly(pin string) bool {
	l := len(pin)
	if strings.HasPrefix(pin, "-") {
		l = l - 1
		pin = pin[1:]
	}

	reg := fmt.Sprintf("\\d{%d}", l)

	rs, err := regexp.MatchString(reg, pin)

	if err != nil {
		return false
	}

	return rs
}

func NormalizePhone(phone string) string {
	phone = strings.TrimSpace(phone)
	if strings.HasPrefix(phone, "0") {
		return "+62" + phone[1:]
	}
	if strings.HasPrefix(phone, "62") {
		return "+" + phone
	}
	return phone
}
func IsPhoneDigitsOnly(phone string) bool {
	phone = strings.TrimPrefix(phone, "+")
	return IsDigitsOnly(phone)
}
