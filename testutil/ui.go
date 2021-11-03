package testutil

import (
	"strings"
	"unicode"
)

func GetTableOutputLines(headers []string, rows [][]string) []string {
	const padding = 3
	var lines []string
	maxLen := make(map[int]int)
	getMaxLen(headers, rows, maxLen)
	lines = append(lines, getLine(headers, maxLen, padding))
	for _, row := range rows {
		lines = append(lines, getLine(row, maxLen, padding))
	}
	return lines
}

func getLine(row []string, maxLen map[int]int, padding int) string {
	var line string
	for i := range row {
		line = line + row[i] + strings.Repeat(" ", getSpaceCount(i, row, maxLen, padding))
	}
	return strings.TrimRightFunc(line, unicode.IsSpace)
}

func getSpaceCount(i int, row []string, maxLen map[int]int, padding int) int {
	var count int
	if i < len(row)-1 {
		count = maxLen[i] - len(row[i]) + padding
	} else {
		count = padding
	}
	return count
}

func getMaxLen(headers []string, rows [][]string, maxLen map[int]int) {
	getRowMaxLen(headers, maxLen)
	for _, row := range rows {
		getRowMaxLen(row, maxLen)
	}
}

func getRowMaxLen(row []string, maxLen map[int]int) {
	for i, col := range row {
		if len(col) > maxLen[i] {
			maxLen[i] = len(col)
		}
	}
}
