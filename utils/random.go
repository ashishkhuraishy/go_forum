package utils

import (
	"math/rand"
	"strings"
	"time"
)

var (
	alphabets = "abcdefghijklmnopqrstuvwxyz"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomString takes in a length and returns a
// random string based on the length
func RandomString(length int) string {
	var res strings.Builder
	k := len(alphabets)

	for i := 0; i < length; i++ {
		letter := alphabets[rand.Intn(k)]
		res.WriteByte(letter)
	}

	return res.String()
}
