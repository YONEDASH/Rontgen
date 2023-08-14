package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
)

func main() {
	config := handleFlags()

	if config == nil {
		os.Exit(1)
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

		fmt.Printf("%s %d:%d\n", match.Path, match.Row, match.Column)
		fmt.Printf("%s\n", match.Line)
	}
}

func handleFlags() *Configuration {
	flag.Usage = printHelp

	verboseFlag := flag.Bool("verbose", false, "Verbose")
	versionFlag := flag.Bool("v", false, "Show version")

	flag.Parse()

	if *versionFlag {
		printVersion()
		return nil
	}

	args := flag.Args()
	argCount := len(args)

	if argCount == 0 {
		printHelp()
		return nil
	}

	if argCount > 2 {
		fmt.Println("Invalid argument count")
		return nil
	}

	pattern, err := regexp.Compile(args[0])

	if err != nil {
		fmt.Println("Could not compile pattern:", err)
		return nil
	}

	path := "."

	if argCount > 1 {
		path = args[1]
	}

	config := Configuration{
		Verbose: *verboseFlag,
		Path:    path,
		Pattern: pattern,
	}

	return &config
}

func printHelp() {
	fmt.Printf("Usage: %s [-v] [-verbose] <path> <pattern>\n", os.Args[0])
	fmt.Println("Flags:")
	flag.PrintDefaults()
	fmt.Println("  <path>\tPath to directory or file")
	fmt.Println("  <pattern>\tPattern to search for")
}

func printVersion() {
	fmt.Printf("Rontgen v%s\n", Version)
}
