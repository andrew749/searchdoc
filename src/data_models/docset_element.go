package data_models

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
