package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"reflect"
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
 * Query the database for a specific result, will automatically fill type
 */
func getAllResults(t reflect.Type, databaseName string, db *sql.DB) []interface{} {

	res := reflect.MakeSlice(reflect.SliceOf(t), 0, 10)

	query := fmt.Sprint("SELECT * from %s", databaseName)

	rows, err := db.Query(query)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	// get the columns of the results
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {

		for i, _ := range columns {
			valuePtrs[i] = &values[i]
		}

		err := rows.Scan(valuePtrs...)

		if err != nil {
			log.Fatal(err)
		}

		element := reflect.New(t)

		res = reflect.Append(res, element)
	}
	return res.Interface().([]*t)
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
