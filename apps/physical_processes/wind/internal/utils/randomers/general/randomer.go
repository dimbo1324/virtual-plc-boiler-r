package general

import (
	"crypto/rand"
	"math/big"
)

// TODO: docs and discription
func randomInt(min, max int64) (int64, error) {
	rangeSize := big.NewInt(max - min + 1)
	n, err := rand.Int(rand.Reader, rangeSize)
	if err != nil {
		return 0, err
	}
	return n.Int64() + min, nil
}
