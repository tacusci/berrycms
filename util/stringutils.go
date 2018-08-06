package utils

import (
	"regexp"

	"github.com/tacusci/logging"
	"golang.org/x/crypto/bcrypt"
)

type CompiledRegex struct {
	*regexp.Regexp
}

func (cr *CompiledRegex) GetMatchGroupContent(s string, gi int) string {
	result := cr.FindStringSubmatch(s)
	if len(result) >= gi {
		return result[gi]
	}
	return ""
}

func HashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		logging.ErrorAndExit(err.Error())
	}
	return string(hash)
}
