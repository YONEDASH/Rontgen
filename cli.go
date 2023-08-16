package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
)

func main() {
	config, errorCode := handleFlags()

	if config == nil {
		if errorCode {
			os.Exit(1)
		}

		return
	}

	matches := Rontgen(config)

	if len(matches) == 0 {
		return
	}

	// TODO: better print out
	currentPath := ""

	for _, match := range matches {

		if len(currentPath) != 0 {
			fmt.Printf("\n")
		}

		if match.Path != currentPath {
			currentPath = match.Path
		}

		if match.NameMatch {
			printMatchedPath(match)
			Ansi(Reset)
			fmt.Println("(filename)")
			Ansi(Reset)
		} else {
			Ansi(Green, Italic)
			fmt.Printf("%s ", match.Path)
			Ansi(Reset)
			Ansi(Yellow)
			fmt.Printf("%d:%d\n", match.Row, match.Column)
			Ansi(Reset)
			printMatchedLine(match)
			//fmt.Printf("%s\n", match.Line)
		}

	}
}

func printMatchedLine(match Match) {
	col := match.Column - 1
	if col < 0 {
		col = 0
	}

	left := match.Line[0:col]
	right := match.Line[col+match.Length : len(match.Line)]

	Ansi(Reset)
	fmt.Print(left)
	Ansi(Red, Bold)
	fmt.Print(match.Matched)
	Ansi(Reset)
	fmt.Println(right)
}

func printMatchedPath(match Match) {
	left := match.Path[0:match.Column]
	right := match.Path[match.Column+match.Length : len(match.Path)]

	Ansi(Reset)
	fmt.Print(left)
	Ansi(Red, Bold)
	fmt.Print(match.Matched)
	Ansi(Reset)
	fmt.Println(right)
}

func handleFlags() (*Configuration, bool) {
	flag.Usage = printHelp

	verboseFlag := flag.Bool("verbose", false, "Verbose")
	versionFlag := flag.Bool("v", false, "Show version")
	noAnsiFlag := flag.Bool("n", false, "No colors")
	depthCapFlag := flag.Int("dc", 10, "Maximum directory depth")
	sizeCapFlag := flag.Int64("fs", 20_000, "Maximum file size in kilobytes") // 20 MB by default
	countCapFlag := flag.Int("fc", 100_000, "Maximum file count")
	matchCapFlag := flag.Int("mc", 1_000, "Maximum matches per file")

	flag.Parse()

	setAnsiEnabled(!*noAnsiFlag)

	if *versionFlag {
		printVersion()
		return nil, false
	}

	if *depthCapFlag < 0 {
		Ansi(Red)
		fmt.Println("Directory depth cap needs to be bigger or equal to zero")
		Ansi(Reset)
		return nil, true
	}

	if *sizeCapFlag <= 0 {
		Ansi(Reset)
		fmt.Println("File size cap needs to be bigger than zero")
		Ansi(Reset)
		return nil, true
	}

	if *countCapFlag <= 0 {
		Ansi(Red)
		fmt.Println("File count cap needs to be bigger than zero")
		Ansi(Reset)
		return nil, true
	}

	if *matchCapFlag <= 0 {
		Ansi(Red)
		fmt.Println("Match cap needs to be bigger than zero")
		Ansi(Reset)
		return nil, true
	}

	args := flag.Args()
	argCount := len(args)

	if argCount == 0 {
		printHelp()
		return nil, false
	}

	if argCount > 2 {
		fmt.Println(Red, "\bInvalid argument count", Reset)
		return nil, true
	}

	pattern, err := regexp.Compile(args[0])

	if err != nil {
		fmt.Println(Red, "\bCould not compile pattern:", err, Reset)
		return nil, true
	}

	path := "."

	if argCount > 1 {
		path = args[1]
	}

	config := Configuration{
		Verbose:  *verboseFlag,
		Path:     path,
		Pattern:  pattern,
		DepthCap: *depthCapFlag,
		SizeCap:  *sizeCapFlag * 1_000, // converts KB to bytes
		CountCap: *countCapFlag,
		MatchCap: *matchCapFlag,
	}

	return &config, false
}

func printHelp() {
	fmt.Printf("Usage: %s [flags...] <pattern> <path>\n", os.Args[0])
	fmt.Println("Flags:")
	flag.PrintDefaults()
	fmt.Println("  <pattern>\n  \tPattern to search for")
	fmt.Println("  <path>\n  \tPath to directory or file")
}

func printVersion() {
	Ansi(Green)
	fmt.Printf("Rontgen version %s\n", Version)
	Ansi(Reset)
}
