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
	DocsetDir       = "docsets"
)

func PreferencesDirectory() string {
	path := filepath.Join(os.Getenv("HOME"), ".searchdoc")

	return path
}

func PreferencesPath() string {
	path := filepath.Join(PreferencesDirectory(), PreferencesFile)

	return path
}

type Preferences struct {
	DocsetPath string `json:"docset_path"`
}

func loadPreferences() Preferences {
	var data Preferences

	var file *os.File
	if _, err := os.Stat(PreferencesPath()); os.IsNotExist(err) {

		if _, err = os.Stat(PreferencesDirectory()); os.IsNotExist(err) {
			log.Println("Trying to create directory")

			err2 := os.Mkdir(PreferencesDirectory(), 0777)

			if err2 != nil {
				log.Fatal(err2)
			}

		}

		log.Println("Trying to create file")
		file, err2 := os.Create(PreferencesPath())

		if err2 != nil {
			log.Println("Cant create")
			log.Fatal(err2)
		}

		docsetPath := filepath.Join(os.Getenv("HOME"), DocsetDir)
		defaultSettings := fmt.Sprintf("{\"docset_path\": \"%s\"}", docsetPath)

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
