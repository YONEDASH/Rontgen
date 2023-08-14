package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"unicode"
)

const (
	Version = "1.0"
)

type Configuration struct {
	Verbose bool
	Path    string
	Pattern *regexp.Regexp
}

type Match struct {
	Path    string
	Row     int
	Column  int
	Length  int
	Matched string
	Line    string
}

func Rontgen(config *Configuration) []Match {
	matches := []Match{}

	isRootDir, err := isDir(config.Path)

	if err != nil {
		fmt.Println("Error while determining root type:", err)
		return matches
	}

	if isRootDir {
		scanDir(config.Path, config, &matches)
	} else {
		scanFile(config.Path, config, &matches)
	}

	return matches
}

func isDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

func scanDir(path string, config *Configuration, matches *[]Match) {
	dirEntry, err := os.ReadDir(path)

	if err != nil {
		fmt.Println("Error opening directory:", err)
		return
	}

	for _, entry := range dirEntry {
		entryPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			scanDir(entryPath, config, matches)
			continue
		}

		scanFile(entryPath, config, matches)
	}
}

func scanFile(path string, config *Configuration, matches *[]Match) {
	content, err := os.ReadFile(path)

	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	if isContentBinary(content) {
		if config.Verbose {
			fmt.Printf("%s is binary\n", path)
		}

		return
	}

	text := string(content)
	locations := config.Pattern.FindAllStringIndex(text, -1)

	if len(locations) == 0 {
		if config.Verbose {
			fmt.Printf("%s has 0 matches\n", path)
		}

		return
	}

	textLength := len(text)
	lineFeedIndices := []int{}

	for i := 0; i < textLength; i++ {
		if text[i] == '\n' {
			lineFeedIndices = append(lineFeedIndices, i)
		}
	}

	rowCount := len(lineFeedIndices)

	for _, location := range locations {
		startIndex, endIndex := location[0], location[1]
		matchLen := endIndex - startIndex

		lineIndex, row := getIndexAndRow(startIndex, lineFeedIndices)
		column := startIndex - lineIndex

		viewStart := lineIndex + 1 // + 1 to skip \n
		if row == 0 {
			viewStart = 0
		}

		viewEnd := textLength
		if row < rowCount {
			viewEnd = lineFeedIndices[row]
		}

		line := text[viewStart:viewEnd]
		matched := text[startIndex:endIndex]

		match := Match{
			Path:    path,
			Row:     row,
			Column:  column,
			Length:  matchLen,
			Matched: matched,
			Line:    line,
		}

		*matches = append(*matches, match)
	}
}

func isContentBinary(data []byte) bool {
	for _, b := range data {
		if b == 0 {
			return true
		}
		if !unicode.IsPrint(rune(b)) && !unicode.IsSpace(rune(b)) {
			return true
		}
	}
	return false
}

func getIndexAndRow(where int, lineFeedIndices []int) (int, int) {
	index := 0
	row := 0

	for _, lineIndex := range lineFeedIndices {
		if lineIndex < where {
			index = lineIndex
		} else {
			break
		}
		row++
	}

	return index, row
}
