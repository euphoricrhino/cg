package main

import (
	"flag"

	cg "github.com/euphoricrhino/cg/lib"
)

var (
	j1 = flag.String("j1", "", "j1 value")
	j2 = flag.String("j2", "", "j2 value")
)

func main() {
	flag.Parse()

	twoj1, err := cg.ParseHalfInteger(*j1)
	if err != nil {
		panic(err)
	}
	twoj2, err := cg.ParseHalfInteger(*j2)
	if err != nil {
		panic(err)
	}
	t := cg.ComputeCG(twoj1, twoj2)
	t.RenderHTML()
}
