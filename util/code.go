package util

import "crypto/rand"

var (
	potentialChars = "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	length         = 6
)

func GenerateRandomCode() (string, error) {

	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(potentialChars)
	for i := 0; i < length; i++ {
		buffer[i] = potentialChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}
