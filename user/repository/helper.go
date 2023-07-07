package repository

import (
	"encoding/base64"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00"
)

// decode cursor from user for mysql
func DecodeCursor(encodedTime string) (time.Time, error) {
	byt, err := base64.StdEncoding.DecodeString(encodedTime)
	if err != nil {
		return time.Time{}, err
	}

	timeString := string(byt)
	t, err := time.Parse(timeFormat, timeString)

	return t, err
}

// encode cursor from mysql
func EncodeCursor(t time.Time) string {
	timeString := t.Format(timeFormat)

	return base64.StdEncoding.EncodeToString([]byte(timeString))
}

func EncodePassword(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
