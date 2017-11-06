package utils

import (
	plist "howett.net/plist"
	"io/ioutil"
	"path/filepath"
	data_models "searchdoc/src/data_models"
)

/**
* Convert a file corresponding to a plist to the appropriate datatype
 */
func GetDocsetPList(docset string) (data_models.DocsetPlist, error) {
	var data data_models.DocsetPlist

	path := GetDocsetPath(docset)

	fileData, err := ioutil.ReadFile(filepath.Join(path, "Contents", "Info.plist"))

	if err != nil {
		return data, err
	}

	// get the data from the file
	_, err = plist.Unmarshal(fileData, &data)

	if err != nil {
		return data, err
	}

	return data, nil
}
