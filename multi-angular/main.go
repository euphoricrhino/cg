package main

import (
	"flag"
)

var (
	states = flag.String("states", "", "j1,m1;j2,m2[;...;jk,mk]")
)

func main() {
	flag.Parse()

	ma, err := computeMultiAngular(*states)
	if err != nil {
		panic(err)
	}

	ma.RenderHTML()
}
