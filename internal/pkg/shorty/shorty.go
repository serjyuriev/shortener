package shorty

import (
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyz")

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateShortPath generates a pseudorandom
// letter sequence of fixed length.
func GenerateShortPath() string {
	b := make([]rune, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
