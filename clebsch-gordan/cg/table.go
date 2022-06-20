package cg

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

// Table represents the table of the CG coefficient.
type Table struct {
	twoj1   int
	twoj2   int
	columns []*column
}

func ComputeCG(twoj1, twoj2 int) *Table {
	if twoj1 <= 0 || twoj2 <= 0 {
		panic(fmt.Sprintf("invalid j1 or j2: %v, %v", twoj1, twoj2))
	}
	if twoj1 < twoj2 {
		twoj1, twoj2 = twoj2, twoj1
	}
	t := &Table{
		twoj1:   twoj1,
		twoj2:   twoj2,
		columns: make([]*column, twoj2+1),
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

func (t *Table) getTableData() *tableData {
	return &tableData{
		J1:       formatHalfInteger(t.twoj1),
		J2:       formatHalfInteger(t.twoj2),
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
	mStr := formatHalfInteger(twom)
	data := &sectionData{
		M:            mStr,
		PrintHeading: !mirrored && twom >= (t.twoj1-t.twoj2),
		Rows:         make([]*rowData, 0, len(col0.cells[i].c)),
	}
	if mirrored {
		data.M = formatHalfInteger(-twom)
	}
	for l := range col0.cells[i].c {
		twom1 := col0.cells[i].twom1ForIndex(l)
		twom2 := twom - twom1
		row := &rowData{
			M1:     formatHalfInteger(twom1),
			M2:     formatHalfInteger(twom2),
			Values: make([]string, 0, i+1),
		}
		if mirrored {
			if twom1 != 0 {
				row.M1 = formatHalfInteger(-twom1)
			}
			if twom2 != 0 {
				row.M2 = formatHalfInteger(-twom2)
			}
		}
		for dj := 0; dj < i+1 && dj < len(t.columns); dj++ {
			col := t.columns[dj]
			value := ""
			if mirrored && dj%2 != 0 {
				value = formatRat(blankRat().Neg(col.cells[i-dj].c[l]))
			} else {
				value = formatRat(col.cells[i-dj].c[l])
			}
			row.Values = append(row.Values, value)
		}
		data.Rows = append(data.Rows, row)
	}
	return data
}

func formatHalfInteger(twiceValue int) string {
	if twiceValue%2 == 0 {
		return strconv.Itoa(twiceValue / 2)
	}
	return fmt.Sprintf("%v/2", twiceValue)
}

func formatRat(r *big.Rat) string {
	if r.Sign() == 0 {
		return "0"
	}
	return r.String()
}
