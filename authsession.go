package main

import "regexp"

var (
	tokenPattern = regexp.MustCompile("^([0-9]|[A-Z]|[a-z]){6}$")
)

func isValidToken(token string) bool {
	return tokenPattern.MatchString(token)
}
