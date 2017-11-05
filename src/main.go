package main

import (
	"flag"
	"fmt"
	docset_logic "searchdoc/src/docset_logic"
)

var language string
var queryType string

func processCommand(query string, language string) {
	queryEngine := docset_logic.GetQueryEngine()

	docset := queryEngine.GetIndicesForLanguage(language)
	fmt.Printf(docset.Name)

	for _, x := range docset.Filter(query) {
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

	download_list := flag.Bool("download_list", false, "Indicate if you want to list the downloadable packages.")
	installed_list := flag.Bool("list", false, "List all installed packages.")
	package_to_download := flag.String("download", "", "Download the specified package.")

	flag.Parse()

	if *download_list {
		for _, x := range docset_logic.GetQueryEngine().GetDownloadableDocsets() {
			fmt.Println(x)
		}
		return
	} else if *installed_list {
		for _, x := range docset_logic.GetQueryEngine().GetDownloadedDocsets() {
			fmt.Println(x)
		}
		return
	} else if *package_to_download != "" {
		_ = docset_logic.GetQueryEngine().DownloadDocset(*package_to_download)
		return

	}

	fmt.Printf("language: %s\nquery: %s\n", language, query)
	// process the command
	// TODO: replace with connection to ui
	processCommand(query, language)

}
