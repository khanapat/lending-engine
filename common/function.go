package common

import "time"

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
