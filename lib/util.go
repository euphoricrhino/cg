package cg

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

// ParseHalfInteger parses half integer from string representation of half integer and returns its twice value integer.
func ParseHalfInteger(str string) (int, error) {
	parts := strings.Split(str, "/")
	if len(parts) == 1 {
		v, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, fmt.Errorf("invalid value for integer: '%v'", str)
		}
		return 2 * v, nil
	}
	if len(parts) != 2 || parts[1] != "2" {
		return 0, fmt.Errorf("invalid format for half integer: '%v'", str)
	}
	v, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid value for half integer: '%v'", str)
	}
	return v, nil
}

// FormatHalfInteger prints the half integer's value.
func FormatHalfInteger(twiceValue int) string {
	if twiceValue%2 == 0 {
		return strconv.Itoa(twiceValue / 2)
	}
	return fmt.Sprintf("%v/2", twiceValue)
}

// FormatRat formats the big.Rat.
func FormatRat(r *big.Rat) string {
	if r.Sign() == 0 {
		return "0"
	}
	return r.String()
}

// BlankInt creates a new blank big.Int.
func BlankInt() *big.Int { return big.NewInt(0) }

// BlankRat creates a new blank big.Rat.
func BlankRat() *big.Rat { return big.NewRat(0, 1) }
