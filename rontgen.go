package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"unicode"
)

var Version = "dev"

type Configuration struct {
	Verbose  bool
	Path     string
	Pattern  *regexp.Regexp
	DepthCap int
}

type Match struct {
	Path      string
	NameMatch bool
	Row       int
	Column    int
	Length    int
	Matched   string
	Line      string
}

func Rontgen(config *Configuration) []Match {
	matches := []Match{}

	isRootDir, err := isDir(config.Path)

	if err != nil {
		fmt.Println(Red, "\bError while determining root type:", err, Reset)
		return matches
	}

	if isRootDir {
		scanDir(config.Path, config, &matches, 0)
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

func scanDir(path string, config *Configuration, matches *[]Match, depth int) {
	if depth > config.DepthCap {
		if config.Verbose {
			fmt.Println("Reached depth cap of", config.DepthCap)
		}

		return
	}

	dirEntry, err := os.ReadDir(path)

	if err != nil {
		fmt.Println(Red, "\bError opening directory:", err, Reset)
		return
	}

	for _, entry := range dirEntry {
		entryPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			scanDir(entryPath, config, matches, depth+1)
			continue
		}

		scanFile(entryPath, config, matches)
	}
}

func scanFile(path string, config *Configuration, matches *[]Match) {
	// Check if pattern matches file name
	nameLocation := config.Pattern.FindStringIndex(path)
	if len(nameLocation) > 0 {
		match := Match{
			Path:      path,
			NameMatch: true,
			Matched:   path[nameLocation[0]:nameLocation[1]],
			Column:    nameLocation[0],
			Length:    nameLocation[1] - nameLocation[0],
		}
		*matches = append(*matches, match)
	}

	// Check if pattern matches file content
	content, err := os.ReadFile(path)

	if err != nil {
		fmt.Println(Red, "\bError reading file:", err, Reset)
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
			Path:      path,
			NameMatch: false,
			Row:       row,
			Column:    column,
			Length:    matchLen,
			Matched:   matched,
			Line:      line,
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
