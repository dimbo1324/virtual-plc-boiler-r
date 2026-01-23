package converter

import (
	"math"
	"strconv"
	"strings"
)

func ConvertFromFloatToInt(num float64) (int64, int64) {
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
