package util

import "regexp"

func IndetifyPhone(phone string) bool {
	var phoneRegex = regexp.MustCompile(`^1(3[0-9]|4[57]|5[0-35-9]|7[0-9]|8[0-9]|9[8])\d{8}$`)
	return phoneRegex.MatchString(phone)
}
