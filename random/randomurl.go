package random

import (
	"crypto/rand"
	"io"
	"math"
	"math/big"
	"os"
)

var (
	randomFile = "/dev/urandom"
)

type Random struct {
	randomReader io.ReadCloser
}

func NewRandom() (*Random, error) {
	file, err := os.Open(randomFile)
	if err != nil {
		return nil, err
	}

	return &Random{
		randomReader: file,
	}, nil
}

func (r *Random) GetRandomUrl(length int) (string, error) {
	// I don't care about accuracy.
	max, _ := big.NewFloat(math.Pow(10.0, float64(length))).Int(nil)
	bigInt, err := rand.Int(r.randomReader, max)
	if err != nil {
		return "", err
	}
	return bigInt.String(), nil
}
