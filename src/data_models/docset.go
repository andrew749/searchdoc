package data_models

/**
* Data class to hold a conceptual docset to be passed throughout the applciation
 */

type Docset struct {
	Name        string
	Path        string
	DocsetPlist DocsetPlist
	Data        []DocsetElement
}
