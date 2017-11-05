package data_models

import (
	"strings"
)

/**
* Data class to hold a conceptual docset to be passed throughout the applciation
 */

type Docset struct {
	Name        string
	Path        string
	DocsetPlist DocsetPlist
	Data        []DocsetElement
}

func (docset *Docset) Filter(query string) []DocsetElement {
	res := make([]DocsetElement, 0)

	for _, v := range docset.Data {

		if strings.Contains(v.Name, query) {
			res = append(res, v)
		}
	}

	return res
}
