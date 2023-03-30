package cg

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"sync"
)

// Table represents the table of the CG coefficient.
type Table struct {
	// Whether we are storing index 1 and 2 in exchanged order due to original j1 < j2.
	exchanged bool
	twoj1     int
	twoj2     int
	columns   []*column
}

// ComputeCG computes the CG table for the given j1 and j2.
// Arguments are twice the value of actual j1 and j2 so they are integers.
func ComputeCG(twoj1, twoj2 int) *Table {
	if twoj1 <= 0 || twoj2 <= 0 {
		panic(fmt.Sprintf("invalid j1 or j2: %v, %v", twoj1, twoj2))
	}
	exchanged := false
	if twoj1 < twoj2 {
		twoj1, twoj2 = twoj2, twoj1
		exchanged = true
	}
	t := &Table{
		exchanged: exchanged,
		twoj1:     twoj1,
		twoj2:     twoj2,
		columns:   make([]*column, twoj2+1),
	}

	for dj := 0; dj <= twoj2; dj++ {
		t.columns[dj] = newColumn(t, dj)
	}

	// One goroutine per column.
	var wg sync.WaitGroup
	wg.Add(twoj2 + 1)
	for _, col := range t.columns {
		go func(c *column) {
			c.compute()
			wg.Done()
		}(col)
	}
	wg.Wait()
	return t
}

// Gets the cell representing state |j1+j2-dj,j1+j2-dm>.
func (t *Table) cell(dj, dm int) *cell {
	return t.columns[dj].cells[dm-dj]
}

func (t *Table) RenderHTML() {
	data := t.getTableData()
	filename := filepath.Join(os.TempDir(), "clebsch-gordon.html")
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := rootTmpl.Execute(f, data); err != nil {
		panic(err)
	}

	fmt.Println(filename)
}

// Query queries the CG table for the value
// ⟨j1,m1;j2,m2|j,m⟩, where j1 and j2 are the same values (in this order) used to create this table.
// All arguments are twice the actual values so they are integers.
func (t *Table) Query(twoj, twom, twom1, twom2 int) *big.Rat {
	return t.queryHelper(twoj, twom, twom1, twom2, false)
}

// ExchangedQuery queries the CG table for the value
// ⟨j2,m2;j1,m1|j,m⟩, where j1 and j2 are the same values (in this order) used to create this table.
// All arguments are twice the actual values so they are integers.
func (t *Table) ExchangedQuery(twoj, twom, twom1, twom2 int) *big.Rat {
	return t.queryHelper(twoj, twom, twom1, twom2, true)
}

func (t *Table) queryHelper(twoj, twom, twom1, twom2 int, exchangedQuery bool) *big.Rat {
	dj := t.twoj1 + t.twoj2 - twoj
	if dj%2 != 0 {
		return BlankRat()
	}
	if twom != twom1+twom2 {
		return BlankRat()
	}
	dj /= 2
	if dj < 0 || dj > t.twoj2 {
		return BlankRat()
	}
	col := t.columns[dj]
	mneg := false
	if twom < 0 {
		mneg = true
		twom = -twom
		twom1 = -twom1
	}
	// If exchanged, remember to exchange 1<->2.
	if t.exchanged {
		twom1 = twom - twom1
	}
	dm := t.twoj1 + t.twoj2 - twom
	if dm%2 != 0 {
		return BlankRat()
	}
	dm /= 2
	if dm-dj < 0 || dm-dj >= len(col.cells) {
		return BlankRat()
	}
	cell := col.cells[dm-dj]
	if !cell.isGoodTwom1(twom1) {
		return BlankRat()
	}
	ret := BlankRat().Set(cell.get(twom1))
	// Use CG coefficient symmetry property:
	// 1. ⟨j1,m1;j2,m2|j,m⟩=(-1)^{j1+j2-j}⟨j2,m2;j1,m1|j,m⟩
	// 2. ⟨j1,-m1;j2,-m2|j,-m⟩=(-1)^{j1+j2-j}⟨j1,m1;j2,m2|j,m⟩
	if (mneg != (t.exchanged != exchangedQuery)) && dj%2 != 0 {
		ret.Neg(ret)
	}
	return ret
}

func (t *Table) getTableData() *tableData {
	twoj1, twoj2 := t.twoj1, t.twoj2
	if t.exchanged {
		twoj1, twoj2 = twoj2, twoj1
	}
	return &tableData{
		J1:       FormatHalfInteger(twoj1),
		J2:       FormatHalfInteger(twoj2),
		Sections: t.getSectionsData(),
	}
}

func (t *Table) getSectionsData() []*sectionData {
	col0 := t.columns[0]
	data := make([]*sectionData, 0, col0.twoj+1)
	for i := range col0.cells {
		data = append(data, t.getSectionData(i, false))
	}
	rbegin := len(col0.cells) - 1
	// For whole integer j1+j2, don't include m=0 twice.
	if col0.twoj%2 == 0 {
		rbegin--
	}
	for i := rbegin; i >= 0; i-- {
		data = append(data, t.getSectionData(i, true))
	}
	return data
}

func (t *Table) getSectionData(i int, mirrored bool) *sectionData {
	col0 := t.columns[0]
	twom := col0.twoj - 2*i
	mStr := FormatHalfInteger(twom)
	data := &sectionData{
		M:            mStr,
		PrintHeading: !mirrored && twom >= (t.twoj1-t.twoj2),
		Rows:         make([]*rowData, 0, len(col0.cells[i].c)),
	}
	if mirrored {
		data.M = FormatHalfInteger(-twom)
	}
	for l := range col0.cells[i].c {
		twom1 := col0.cells[i].twom1ForIndex(l)
		twom2 := twom - twom1
		if t.exchanged {
			twom1, twom2 = twom2, twom1
		}
		row := &rowData{
			M1:     FormatHalfInteger(twom1),
			M2:     FormatHalfInteger(twom2),
			Values: make([]string, 0, i+1),
		}
		if mirrored {
			if twom1 != 0 {
				row.M1 = FormatHalfInteger(-twom1)
			}
			if twom2 != 0 {
				row.M2 = FormatHalfInteger(-twom2)
			}
		}
		for dj := 0; dj < i+1 && dj < len(t.columns); dj++ {
			col := t.columns[dj]
			value := ""
			// Use CG coefficient symmetry property:
			// 1. ⟨j1,m1;j2,m2|j,m⟩=(-1)^{j1+j2-j}⟨j2,m2;j1,m1|j,m⟩
			// 2. ⟨j1,-m1;j2,-m2|j,-m⟩=(-1)^{j1+j2-j}⟨j1,m1;j2,m2|j,m⟩
			if (t.exchanged != mirrored) && dj%2 != 0 {
				value = FormatRat(BlankRat().Neg(col.cells[i-dj].c[l]))
			} else {
				value = FormatRat(col.cells[i-dj].c[l])
			}
			row.Values = append(row.Values, value)
		}
		data.Rows = append(data.Rows, row)
	}
	return data
}
