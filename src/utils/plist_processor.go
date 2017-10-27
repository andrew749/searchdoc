package main

import (
	plist "howett.net/plist"
	"io/ioutil"
	"path/filepath"
)

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
func GetDocsetPList(docset string) (DocsetPlist, error) {
	var data DocsetPlist

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
