package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	PreferencesFile = "preferences.json"
	PreferencesDir  = ".searchdoc"
	DocsetDir       = "docsets"
)

func PreferencesDirectory() string {
	return filepath.Join(os.Getenv("HOME"), PreferencesDir)
}

func PreferencesPath() string {
	return filepath.Join(PreferencesDirectory(), PreferencesFile)
}

func DocsetPath() string {
	return filepath.Join(PreferencesDirectory(), DocsetDir)
}

type Preferences struct {
	DocsetPath    string `json:"docset_path"`
	SearchDocPath string `json:"search_doc_path"`
}

func loadPreferences() Preferences {
	var (
		data Preferences
		file *os.File
	)

	if _, err := os.Stat(PreferencesPath()); os.IsNotExist(err) {

		if _, err = os.Stat(PreferencesDirectory()); os.IsNotExist(err) {

			log.Println("Trying to create directory")
			err = os.Mkdir(PreferencesDirectory(), 0777)

			if err != nil {
				log.Fatal(err)
			}

		}

		log.Println("Trying to create file")
		file, err := os.Create(PreferencesPath())

		if err != nil {
			log.Println("Cant create")
			log.Fatal(err)
		}

		defaultSettings := fmt.Sprintf(`{
			"docset_path": "%s",
			"search_doc_path":"%s"
		}`, DocsetPath(), PreferencesDirectory())

		file.WriteString(defaultSettings)

		log.Println("Created file")
	}

	file, err := os.Open(PreferencesPath())

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(fileBytes, &data)

	if err != nil {
		log.Fatal(err)
	}

	return data
}

var _Preferences = loadPreferences()

func GetPreferences() Preferences {
	return _Preferences
}
