package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"path/filepath"
	"searchdoc/src/utils"
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

	db := OpenDatabaseFile(query.Path)

	defer db.Close()

	queryResults, err := db.Query("SELECT id, name, type, path FROM searchIndex")

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
	sqLitePath := filepath.Join(
		utils.GetDocsetPath(language),
		"/Contents/Resources/docSet.dsidx")
	return sqLitePath
}

/**
* Open an sqlite database file
 */
func OpenDatabaseFile(filePath string) *sql.DB {
	db, err := sql.Open("sqlite3", filePath)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
