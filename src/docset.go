package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io"
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

type FeedData struct {
	Name     string
	Version  string   `xml:"version"`
	Urls     []string `xml:"url"`
	Versions []string `xml:"other-versions>version>name"`
}

func (data *FeedData) Print() {
	fmt.Println(data.Name)
	fmt.Println(data.Version)
	fmt.Println(data.Urls)
	fmt.Print("\n\n")
}

/**
* Connect to github and get the latest feeds from kapeli's repo.
 */
func GetDocsetFeeds() []FeedData {
	repoUrl := `https://github.com/Kapeli/feeds/archive/master.tar.gz`
	resp, err := http.Get(repoUrl)

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	byteReader := bytes.NewReader(body)

	gzipReader, err := gzip.NewReader(byteReader)

	if err != nil {
		log.Fatal(err)
	}

	res := make([]FeedData, 0)

	reader := tar.NewReader(gzipReader)

	for {

		header, err := reader.Next()

		// check if the error is expected
		if err == io.EOF {
			return res
		} else if err != nil {
			log.Fatal(err)
		} // if

		if filepath.Ext(header.FileInfo().Name()) != ".xml" {
			continue
		}

		switch header.Typeflag {
		// process files
		case tar.TypeReg:
			{
				var data FeedData
				buf := new(bytes.Buffer)
				buf.ReadFrom(reader)
				xml.Unmarshal(buf.Bytes(), &data)
				data.Name = header.FileInfo().Name()

				res = append(res, data)
			}
		} // switch
	} // for processing file entries
}

/**
* Downloads a docset and places it into a specific directory as a tar.gz file
 */
func DownloadDocset(language string, url string) error {
	log.Printf("Downloading  %s", language)
	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	// add the newly downloaded docset to the appropriate directory
	pathToUnzip := DocsetPath()
	log.Printf("Saving to  %s", pathToUnzip)
	Untar(pathToUnzip, resp.Body)

	return nil
}

func GetAvailableDocsets() []string {
	docsetNames := make([]string, 0)

	docsetDirectories, err := ioutil.ReadDir(DocsetPath())

	if err != nil {
		log.Fatal(err)
	}

	for _, directory := range docsetDirectories {
		docsetNames = append(docsetNames, directory.Name())
	}

	return docsetNames
}

/**
* Provided with a docset name, read the sqlite index and populate a docset object.
 */
func LoadSQLiteIndex(languageName string) Docset {
	databasePath := filepath.Join(
		GetDocsetPath(languageName),
		"Contents",
		"Resources",
		"docSet.dsidx")

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
