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
	ResultName     string
	ResultType     string
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

func getResults(q Query, database) QueryResult {

}

func insertElement(resultName string, resultType string, resultLanguage string) {
	stmt, err := tx.Prepare("INSERT INTO TABLE values(?, ?)")

	if err != nil {
		log.Fatal(err);
	}
}

func createDatabase() {
	stmt, err := tx.Prepare("CREATE TABLE searchIndex(id INTEGER PRIMARY KEY, name TEXT, type TEXT, path TEXT)");
}

/**
 * connect to an sql database
 */
func connect() {
	db, err := sql.Open("sqlite3", "./database.db")

	if err != nil {
		log.Fatal(err)
	}

}
