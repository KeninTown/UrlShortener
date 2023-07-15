package random

import (
	"math/rand"
)

const alph = "abcdefghijklmnopqrstuvwxyz0123456789"

func NewRandomString(size int) string {
	var alias string
	for i := 0; i < size; i++ {
		alias += string(alph[rand.Intn(len(alph))])
	}

	return alias
}
