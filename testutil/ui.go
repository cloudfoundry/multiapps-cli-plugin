package testutil

import "strings"

func GetTableOutputLines(headers []string, rows [][]string) []string {
	const padding = 3
	lines := []string{}
	maxlen := make(map[int]int)
	getMaxlen(headers, rows, maxlen)
	lines = append(lines, getLine(headers, maxlen, padding))
	for _, row := range rows {
		lines = append(lines, getLine(row, maxlen, padding))
	}
	return lines
}

func getLine(row []string, maxlen map[int]int, padding int) string {
	var line string
	for i := range row {
		line = line + row[i] + strings.Repeat(" ", getSpaceCount(i, row, maxlen, padding))
	}
	return line + "\n"
}

func getSpaceCount(i int, row []string, maxlen map[int]int, padding int) int {
	var count int
	if i < len(row)-1 {
		count = maxlen[i] - len(row[i]) + padding
	} else {
		count = padding
	}
	return count
}

func getMaxlen(headers []string, rows [][]string, maxlen map[int]int) {
	getRowMaxlen(headers, maxlen)
	for _, row := range rows {
		getRowMaxlen(row, maxlen)
	}
}

func getRowMaxlen(row []string, maxlen map[int]int) {
	for i, col := range row {
		if len(col) > maxlen[i] {
			maxlen[i] = len(col)
		}
	}
}
