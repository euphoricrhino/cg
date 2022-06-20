package cg

import (
	"math/big"
)

// Represents a cell in the CG table. Each cell is identified by the the state |j,m>.
// Each cell contains a set of non-vanishing CG coefficients, which corresponds different combinations of
// <m1,m2|j,m>, where m1+m2=m.
type cell struct {
	// The range of the m1 value (doubled so we store only integers).
	minTwom1 int
	maxTwom1 int
	// Indices correspond to decreasing m1 value.
	c []*big.Rat
}

func newCell(minTwom1, maxTwom1 int) *cell {
	return &cell{
		minTwom1: minTwom1,
		maxTwom1: maxTwom1,
		c:        make([]*big.Rat, (maxTwom1-minTwom1)/2+1),
	}
}

// Given index to the coefficient slice, returns the corresponding 2m1 value.
func (c *cell) twom1ForIndex(idx int) int {
	return c.maxTwom1 - 2*idx
}

// Checks if the given 2m1 value is valid for this cell.
func (c *cell) isGoodTwom1(twom1 int) bool {
	return twom1 >= c.minTwom1 && twom1 <= c.maxTwom1 && (c.maxTwom1-twom1)%2 == 0
}

// Gets the coefficient for the given 2m1 value.
func (c *cell) get(twom1 int) *big.Rat {
	return c.c[(c.maxTwom1-twom1)/2]
}
