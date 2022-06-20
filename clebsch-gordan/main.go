package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/euphoricrhino/cg/cg"
)

var (
	j1 = flag.String("j1", "", "j1 value")
	j2 = flag.String("j2", "", "j2 value")
)

// Parses half integer from command line and returns its twice value integer.
func parseHalfInteger(str string) int {
	parts := strings.Split(str, "/")
	if len(parts) == 1 {
		v, err := strconv.Atoi(parts[0])
		if err != nil {
			panic(fmt.Sprintf("invalid value for integer: '%v'", str))
		}
		return 2 * v
	}
	if len(parts) != 2 || parts[1] != "2" {
		panic(fmt.Sprintf("invalid format for half integer: '%v'", str))
	}
	v, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(fmt.Sprintf("invalid value for half integer: '%v'", str))
	}
	return v
}

func main() {
	flag.Parse()

	twoj1 := parseHalfInteger(*j1)
	twoj2 := parseHalfInteger(*j2)
	t := cg.ComputeCG(twoj1, twoj2)
	t.RenderHTML()
}
