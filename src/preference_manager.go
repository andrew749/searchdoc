package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func PreferencesDirectory() string {
	path := filepath.Join(os.Getenv("HOME"), ".searchdoc")

	return path
}

const PreferencesFile = "preferences.json"

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

		file.WriteString("{}")

		log.Println("Created file")
	}

	log.Println("DONE")

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
