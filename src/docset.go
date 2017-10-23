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
* Each row of a docset index
 */
type DocsetElement struct {
	Id   int
	Name string
	Type string
	Path string
}

type DocsetData struct {
	Elements []DocsetElement
}

type Docset struct {
	Name string
	Path string
	Data DocsetData
}

/**
* Hold info.plist information from the docset
 */
type DashPlist struct {
	Identifier           string `plist:"CFBundleIdentifier"`
	Name                 string `plist:"CFBundleName"`
	DocsetPlatformFamily string `plist:"DocSetPlatformFamily"`
	isDashDocset         bool   `plist:"isDashDocset"`
}

func getDocsetPList(fileData []byte) DashPlist {
	var data DashPlist
	_, err := plist.Unmarshal(fileData, data)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

/**
* Downloads a docset and places it into a specific directory as a tar.gz file
 */
func downloadDocset(docsetName string, saveDirectory string) {
	baseString := "https://github.com/Kapeli/Dash-User-Contributions/tree/master/docsets/%s.tar.gz"
	docsetString := fmt.Sprint(baseString, docsetName)
	_, err := http.Get(docsetString)

	if err != nil {
		log.Fatal(err)
	}

}

/**
* Provided with a docset name read the sqlite index and populate the local db.
 */
func loadSQLiteIndex(languageName string, docsetPath string) []DocsetData {
	docsetIndexPath := "%s/Contents/Resources/docSet.dsidx"
	database := OpenDatabaseFile(docsetIndexPath)

	defer database.Close()

	for _, result := range getAllResults(database) {
		insertSearchIndexElement(result.ResultName, result.ResultType, languageName, result.ResultPath, database)
	}

	return nil
}
