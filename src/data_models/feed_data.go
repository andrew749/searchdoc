package data_models

/**
* A data class to hold data returned from the language feed from
* https://github.com/Kapeli/feeds
 */

import (
	"fmt"
)

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
