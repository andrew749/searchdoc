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
				data.Print()
			}
		} // switch
	} // for processing file entries
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
		GetDocsetPath(languageName),
		"/Contents/Resources/docSet.dsidx")

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
