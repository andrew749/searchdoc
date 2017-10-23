package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

const PreferencesPath = "~/.searchdoc/preferences.json"

type Preferences struct {
	DatabasePath string `json:"database_path"`
}

func loadPreferences() Preferences {
	var data Preferences

	file, err := os.Open(PreferencesPath)

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

func getPreferences() Preferences {
	return _Preferences
}
