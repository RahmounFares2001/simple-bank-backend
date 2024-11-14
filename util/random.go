package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	// make sure that every time we run, it generate diffrente data
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1) // random nb between 0 and max - min
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// random owner
func RandomOwner() string {
	return RandomString(6)
}

// random money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// random currency
func RandomCurrency() string {
	currencies := []string{"USD", "EUR"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

// random email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomOwner())
}
