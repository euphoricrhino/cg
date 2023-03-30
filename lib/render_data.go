package cg

type tableData struct {
	J1       string
	J2       string
	Sections []*sectionData
}

type sectionData struct {
	M            string
	PrintHeading bool
	Rows         []*rowData
}

type rowData struct {
	M1     string
	M2     string
	Values []string
}
