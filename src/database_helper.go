package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

/**
 * Result from a query into the db storing language specific results.
 */
type QueryResult struct {
	id             int
	ResultName     string
	ResultType     string
	ResultPath     string
	ResultLanguage string
}

/**
 * Encapsulate a user query
 */
type Query struct {
	QueryString string
	Language    string
	Type        string
}

/**
 * Query the database for a specific result
 */
func getAllResults(db *sql.DB) []QueryResult {
	rows, err := db.Query("SELECT * from searchIndex")
	res := make([]QueryResult, 0)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	// a result row
	var (
		id          int
		name        string
		elementType string
		path        string
		language    string
	)

	for rows.Next() {
		err := rows.Scan(&id, &name, &language, &elementType)

		if err != nil {
			log.Fatal(err)
		}

		res = append(res, QueryResult{id, name, elementType, language, path})
	}

	return res
}

func insertSearchIndexElement(name string, elementType string, language string, path string, db *sql.DB) {

	stmt, err := db.Prepare("INSERT INTO SearchIndex values(?, ?, ?, ?)")

	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(name, elementType, path, language)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Inserted index successfully")
}

func createSearchIndexTable(language, db *sql.DB) {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS SearchIndex(id INTEGER PRIMARY KEY AUTO_INCREMENT, name TEXT, type TEXT, path TEXT, language TEXT)")

	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Created table successfully")
}

func OpenApplicationDatabase() *sql.DB {
	return OpenDatabaseFile(getPreferences().DatabasePath)
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
