package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

/**
 * Result from a query.
 *
 * OUTPUT
 */
type DocsetQueryResult struct {
	Id              int
	QueryResultName string
	QueryResultType string
	QueryResultPath string
}

/**
 * Encapsulate a user query into the database.
 *
 * INPUT
 */
type DocsetQuery struct {
	// provide the database path
	Path string
}

func GetAllIndexResultsForLanguage(query DocsetQuery) []DocsetQueryResult {
	databaseLocation := GetSQLiteLocation(query.Path)

	db := OpenDatabaseFile(databaseLocation)
	defer db.Close()

	queryResults, err := db.Query("SELECT * FROM searchIndex")

	if err != nil {
		log.Fatal(err)
	}

	defer queryResults.Close()

	// hold results from query
	res := make([]DocsetQueryResult, 0, 1)

	for queryResults.Next() {
		var queryResult DocsetQueryResult

		err := queryResults.Scan(
			&queryResult.Id,
			&queryResult.QueryResultName,
			&queryResult.QueryResultType,
			&queryResult.QueryResultPath)
		if err != nil {
			log.Fatal(err)
		}
		res = append(res, queryResult)
	}

	return res
}

/**
* Get the location of an sqlite file.
 */
func GetSQLiteLocation(language string) string {
	return "ERROR NOT IMPLEMENTED"
}

/**
 * Open an sqlite database file
 */
func OpenDatabaseFile(filePath string) *sql.DB {
	db, err := sql.Open("sqlite3", filePath)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to database successfully")

	return db
}
