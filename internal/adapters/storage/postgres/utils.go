package postgres

import "strconv"

func ParseFromStringToFloat64(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0, err
	}

	return f, nil
}

func ParseFloat64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', 4, 64)
}