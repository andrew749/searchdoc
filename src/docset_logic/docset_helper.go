package docset_logic

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	data_models "searchdoc/src/data_models"
	database "searchdoc/src/database"
	utils "searchdoc/src/utils"
	"strings"
)

/**
* Connect to github and get the latest feeds from kapeli's repo.
 */
func GetDocsetFeeds() []data_models.FeedData {
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

	res := make([]data_models.FeedData, 0)

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
				var data data_models.FeedData
				buf := new(bytes.Buffer)
				buf.ReadFrom(reader)
				xml.Unmarshal(buf.Bytes(), &data)
				data.Name = strings.TrimSuffix(header.FileInfo().Name(), filepath.Ext(header.FileInfo().Name()))

				res = append(res, data)
			}
		} // switch
	} // for processing file entries
}

/**
* Downloads a docset and places it into a specific directory as a tar.gz file
* specify the language name and the url to query from
 */
func DownloadDocset(url string) error {
	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	// add the newly downloaded docset to the appropriate directory
	pathToUnzip := utils.DocsetPath()
	log.Printf("Saving to  %s", pathToUnzip)
	utils.Untar(pathToUnzip, resp.Body)

	return nil
}

/**
* Get the available docsets which are stored locally in the
* directory specified in preferences.
 */
func GetAvailableDocsets() []string {

	docsetNames := make([]string, 0)

	docsetDirectories, err := ioutil.ReadDir(utils.DocsetPath())

	if err != nil {
		log.Fatal(err)
	}

	for _, directory := range docsetDirectories {

		docsetPath := filepath.Join(utils.DocsetPath(), directory.Name())

		dashDir := false
		// if the check the docsetpath
		err = filepath.Walk(docsetPath, func(path string, info os.FileInfo, err error) error {
			// check that a plist exists
			if info.Name() == "docSet.dsidx" {
				// found a dash docset
				dashDir = true
			}

			return nil
		})

		if dashDir {
			docsetNames = append(docsetNames, directory.Name())
		}
	}

	return docsetNames
}

/**
* Provided with a docset name, read the sqlite index and populate a Docset object.
 */
func DocsetForLanguage(language string) data_models.Docset {
	databasePath := filepath.Join(
		utils.GetDocsetPath(language),
		"Contents",
		"Resources",
		"docSet.dsidx")

	query := database.DocsetQuery{databasePath}

	languageResults := database.GetAllIndexResultsForLanguage(query)

	var docset data_models.Docset
	plist, err := utils.GetDocsetPList(language)

	if err != nil {
		log.Fatal(err)
	}

	docset.DocsetPlist = plist
	docset.Data = make([]data_models.DocsetElement, 0, 0)

	// iterate over the results for a specific language
	for _, element := range languageResults {
		var tempElement data_models.DocsetElement = data_models.DocsetElement{
			element.Id,
			element.QueryResultName,
			element.QueryResultType,
			element.QueryResultPath,
		}

		docset.Data = append(docset.Data, tempElement)
	}

	return docset
}
