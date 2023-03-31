package main

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	cg "github.com/euphoricrhino/cg/lib"
)

var (
	errFormat = fmt.Errorf("input must be in the format j1,m1;j2,m2[;...;jk,mk]")
)

// Represents an angular momentum eigenstate |j,m⟩ times a complex coefficient.
type state struct {
	c    *big.Rat
	twoj int
	twom int
	// Unique path identifying the subspace of dimension 2j+1.
	subspacePath []int
}

// Represents the expansion of multiple |ji,mi⟩ tensor products into total angular momentum basis |j,m⟩.
type multiAngular struct {
	// Input states in the tensor product form.
	inputStates []*state
	// Eventual states expanded in total angular momentum basis |j,m⟩.
	expandedStates []*state

	// Alphabetically sorted subspace paths.
	subspacePaths [][]int
	// Lookup map from subspace path to the subspace index.
	subspaceIndex map[string]int
}

func newState(jmStr string) (*state, error) {
	jmParts := strings.Split(jmStr, ",")
	if len(jmParts) != 2 {
		return nil, errFormat
	}
	twoj, err := cg.ParseHalfInteger(jmParts[0])
	if err != nil {
		return nil, err
	}
	if twoj < 0 {
		return nil, fmt.Errorf("invalid j value: %v", jmParts[0])
	}
	twom, err := cg.ParseHalfInteger(jmParts[1])
	if err != nil {
		return nil, err
	}
	if twom > twoj || twom < -twoj {
		return nil, fmt.Errorf("invalid m value for j=%v: %v", jmParts[1], jmParts[1])
	}
	return &state{
		c:            big.NewRat(1, 1),
		twoj:         twoj,
		twom:         twom,
		subspacePath: []int{twoj},
	}, nil
}

// Computes the multi angular decomposition given the input states.
func computeMultiAngular(statesStr string) (*multiAngular, error) {
	parts := strings.Split(statesStr, ";")
	if len(parts) <= 1 {
		return nil, errFormat
	}
	var inputStates []*state
	// We are putting partially expanded states in head.
	head := make([]*state, 1)
	// All remaining input states are in tail.
	tail := make([]*state, len(parts)-1)
	for i, part := range parts {
		st, err := newState(part)
		if err != nil {
			return nil, err
		}
		inputStates = append(inputStates, st)
		if i == 0 {
			head[0] = st
		} else {
			tail[i-1] = st
		}
	}
	// Calculate dimensions of all irreducible subspaces.
	var prefix []int
	queue := [][]int{{head[0].twoj}}
	for _, st := range tail {
		twoj2 := st.twoj
		qlen := len(queue)
		for i := 0; i < qlen; i++ {
			prefix, queue = queue[0], queue[1:]
			twoj1 := prefix[len(prefix)-1]
			lo := twoj1 - twoj2
			if lo < 0 {
				lo = -lo
			}
			for twoj := lo; twoj <= twoj1+twoj2; twoj += 2 {
				queue = append(queue, appendCopy(prefix, twoj))
			}
		}
	}
	// Sort the subspaces by the total angular momentum (i.e., last element in the path).
	sort.Slice(queue, func(i, j int) bool {
		pi, pj := queue[i], queue[j]
		for i := len(pi) - 1; i >= 0; i-- {
			if pi[i] < pj[i] {
				return true
			}
			if pi[i] > pj[i] {
				return false
			}
		}
		return false
	})
	subspaceIndex := make(map[string]int)
	// Assign each path a unique key for quick reference later.
	for i, path := range queue {
		subspaceIndex[pathToSubspaceKey(path)] = i
	}

	tables := make(map[string]*cg.Table)
	// Now expand the tensor products into total angular momentum |j,m⟩ basis, consuming tail states one by one.
	for len(tail) > 0 {
		var st1, st2 *state
		st2, tail = tail[0], tail[1:]
		var hd []*state
		for _, st1 = range head {
			exchanged := false
			jmax, jmin := st1.twoj, st2.twoj
			if jmax < jmin {
				jmax, jmin = jmin, jmax
				exchanged = true
			}
			twom := st1.twom + st2.twom
			// Trivial case, one of the j's is zero.
			if jmin == 0 {
				st := &state{
					c:            st1.c,
					twoj:         jmax,
					twom:         twom,
					subspacePath: appendCopy(st1.subspacePath, jmax),
				}
				hd = append(hd, st)
				continue
			}
			tableKey := fmt.Sprintf("%v,%v", jmax, jmin)
			t, found := tables[tableKey]
			if !found {
				// Construct the CG table for j1,j2.
				fmt.Printf("constructing C-G table for j1=%v, j2=%v ...\n", cg.FormatHalfInteger(jmax), cg.FormatHalfInteger(jmin))
				t = cg.ComputeCG(jmax, jmin)
				tables[tableKey] = t
			}
			for twoj := jmax - jmin; twoj <= jmax+jmin; twoj += 2 {
				var c *big.Rat
				if !exchanged {
					c = t.Query(twoj, twom, st1.twom, st2.twom)
				} else {
					c = t.ExchangedQuery(twoj, twom, st2.twom, st1.twom)
				}
				if c.Sign() != 0 {
					st := &state{
						c:            cg.BlankRat().Mul(st1.c, c),
						twoj:         twoj,
						twom:         twom,
						subspacePath: appendCopy(st1.subspacePath, twoj),
					}
					hd = append(hd, st)
				}
			}
		}
		head = hd
	}
	return &multiAngular{
		inputStates:    inputStates,
		expandedStates: head,
		subspacePaths:  queue,
		subspaceIndex:  subspaceIndex,
	}, nil
}

func (ma *multiAngular) lookupSubspaceIndex(path []int) int {
	idx, found := ma.subspaceIndex[pathToSubspaceKey(path)]
	if !found {
		panic("subspaceIndex lookup failure")
	}
	return idx
}

// RenderHTML renders the multi angular decomposition.
func (ma *multiAngular) RenderHTML() {
	// Subspace compositions.
	latexStr := "\\mbox{irreducible subspace compositions} & &"
	for i, path := range ma.subspacePaths {
		if i == 0 {
			latexStr += fmt.Sprintf("%v&:%v", i, pathLatex(path))
		} else {
			latexStr += fmt.Sprintf("\\qquad %v:%v", i, pathLatex(path))
		}
	}
	latexStr += "\\\\\n"

	latexStr += "\\mbox{irreducible subspace dimensions} & &"
	// Tensor products of input dimensions.
	var s []string
	for _, st := range ma.inputStates {
		s = append(s, strconv.Itoa(st.twoj+1))
	}
	latexStr += strings.Join(s, "\\otimes ")
	latexStr += " &= "
	// Direct sums of irreducible subspace dimensions.
	s = nil
	for _, path := range ma.subspacePaths {
		twoj := path[len(path)-1]
		s = append(s, fmt.Sprintf("%v_{%v}", twoj+1, ma.lookupSubspaceIndex(path)))
	}
	latexStr += strings.Join(s, "\\oplus ")
	latexStr += "\\\\\n"

	// Tensor products of angular momenta.
	s = nil
	latexStr += "\\mbox{irreducible subspace total angular momenta} & &"
	for _, st := range ma.inputStates {
		s = append(s, halfIntegerLatex(st.twoj))
	}
	latexStr += strings.Join(s, "\\otimes ")
	latexStr += " &= "
	s = nil
	for _, path := range ma.subspacePaths {
		twoj := path[len(path)-1]
		if twoj%2 != 0 {
			s = append(s, fmt.Sprintf("\\left(%v\\right)_{%v}", halfIntegerLatex(twoj), ma.lookupSubspaceIndex(path)))
		} else {
			s = append(s, fmt.Sprintf("%v_{%v}", twoj, ma.lookupSubspaceIndex(path)))
		}
	}
	latexStr += strings.Join(s, "\\oplus ")
	latexStr += "\\\\\n"

	// Tensor product states.
	s = nil
	for _, st := range ma.inputStates {
		s = append(s, jmLatex(st.twoj, st.twom))
	}
	latexStr += "\\mbox{expansion in total angular momenta basis} & &"
	latexStr += strings.Join(s, "\\otimes ")
	latexStr += " &= "
	// Expanded states.
	if len(ma.expandedStates) == 0 {
		latexStr += "0"
	} else {
		for i, st := range ma.expandedStates {
			termStr := fmt.Sprintf("%v_{%v}", stateLatex(st), ma.lookupSubspaceIndex(st.subspacePath))
			if i != 0 && termStr[0] != '-' {
				latexStr += "+"
			}
			latexStr += termStr
		}
	}
	latexStr += "\\\\\n"

	filename := filepath.Join(os.TempDir(), "multi-angular.html")
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := tmpl.Execute(f, latexStr); err != nil {
		panic(err)
	}

	fmt.Println(filename)
}

func appendCopy(path []int, v int) []int {
	np := append([]int{}, path...)
	np = append(np, v)
	return np
}

func pathToSubspaceKey(path []int) string {
	s := make([]string, len(path))
	for i, v := range path {
		s[i] = strconv.Itoa(v)
	}
	return strings.Join(s, ",")
}
