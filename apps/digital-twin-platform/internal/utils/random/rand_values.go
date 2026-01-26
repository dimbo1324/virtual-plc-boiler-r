package random

import (
	"crypto/rand"
	"errors"
	"math"
	"math/big"
	"strconv"
	"strings"
)

func createRandValue(min, max int64) (int64, error) {
	if min == max {
		return min, nil
	}
	if min > max {
		min, max = max, min
	}
	bgMax := big.NewInt(max)
	bgMin := big.NewInt(min)
	diff := new(big.Int).Sub(bgMax, bgMin)
	plusOne := big.NewInt(1)
	rangeBig := new(big.Int).Add(diff, plusOne)
	num, err := rand.Int(rand.Reader, rangeBig)
	if err != nil {
		return 0, err
	}
	res := new(big.Int).Add(num, bgMin)
	return res.Int64(), nil
}

func CreateRandInt(min, max int64) (int64, error) {
	return createRandValue(min, max)
}

func CreateRandFloat(min, max float64) (float64, error) {
	convertFromFloatToInt := func(num float64) (int64, int64) {
		str := strconv.FormatFloat(num, 'f', -1, 64)
		idx := strings.Index(str, ".")
		decimalPlaces := 0
		if idx != -1 {
			decimalPlaces = len(str) - idx - 1
		}
		multiplier := int64(math.Pow10(decimalPlaces))
		res := int64(math.Round(num * float64(multiplier)))
		return res, multiplier
	}
	start, startMult := convertFromFloatToInt(min)
	end, endMult := convertFromFloatToInt(max)
	var commonMult int64
	if startMult > endMult {
		commonMult = startMult
	} else {
		commonMult = endMult
	}
	startScaled := start * (commonMult / startMult)
	endScaled := end * (commonMult / endMult)
	randInt, err := createRandValue(startScaled, endScaled)
	if err != nil {
		return 0, err
	}
	result := float64(randInt) / float64(commonMult)
	return result, nil
}

func GetRandArrVal[T any](slice []T) (T, error) {
	if len(slice) == 0 {
		var zero T
		return zero, errors.New("slice is empty")
	}
	randIdx, err := createRandValue(0, int64(len(slice)-1))
	if err != nil {
		var zero T
		return zero, err
	}
	return slice[randIdx], nil
}
