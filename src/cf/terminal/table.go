package terminal

import (
	"fmt"
	"strings"
)

type Table interface {
	Print(rows [][]string)
}

type PrintableTable struct {
	ui            UI
	header        []string
	headerPrinted bool
	maxSizes      []int
}

func NewTable(ui UI, header []string) Table {
	return &PrintableTable{
		ui:       ui,
		header:   header,
		maxSizes: make([]int, len(header)),
	}
}

func (t *PrintableTable) Print(rows [][]string) {
	for _, row := range append(rows, t.header) {
		t.calculateMaxSize(row)
	}

	if t.headerPrinted == false {
		t.printHeader()
		t.headerPrinted = true
	}

	for _, line := range rows {
		t.printRow(line)
	}
}

func (t *PrintableTable) calculateMaxSize(row []string) {
	for index, value := range row {
		cellLength := len(decolorize(value))
		if t.maxSizes[index] < cellLength {
			t.maxSizes[index] = cellLength
		}
	}
}

func (t *PrintableTable) printHeader() {
	output := ""
	for col, value := range t.header {
		output = output + t.cellValue(col, HeaderColor(value))
	}
	t.ui.Say(output)
}

func (t *PrintableTable) printRow(row []string) {
	output := ""
	for columnIndex, value := range row {
		if columnIndex == 0 {
			value = TableContentHeaderColor(value)
		}

		output = output + t.cellValue(columnIndex, value)
	}
	t.ui.Say("%s", output)
}

func (t *PrintableTable) cellValue(col int, value string) string {
	padding := ""
	if col < len(t.header)-1 {
		padding = strings.Repeat(" ", t.maxSizes[col]-len(decolorize(value)))
	}
	return fmt.Sprintf("%s%s   ", value, padding)
}
