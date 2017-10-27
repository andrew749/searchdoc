package main

import (
	"fmt"
	"os"
)

var language string
var queryType string

func printUsage(programName string) {
	fmt.Printf("./%s [--query QUERY] [--language] [--type]* \n", programName)
	os.Exit(1)
}

// process a specific query based on the
func processCommand(query string, language string, queryType string) {
	queryEngine := DocsetQueryEngineImpl{}

	queryEngine.GetIndicesForLanguage(language)

	// get the plist data
	GetDocsetPList("Go")

	// get the feed data
	GetDocsetFeeds()
}

func main() {

	// get the arguments

	if len(os.Args) < 2 {
		// not enough arguments
		printUsage(os.Args[0])
	}

	query := os.Args[1]

	var language string

	if len(os.Args) >= 3 {
		language = os.Args[2]
	}

	var queryType string

	if len(os.Args) >= 4 {
		queryType = os.Args[3]
	}

	// process the command
	processCommand(query, language, queryType)

}
