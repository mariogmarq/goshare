package util

import (
	"math/rand"
	"strings"
	"time"
)

//Create a random string of a custom size
func CreateRandomString(length int) string {
	rand.Seed(time.Now().Unix())
	charSet := "abcdedfghijklmnopqrstuvwxyz"
	var output strings.Builder
	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}

	return output.String()
}
