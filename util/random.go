package util

import (
	"math/rand"
	"strings"
	"time"
)

const letters = "SqudgyFezBlankJimpCrwthVox"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(letters)

	for i := 0; i < n; i++ {
		c := letters[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomHolder() string {
	return RandomString(10)
}

func RandomAmount() int64 {
	return RandomInt(0, 5000)
}

func RandomSymbol() string {
	symbols := []string{"BTC", "ETH", "ADA"}
	n := len(symbols)
	return symbols[rand.Intn(n)]
}
