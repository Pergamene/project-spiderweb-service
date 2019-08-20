// This test file sets up Main so that it:
// 1. Only runs this packages' tests if it can establish a connection to the local database container and pull in the setup.sql file.
// 2. It will create a new database under the root db user, and execute the setup.sql file
// 3. Specific stores' tests should assert that mysqldb is setup and ready to pass to the store.
// 4. Once the tests run, it will close and remove the temporarly database.
// It is the responsibility of the individual tests to reset the tables to a testable state before
// runnning their tests: you can only assume that the setup.sql command added the neccesary tables
// and that the tables likely contain junk content that needs to be deleted.
// The helper functions `clearTableForTest` and `execPreTestQueries` can be used by the tests to
// prepare the test before execution.
package mysqlstore

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Pergamene/project-spiderweb-service/internal/util/env"
)

var mysqldb *sql.DB
var mysqldbName string

func getDb() (*sql.DB, string, bool, error) {
	filePath := env.Get("SETUP_SQL_FILEPATH", "/Users/rhyeen/Documents/repos/project-spiderweb/project-spiderweb-db/setup.sql")
	if filePath == "" {
		return nil, "", false, nil
	}
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, "", false, err
	}
	fileString := string(fileBytes)
	queries := getQueries(fileString)
	newDB, newDBName, err := createAndOpenNewDB()
	if err != nil {
		return nil, "", false, err
	}
	err = executeQueries(newDB, queries)
	if err != nil {
		return newDB, newDBName, false, err
	}
	return newDB, newDBName, true, nil
}

func getQueries(fileString string) []string {
	return strings.Split(string(fileString), ";\n")
}

func createAndOpenNewDB() (*sql.DB, string, error) {
	newDBName := getRandomDBName()
	rootDB, err := SetupRootMySQL("")
	if err != nil {
		return nil, newDBName, err
	}
	defer rootDB.Close()
	_, err = rootDB.Exec("CREATE DATABASE IF NOT EXISTS " + newDBName)
	if err != nil {
		return nil, newDBName, err
	}
	rootDB.Close()

	db, err := SetupRootMySQL(newDBName)
	if err != nil {
		return nil, newDBName, err
	}
	return db, newDBName, nil
}

func getRandomDBName() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func executeQueries(db *sql.DB, queries []string) error {
	if db == nil {
		return nil
	}
	if queries == nil {
		return nil
	}
	for _, query := range queries {
		if query == "" {
			continue
		}
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func closeAndRemoveDb(db *sql.DB, dbName string) error {
	db.Close()
	rootDB, err := SetupRootMySQL("")
	if err != nil {
		return err
	}
	defer rootDB.Close()
	_, err = rootDB.Exec("DROP DATABASE IF EXISTS " + dbName)
	if err != nil {
		return err
	}
	return nil
}

func TestMain(m *testing.M) {
	db, dbName, isTestible, err := getDb()
	mysqldb = db
	mysqldbName = dbName
	fmt.Printf("Initialized db: %v\n", mysqldbName)
	if err != nil {
		if mysqldb != nil {
			closeAndRemoveDb(mysqldb, dbName)
		}
		fmt.Printf("Unable to bootstrap DB:\n%v", err)
		os.Exit(1)
	}
	if !isTestible {
		fmt.Printf("Not configured to run mysqlstore tests")
		os.Exit(0)
	}
	result := m.Run()
	if mysqldb != nil {
		err := closeAndRemoveDb(mysqldb, dbName)
		if err != nil {
			fmt.Printf("Unable to close and remove DB %v:\n%v", dbName, err)
			os.Exit(1)
		}
	}
	os.Exit(result)
}

func clearTableForTest(db *sql.DB, table string) error {
	if db == nil {
		return nil
	}
	statement, err := db.Prepare(fmt.Sprintf("TRUNCATE TABLE `%v`", table))
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec()
	return err
}

func execPreTestQueries(db *sql.DB, queries []string) error {
	return executeQueries(db, queries)
}
