package cg

import (
	"fmt"
	"math/big"
	"sync"
)

// A column of states, indexed by dj, for all states |j1+j2-dj,m>.
type column struct {
	t     *Table
	dj    int
	twoj  int
	cells []*cell
	wg    sync.WaitGroup
}

func newColumn(t *Table, dj int) *column {
	cellCount := (t.twoj1+t.twoj2)/2 + 1 - dj
	col := &column{
		t:     t,
		dj:    dj,
		twoj:  t.twoj1 + t.twoj2 - 2*dj,
		cells: make([]*cell, cellCount),
	}

	for i := 0; i < cellCount; i++ {
		// Determine the m1 range:
		// * -j1 <= m1 <= j1
		// * -j2 <= m2=m-m1 <= j2 ==> m-j2 <= m1 <= m+j2
		// We take the tighter bound max(m-j2, -j1) <= m1 <= min(m+j2, j1).
		min := -t.twoj1
		twom := col.twoj - 2*i
		if min < twom-t.twoj2 {
			min = twom - t.twoj2
		}
		max := t.twoj1
		if max > twom+t.twoj2 {
			max = twom + t.twoj2
		}
		col.cells[i] = newCell(min, max)
	}
	// Top cell of this column depends on its row peers before.
	col.wg.Add(dj)
	return col
}

func (col *column) compute() {
	// Wait for all dependency of the top cell of this column to be ready.
	col.wg.Wait()
	col.computeTop()

	// Go down the ladder by applying lowering operator.
	// Cross-referencing to group-nut pp225 eq (18), here is the mapping:
	// * j=j1+j2-dj;
	// * m is the "upper" state's z-spin, which corresponds to iteration i, which has z-spin value of j1+j2-dj-i, therefore m=j1+j2-dj-i;
	// * At level i+1 (lower rung), the coefficient corresponding to m1 is contributed from both the m1 and m1+1 coefficient of level i.
	//
	//                  column dj
	//
	//   i       ...m1-1 m1 m1+1 ...
	//            |  /|  /|  /|  /|
	//            | / | / | / | / |
	//            |/  |/  |/  |/  |
	//  i+1      ...m1-1 m1 m1+1 ...
	//
	//

	// All relevant values are scaled by 2 so we deal only with integers.
	for i := 0; i < len(col.cells)-1; i++ {
		twom := col.twoj - 2*i
		current, lower := col.cells[i], col.cells[i+1]
		for l := range lower.c {
			twom1 := lower.twom1ForIndex(l)
			twom2 := twom - 2 - twom1
			lower.c[l] = big.NewRat(0, 1)
			// Contribution from m1+1 term in current.
			if current.isGoodTwom1(twom1 + 2) {
				// √((j1+1+m1)(j1-m1))
				r := big.NewRat(int64((col.t.twoj1+2+twom1)*(col.t.twoj1-twom1)), 4)
				accum(lower.c[l], r.Mul(r, current.get(twom1+2)))
			}
			// Contribution from m1 term in current.
			if current.isGoodTwom1(twom1) {
				// √((j2+1+m2)(j2-m2))
				r := big.NewRat(int64((col.t.twoj2+2+twom2)*(col.t.twoj2-twom2)), 4)
				accum(lower.c[l], r.Mul(r, current.get(twom1)))
			}
			// 1/√((j+1-m)(j+m))
			r := big.NewRat(4, int64((col.twoj+2-twom)*(col.twoj+twom)))
			lower.c[l].Mul(lower.c[l], r)
		}
		// Unblock one dependency of the last column of the row.
		if col.dj+i+1 < len(col.t.columns) {
			col.t.columns[col.dj+i+1].wg.Done()
		}
	}
}

func (col *column) computeTop() {
	topCell := col.cells[0]
	topCell.c[0] = big.NewRat(1, 1)
	if col.dj == 0 {
		// Init case.
		return
	}

	rowPeer := func(dj int) *cell { return col.t.cell(dj, col.dj) }

	c0 := topCell.c[0]
	// Normalization constraint for 0th coefficient (one corresponding to the max m1), sign is positive by convention, see Shankar (15.2.10).
	for dj := 0; dj < col.dj; dj++ {
		c0.Sub(c0, blankRat().Abs(rowPeer(dj).c[0]))
	}
	if c0.Sign() <= 0 {
		panic(fmt.Sprintf("non-postive sign for square of coefficient <m1,m2| at column %v", col.dj))
	}

	// The remaining coefficients in topCell.
	// Calculated using the orthogonality constraint between lth and 0th.
	for l := 1; l <= col.dj; l++ {
		topCell.c[l] = big.NewRat(0, 1)
		cl := topCell.c[l]
		for dj := 0; dj < col.dj; dj++ {
			peer := rowPeer(dj)
			accum(cl, blankRat().Mul(peer.c[0], peer.c[l]))
		}
		cl.Quo(cl, c0).Neg(cl)
	}
}

// Accumulates v onto sum (both are to be interpreted as square of the underlying rational values with sign on the enumerator).
func accum(sum, v *big.Rat) {
	// Determine the overall sign.
	n1 := blankInt().Mul(sum.Num(), v.Denom())
	n2 := blankInt().Mul(sum.Denom(), v.Num())
	overallSign := n1.Add(n1, n2).Sign()

	if overallSign == 0 {
		sum.SetFrac64(0, 1)
		return
	}

	abs1 := blankRat().Abs(sum)
	abs2 := blankRat().Abs(v)

	// Cross term.
	cross := blankRat().Mul(abs1, abs2)
	// Verify cross term is perfect square of rational.
	// Note this is a much stronger condition than the condition that CG coefficients themselves are rational squares.
	numRoot := blankInt().Sqrt(cross.Num())
	r := blankInt().Mul(numRoot, numRoot)
	if r.Cmp(cross.Num()) != 0 {
		panic(fmt.Sprintf("numerator of cross term (%v, %v) is not square", sum, v))
	}
	denomRoot := blankInt().Sqrt(cross.Denom())
	r = r.Mul(denomRoot, denomRoot)
	if r.Cmp(cross.Denom()) != 0 {
		panic(fmt.Sprintf("denominator of cross term (%v, %v) is not square", sum, v))
	}
	cross = cross.SetFrac(numRoot, denomRoot)
	crossFactor := big.NewRat(2*int64(sum.Num().Sign()*v.Num().Sign()), 1)

	// Add to sum the sum of squares and cross term with proper sign and factor 2.
	sum.Add(abs1, abs2).Add(sum, cross.Mul(cross, crossFactor))

	// Apply the overall sign.
	if overallSign < 0 {
		sum.Neg(sum)
	}
}

func blankInt() *big.Int { return big.NewInt(0) }

func blankRat() *big.Rat { return big.NewRat(0, 1) }
