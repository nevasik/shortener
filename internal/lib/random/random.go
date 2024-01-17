package random

import "math/rand"

func NewRandomString(aliasLength int) string {
	return randomString(aliasLength)
}

func randomString(length int) string {
	bytes := make([]byte, length) // обьявление слайса(массива) с элементами byte и начальным размером массива length

	for i := 0; i < length; i++ {
		bytes[i] = byte(randInt(97, 122))
	}

	return string(bytes) // обернули строку bytes
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
