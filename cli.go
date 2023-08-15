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
			PrintReset()
			fmt.Println("(filename)")
			PrintReset()
		} else {
			fmt.Print(Green, Italic)
			fmt.Printf("%s ", match.Path)
			PrintReset()
			fmt.Print(Yellow)
			fmt.Printf("%d:%d\n", match.Row, match.Column)
			PrintReset()
			printMatchedLine(match)
			//fmt.Printf("%s\n", match.Line)
		}

	}
}

func printMatchedLine(match Match) {
	left := match.Line[0 : match.Column-1]
	right := match.Line[match.Column+match.Length-1 : len(match.Line)]

	PrintReset()
	fmt.Print(left)
	fmt.Print(Red, Bold, match.Matched)
	PrintReset()
	fmt.Println(right)
}

func printMatchedPath(match Match) {
	left := match.Path[0:match.Column]
	right := match.Path[match.Column+match.Length : len(match.Path)]

	PrintReset()
	fmt.Print(left)
	fmt.Print(Red, Bold, match.Matched)
	PrintReset()
	fmt.Println(right)
}

func handleFlags() (*Configuration, bool) {
	flag.Usage = printHelp

	verboseFlag := flag.Bool("verbose", false, "Verbose")
	versionFlag := flag.Bool("v", false, "Show version")
	depthCapFlag := flag.Int("dc", 10, "Maximum directory depth")

	flag.Parse()

	if *versionFlag {
		printVersion()
		return nil, false
	}

	if *depthCapFlag <= 0 {
		fmt.Print(Red)
		fmt.Print("Directory depth cap needs to be bigger than zero")
		fmt.Println(Reset)
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
	}

	return &config, false
}

func printHelp() {
	fmt.Printf("Usage: %s [-v] [-verbose] <path> <pattern>\n", os.Args[0])
	fmt.Println("Flags:")
	flag.PrintDefaults()
	fmt.Println("  <path> Path to directory or file")
	fmt.Println("  <pattern> Pattern to search for")
}

func printVersion() {
	fmt.Print(Green)
	fmt.Printf("Rontgen version %s\n", Version)
	PrintReset()
}

func PrintReset() {
	fmt.Print(Reset)
}
