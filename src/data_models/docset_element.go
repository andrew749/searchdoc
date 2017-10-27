package main

/**
* Each row of a docset index from the sqlite file.
 */
type DocsetElement struct {
	Id   int
	Name string
	Type string
	Path string
}
