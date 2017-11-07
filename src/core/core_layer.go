package core

/**
* This is the core layer.
*
* This class should interface with the database to get results based on a provided query.
* Results should be returned as if a fuzzy search is performed.
 */

import (
    "fmt"
    "log"
    "strings"
    "searchdoc/src/data_models"
    docset_logic "searchdoc/src/docset_logic"
)

// type SearchQuery struct {
// 	  QueryString string
//    Language string
// }

/**
 * TODO:
 * - make this a channel of search queries (look into this)
 * - use fuzzy search on input query string
 * - connect to UI
 */

type CoreLayer interface {
    // TODO: chan SearchQuery?
    Query(query string, language string) []string
    DownloadDocset()
    ListDownloadableDocsets() []string
    ListInstalledDocsets() []string
}

// TODO: save language type and make it mutable while program is running
func Query(query string, language string) []string {
    var result []string
    queryEngine := docset_logic.GetQueryEngine()

    var docset data_models.Docset

    if language != "" {
        docset = queryEngine.GetIndicesForLanguage(language)
        fmt.Printf(docset.Name)
    } else {
        fmt.Printf(docset.Name)
        return append(result, "No language specified.\n")
    }

    filterResults := docset.Filter(query)
    count := 0

    for _, x := range filterResults {
        fmt.Printf("%d) %s\n", count, x.Name)
        count += 1
    }

    if count == 0 {
        return append(result, "No results found.\n")
    }

    var selection = 0
    _, err := fmt.Scanf("%d", &selection)

    if err != nil {
        log.Fatal(err)
    }

    documentationLocation := filterResults[selection].Path
    var cleanedLocation string

    hashIndex := strings.LastIndex(documentationLocation, "#")
    if hashIndex > 0 {
        cleanedLocation = documentationLocation[0:hashIndex]
    } else {
        cleanedLocation = documentationLocation
    }

    documentationData := queryEngine.LoadDocumentationData(language, cleanedLocation)
    result = append(result, string(documentationData))

    return result
}

func DownloadDocset(language string) {
   _ = docset_logic.GetQueryEngine().DownloadDocset(language) 
}

func ListDownloadableDocsets() []string {
    return docset_logic.GetQueryEngine().GetDownloadableDocsets()
}

func ListInstalledDocsets() []string {
    return docset_logic.GetQueryEngine().GetDownloadedDocsets()
}
