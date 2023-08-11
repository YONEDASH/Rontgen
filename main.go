package main

import (
	"flag"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Match struct {
	Path   string
	Line   string
	Row    int
	Column int
}

var Version string = "1.0"
var ViewRange int = 10

func main() {
	fmt.Printf("GrepClone v%s\n", Version)

	patternPtr := flag.String("p", "", "Pattern")
	dirPathPtr := flag.String("d", "", "Path of directory")
	flag.Parse()

	patternStr := *patternPtr

	if len(patternStr) == 0 {
		fmt.Println("No pattern declared")
		return
	}

	dirPath := *dirPathPtr

	if len(dirPath) == 0 {
		dirPath = "."
	}

	matches := []Match{}

	pattern, err := regexp.Compile(patternStr)
	if err != nil {
		fmt.Println("Error compiling regex: ", err)
		return
	}

	ScanDir(pattern, dirPath, &matches)

	for _, match := range matches {
		fmt.Printf("%s @ line %d, column %d: %s\n", match.Path, match.Row, match.Column, match.Line)
	}

	if len(matches) == 0 {
		fmt.Println("No matches found")
	}

}

func ScanDir(pattern *regexp.Regexp, path string, matches *[]Match) {
	// Open dir
	dirEntry, err := os.ReadDir(path)

	if err != nil {
		fmt.Println("Error opening directory: ", err)
		return
	}

	for _, entry := range dirEntry {
		entryPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			ScanDir(pattern, entryPath, matches)
			continue
		}

		ScanFile(entryPath, entry, pattern, matches)
	}
}

func ScanFile(entryPath string, entry fs.DirEntry, pattern *regexp.Regexp, matches *[]Match) {
	fileContent, err := os.ReadFile(entryPath)

	if err != nil {
		fmt.Println("Error reading file: ", err)
		return
	}

	contentString := string(fileContent)

	locations := pattern.FindAllStringIndex(contentString, -1)

	if len(locations) == 0 {
		return
	}

	newLineIndexes := []int{}
	contentLength := len(contentString)

	for i := 0; i < contentLength; i++ {
		if contentString[i] == '\n' {
			newLineIndexes = append(newLineIndexes, i)
		}
	}

	for _, loc := range locations {
		startIndex, endIndex := loc[0], loc[1]

		lineIndex, row := GetLineIndex(startIndex, newLineIndexes)

		col := startIndex - lineIndex

		viewStart := int(math.Max(0, float64(startIndex-ViewRange)))
		viewEnd := int(math.Min(float64(contentLength), float64(endIndex+ViewRange)))

		line := ".." + CleanUpLine(contentString[viewStart:viewEnd]) + ".."

		match := Match{
			Path:   entryPath,
			Line:   line,
			Row:    row,
			Column: col,
		}

		*matches = append(*matches, match)
	}
}

func GetLineIndex(where int, lineIndexes []int) (int, int) {
	index := 0
	inArray := 0

	for _, lineIndex := range lineIndexes {
		inArray++
		if lineIndex < where {
			index = lineIndex
		} else {
			break
		}
	}

	return index, inArray
}

// Remove new lines and tabs
func CleanUpLine(str string) string {
	return strings.ReplaceAll(strings.ReplaceAll(str, "\n", " "), "\t", " ")
}
