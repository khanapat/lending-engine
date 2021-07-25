package common

func RandStringBytes(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = RandomStringInt[SeededRand.Intn(len(RandomStringInt))]
	}
	return string(b)
}
