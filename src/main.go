package main

import (
	"flag"
	"fmt"
	"searchdoc/src/docset_logic"
)

var language string
var queryType string

func processCommand(query string, language string) {
	queryEngine := docset_logic.DocsetQueryEngineImpl{}

	docset := queryEngine.GetIndicesForLanguage(language)
	fmt.Printf(docset.Name)
	for _, x := range docset.Data {
		x.PrintElement()
	}

	// get the feed data
	// To be used for determining what to download
	//feeds := GetDocsetFeeds()

	// downloads work
	//DownloadDocset(feeds[0].Urls[0])
	//fmt.Println(GetAvailableDocsets())
}

func main() {

	// get the arguments
	var (
		query    string
		language string
	)

	flag.StringVar(&query, "query", "", "The query to search")
	flag.StringVar(&language, "language", "", "The query to search")
	flag.Parse()

	fmt.Printf("language: %s\nquery: %s\n", language, query)

	// process the command
	// TODO: replace with connection to ui
	processCommand(query, language)

}
