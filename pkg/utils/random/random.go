package random

import (
	"math/rand"
	"time"
)

func GenerateSixDigitOTP() int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	otp :=100000 + rng.Intn(900000) // Generates a number between 100000 and 999999
	return otp
}
