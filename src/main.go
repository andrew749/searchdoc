package main

import (
	"flag"
	"fmt"
	"log"
	"searchdoc/src/data_models"
	docset_logic "searchdoc/src/docset_logic"
	"searchdoc/src/ui"
	"strings"
)

var language string
var queryType string

func processCommand(query string, language string) {
	queryEngine := docset_logic.GetQueryEngine()

	var docset data_models.Docset

	if language != "" {
		docset = queryEngine.GetIndicesForLanguage(language)
		fmt.Printf(docset.Name)
	} else {
		fmt.Printf("TODO: Not implemented all langauge search.\n")
		return
	}

	filterResults := docset.Filter(query)
	count := 0
	for _, x := range filterResults {
		fmt.Printf("%d) %s\n", count, x.Name)
		count += 1
	}

	if count == 0 {
		fmt.Printf("No results found.\n")
		return
	}

	var selection = 0
	_, err := fmt.Scanf("%d", &selection)

	if err != nil {
		log.Fatal(err)
	}

	documentationLocation := filterResults[selection].Path
	// remove any trailing #
	cleanedLocation := documentationLocation[:strings.LastIndex(documentationLocation, "#")]
	documentationData := queryEngine.LoadDocumentationData(language, cleanedLocation)
    ui.Init()
    ui.SetContent(documentationData)
	//fmt.Println(string(documentationData))

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
