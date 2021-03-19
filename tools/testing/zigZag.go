package testing

import "strings"

func convert(s string, numRows int) string {
	if len(s) == 0 || numRows == 0 { return "" }
	if numRows == 1 { return s }

	var (
		rows = make([][]rune, min(len(s), numRows), min(len(s), numRows))
		curRow = 0
		nextIsDown = true
		zigZag = ""
	)
	strings.Split()
	for _, c := range s {
		rows[curRow] = append(rows[curRow], c)
		if nextIsDown {
			curRow ++
		} else {
			curRow --
		}
		if curRow == 0 || curRow == numRows - 1 {
			nextIsDown = !nextIsDown
		}
	}

	for _, row := range rows {
		zigZag += string(row)
	}

	return zigZag
}

func min(v ...int) int {
	if len(v) == 0 { return 0 }
	min := v[0]
	for _, value := range v {
		if value < min {
			min = value
		}
	}
	return min
}