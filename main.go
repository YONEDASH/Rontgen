package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"unicode"
)

type Match struct {
	Path   string
	Line   string
	Row    int
	Column int
}

var Version string = "1.0"
var ViewRange int = 10
var AllowDuplicates bool = false

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [-v] [-i] <file_path> <pattern>\n", os.Args[0])
		fmt.Println("Flags:")
		flag.PrintDefaults()
		fmt.Println("Pattern:")
		fmt.Printf("  <pattern>\tPattern to search for\n")
	}

	versionPtr := flag.Bool("v", false, "Show version")
	infoPtr := flag.Bool("i", false, "Show info for matches")
	duplicatesPtr := flag.Bool("d", false, "Allow duplicate lines on multiple matches")

	flag.Parse()

	if *versionPtr {
		fmt.Printf("GrepClone v%s\n", Version)
		return
	}

	if *duplicatesPtr {
		AllowDuplicates = true
	}

	nonFlagArgs := flag.Args()

	if len(nonFlagArgs) == 0 {
		fmt.Println("No pattern declared")
		return
	}

	if len(nonFlagArgs) > 2 {
		fmt.Println("Too many arguments")
		return
	}

	patternStr := nonFlagArgs[0]

	if len(patternStr) == 0 {
		fmt.Println("Invalid pattern declared")
		return
	}

	path := "."

	if len(nonFlagArgs) == 2 {
		path = nonFlagArgs[1]
	}

	matches := []Match{}

	pattern, err := regexp.Compile(patternStr)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return
	}

	pathInfo, err := GetPath(path)

	if err != nil {
		fmt.Println("Error while checking path:", err)
		return
	}

	if pathInfo.IsDir() {
		ScanDir(pattern, path, &matches)
	} else {
		ScanFile(pathInfo.Name(), pattern, &matches)
	}

	if *infoPtr {
		for _, match := range matches {
			fmt.Printf("%s @ line %d, column %d: %s\n", match.Path, match.Row, match.Column, match.Line)
		}
	} else {
		linesPrinted := []int{}

		for _, match := range matches {
			if !AllowDuplicates && ContainsInt(match.Row, linesPrinted) {
				continue
			}
			linesPrinted = append(linesPrinted, match.Row)

			fmt.Printf("%s\n", match.Line)
		}

	}

	if len(matches) == 0 {
		fmt.Println("No matches found")
	}

}

func ContainsInt(target int, slice []int) bool {
	for _, num := range slice {
		if num == target {
			return true
		}
	}
	return false
}

func GetPath(path string) (fs.FileInfo, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	return fileInfo, nil
}

func ScanDir(pattern *regexp.Regexp, path string, matches *[]Match) {
	// Open dir
	dirEntry, err := os.ReadDir(path)

	if err != nil {
		fmt.Println("Error opening directory:", err)
		return
	}

	for _, entry := range dirEntry {
		entryPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			ScanDir(pattern, entryPath, matches)
			continue
		}

		ScanFile(entryPath, pattern, matches)
	}
}

func ScanFile(entryPath string, pattern *regexp.Regexp, matches *[]Match) {
	fileContent, err := os.ReadFile(entryPath)

	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	if IsBinary(fileContent) {
		fmt.Printf("%s is a binary\n", entryPath)
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
		startIndex := loc[0]

		lineIndex, row := GetLineIndex(startIndex, newLineIndexes)

		col := startIndex - lineIndex

		viewStart := lineIndex + 1
		if row == 0 {
			viewStart = 0
		}

		viewEnd := contentLength
		if row < len(newLineIndexes) {
			viewEnd = newLineIndexes[row]
		}

		line := contentString[viewStart:viewEnd]

		match := Match{
			Path:   entryPath,
			Line:   line,
			Row:    row,
			Column: col,
		}

		*matches = append(*matches, match)
	}
}

func IsBinary(data []byte) bool {
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

func GetLineIndex(where int, lineIndexes []int) (int, int) {
	index := 0
	inArray := 0

	for _, lineIndex := range lineIndexes {
		if lineIndex < where {
			index = lineIndex
		} else {
			break
		}
		inArray++
	}

	return index, inArray
}
