package core

/**
* This is the core layer.
*
* This class should interface with the database to get results based on a provided query.
* Results should be returned as if a fuzzy search is performed.
 */

import (
	"errors"
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

var NoLanguageError error = errors.New("core: no language specified")

// TODO: save language type and make it mutable while program is running
// Return a list of documents that match the given query.
func Query(query string, language string) ([]data_models.DocsetElement, error) {
	queryEngine := docset_logic.GetQueryEngine()

	var docset data_models.Docset

	if language == "" {
		return []data_models.DocsetElement{}, NoLanguageError
	}

	docset = queryEngine.GetIndicesForLanguage(language)
	filterResults := docset.Filter(query)
	return filterResults, nil
}

// Return the html-string of the document body.
func GetDocContent(doc data_models.DocsetElement, language string) (string, error) {
	documentationLocation := doc.Path
	var cleanedLocation string

	hashIndex := strings.LastIndex(documentationLocation, "#")
	if hashIndex > 0 {
		cleanedLocation = documentationLocation[0:hashIndex]
	} else {
		cleanedLocation = documentationLocation
	}

	queryEngine := docset_logic.GetQueryEngine()
	documentationData := queryEngine.LoadDocumentationData(language, cleanedLocation)
	return string(documentationData), nil
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
