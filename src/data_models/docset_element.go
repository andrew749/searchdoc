package data_models

import (
	"fmt"
)

/**
* Each row of a docset index from the sqlite file
* standard in every kapeli repo.
 */
type DocsetElement struct {
	Id   int
	Name string
	Type string
	Path string
}

func (element DocsetElement) PrintElement() {
	fmt.Printf("%s: %s\n", element.Name, element.Type)
}
