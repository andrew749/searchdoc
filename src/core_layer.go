package main

/**
* This is the core layer.
*
* This class should interface with the database to get results based on a provided query.
* Results should be returned as if a fuzzy search is performed.
 */

type SearchQuery struct {
	QueryString string
}

type CoreLayer interface {
	Query(string, chan SearchQuery)
}
