package phonenumber

import (
	"strconv"
)

func IsValid(phoneNumber string) bool {
	//TODO: - tech debt - we can use regex to support +98 phone code
	if len(phoneNumber) != 11 {
		return false
	}

	if phoneNumber[:2] != "09" {
		return false
	}

	if _, err := strconv.Atoi(phoneNumber[2:]); err != nil {
		return false
	}

	return true
}
