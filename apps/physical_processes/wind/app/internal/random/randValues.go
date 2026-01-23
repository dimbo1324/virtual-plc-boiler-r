package random

import (
	"crypto/rand"
	"math/big"
	"wind-process/internal/utils/converter"
)

func createRandValue(min, max int64) (int64, error) {
	if min > max {
		min, max = max, min
	}
	rangeValue := max - min + 1
	rangeBig := big.NewInt(rangeValue)
	num, err := rand.Int(rand.Reader, rangeBig)
	if err != nil {
		return 0, err
	}
	res := num.Int64() + min
	return res, nil
}

func CreateRandInt(min, max int64) (int64, error) {
	return createRandValue(min, max)
}

func CreateRandFloat(min, max float64) (float64, error) {
	start, startMult := converter.ConvertFromFloatToInt(min)
	end, endMult := converter.ConvertFromFloatToInt(max)
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
