package main

import (
	"fmt"
	plist "howett.net/plist"
	"log"
	"net/http"
)

/**
* Docsets we are currently supporting.
 */
var Docsets = [...]string{"Pandoc"}

/**
* Each row of a docset index from the sqlite file.
 */
type DocsetElement struct {
	Id   int
	Name string
	Type string
	Path string
}

type Docset struct {
	Name        string
	Path        string
	DocsetPlist DocsetPlist
	Data        []DocsetElement
}

/**
* Hold info.plist information from the docset.
 */
type DocsetPlist struct {
	Identifier           string `plist:"CFBundleIdentifier"`
	Name                 string `plist:"CFBundleName"`
	DocsetPlatformFamily string `plist:"DocSetPlatformFamily"`
	isDashDocset         bool   `plist:"isDashDocset"`
}

/**
* Convert a file corresponding to a plist to the appropriate datatype
 */
func getDocsetPList(fileData []byte) DocsetPlist {
	var data DocsetPlist
	_, err := plist.Unmarshal(fileData, data)

	if err != nil {
		log.Fatal(err)
	}

	return data
}

/**
* Downloads a docset and places it into a specific directory as a tar.gz file
 */
func DownloadDocset(docsetName string, saveDirectory string) {
	baseString := "https://github.com/Kapeli/Dash-User-Contributions/tree/master/docsets/%s.tar.gz"
	docsetString := fmt.Sprint(baseString, docsetName)
	_, err := http.Get(docsetString)

	if err != nil {
		log.Fatal(err)
	}
}

/**
* Provided with a docset name, read the sqlite index and populate the local db.
 */
func LoadSQLiteIndex(languageName string, docsetPath string) Docset {
	docsetIndexPath := "%s/Contents/Resources/docSet.dsidx"
	database := OpenDatabaseFile(docsetIndexPath)

	defer database.Close()

	query := DocsetQuery{languageName}

	languageResults := GetAllIndexResultsForLanguage(query)

	var docset Docset
	docset.Data = make([]DocsetElement, 0, 0)

	for _, element := range languageResults {
		var tempElement DocsetElement = DocsetElement{
			element.Id,
			element.QueryResultName,
			element.QueryResultType,
			element.QueryResultPath}
		docset.Data = append(docset.Data, tempElement)
	}

	return docset
}

type DocsetQueryEngine interface {
	GetIndicesForLanguage(language string) Docset
}

/**
* Concrete implementation of Query engine
 */
type DocsetQueryEngineImpl struct {
}

func (engine *DocsetQueryEngineImpl) GetIndicesForLanguage(language string) Docset {
	databasePath := ""
	return LoadSQLiteIndex(language, databasePath)
}
