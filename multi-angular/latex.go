package main

import (
	"fmt"
	"strconv"
	"strings"

	cg "github.com/euphoricrhino/cg/lib"
)

func pathLatex(path []int) string {
	s := make([]string, len(path))
	for i, p := range path {
		s[i] = halfIntegerLatex(p)
	}
	return fmt.Sprintf("\\left[%v\\right]", strings.Join(s, ","))
}

func halfIntegerLatex(twov int) string {
	str := ""
	if twov < 0 {
		str += "-"
		twov = -twov
	}
	if twov%2 == 0 {
		str += strconv.Itoa(twov / 2)
	} else {
		str += fmt.Sprintf("\\frac{%v}{2}", twov)
	}
	return str
}

func jmLatex(twoj, twom int) string {
	return fmt.Sprintf("\\left|%v,%v\\right\\rangle", halfIntegerLatex(twoj), halfIntegerLatex(twom))
}

func stateLatex(st *state) string {
	str := ""
	num := cg.BlankInt().Abs(st.c.Num())
	if st.c.Sign() < 0 {
		str += "-"
	}
	str += fmt.Sprintf("\\sqrt{\\frac{%v}{%v}}", num, st.c.Denom())
	str += jmLatex(st.twoj, st.twom)
	return str
}
