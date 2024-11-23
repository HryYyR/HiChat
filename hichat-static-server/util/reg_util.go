package util

import "regexp"

func PhoneValid(phone string) bool {
	reg := `^(13[0-9]|14[01456879]|15[0-3,5-9]|16[2567]|17[0-8]|18[0-9]|19[0-3,5-9])\d{8}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(phone)
}

func EmailValid(email string) bool {
	reg := `^[a-zA-Z0-9-_]+@([a-zA-Z0-9-_]+.)+[a-zA-Z]{2,}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(email)
}
