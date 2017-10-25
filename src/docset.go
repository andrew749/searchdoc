package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	plist "howett.net/plist"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

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
func GetDocsetPList(fileData []byte) DocsetPlist {
	var data DocsetPlist

	_, err := plist.Unmarshal(fileData, data)

	if err != nil {
		log.Fatal(err)
	}

	return data
}

func GetDocsetFeeds() []string {
	repoUrl := `https://github.com/Kapeli/feeds/tar/master.tar.gz`
	resp, err := http.Get(repoUrl)

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	reader := bytes.NewReader(body)
	//TODO: finish getting docset feeds
	_, err = gzip.NewReader(reader)

	if err != nil {
		log.Fatal(err)
	}

	return []string{"TODO REMOVE"}
}

/**
* Downloads a docset and places it into a specific directory as a tar.gz file
 */
func DownloadDocset(docsetName string, saveDirectory string) {
	baseString := "https://github.com/Kapeli/feeds/%s.tar.gz"
	docsetString := fmt.Sprint(baseString, docsetName)
	_, err := http.Get(docsetString)

	if err != nil {
		log.Fatal(err)
	}

}

/**
* Provided with a docset name, read the sqlite index and populate a docset object.
 */
func LoadSQLiteIndex(languageName string) Docset {
	databasePath := filepath.Join(
		GetPreferences().DocsetPath,
		fmt.Sprintf("%s.docset/Contents/Resources/docSet.dsidx", languageName))

	query := DocsetQuery{databasePath}

	languageResults := GetAllIndexResultsForLanguage(query)

	var docset Docset
	docset.Data = make([]DocsetElement, 0, 0)

	// iterate over the results for a specific language
	for _, element := range languageResults {
		var tempElement DocsetElement = DocsetElement{
			element.Id,
			element.QueryResultName,
			element.QueryResultType,
			element.QueryResultPath,
		}

		docset.Data = append(docset.Data, tempElement)
	}

	return docset
}
