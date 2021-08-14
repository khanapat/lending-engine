package common

import (
	"regexp"
	"time"
)

var EmailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func RandStringBytes(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = RandomStringInt[SeededRand.Intn(len(RandomStringInt))]
	}
	return string(b)
}

func RandIntBytes(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = RandomInt[SeededRand.Intn(len(RandomInt))]
	}
	return string(b)
}

func TimeDateSecondLeft() int {
	timeNow := time.Now()
	y, m, d := timeNow.Date()
	return int(time.Date(y, m, d, 0, 0, 0, 0, time.Local).AddDate(0, 0, 1).Sub(timeNow).Seconds())
}
