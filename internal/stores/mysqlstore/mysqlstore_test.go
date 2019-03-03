package mysqlstore

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Pergamene/project-spiderweb-service/internal/util/env"
	"github.com/Pergamene/project-spiderweb-service/internal/util/testutils"
	"github.com/stretchr/testify/require"
)

var mysqldb *sql.DB

func getDb() (*sql.DB, string, bool, error) {
	var db *sql.DB
	filePath := env.Get("SETUP_SQL_FILEPATH", "")
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
	db = newDB
	err = executeQueries(db, queries)
	if err != nil {
		return nil, newDBName, false, err
	}
	return db, newDBName, true, nil
}

func getQueries(fileString string) []string {
	return strings.Split(string(fileString), ";\n")
}

func createAndOpenNewDB() (*sql.DB, string, error) {
	newDBName := getRandomDBName()
	rootDB, err := SetupMySQL("")
	if err != nil {
		return nil, newDBName, err
	}
	defer rootDB.Close()
	_, err = rootDB.Exec("CREATE DATABASE IF NOT EXISTS " + newDBName)
	if err != nil {
		return nil, newDBName, err
	}
	rootDB.Close()

	db, err := SetupMySQL(newDBName)
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
	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func closeAndRemoveDb(db *sql.DB, dbName string) error {
	db.Close()
	rootDB, err := SetupMySQL("")
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
	if !isTestible {
		fmt.Printf("Not configured to run mysqlstore tests")
		os.Exit(0)
	}
	if err != nil {
		fmt.Printf("Unable to bootstrap DB:\n%v", err)
		os.Exit(1)
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

func TestHealthcheckIsHealthy(t *testing.T) {
	cases := []struct {
		name                   string
		shouldReplaceDBWithNil bool
		preTestQueries         []string
		returnIsHealthy        bool
		returnErr              error
	}{
		{
			name:            "db healthy",
			preTestQueries:  []string{"INSERT INTO `healthcheck` (`status`) VALUES (\"ok\")"},
			returnIsHealthy: true,
		},
		{
			name:            "db not healthy",
			preTestQueries:  []string{"INSERT INTO `healthcheck` (`status`) VALUES (\"error\")"},
			returnIsHealthy: false,
		},
		{
			name:            "db not healthy because entry doesn't exist",
			preTestQueries:  []string{},
			returnIsHealthy: false,
		},
		{
			name:                   "db not setup",
			shouldReplaceDBWithNil: true,
			preTestQueries:         []string{"INSERT INTO `healthcheck` (`status`) VALUES (\"ok\")"},
			returnIsHealthy:        false,
			returnErr:              errors.New("failure"),
		},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			healthcheckStore := HealthcheckStore{
				db: mysqldb,
			}
			if tc.shouldReplaceDBWithNil {
				healthcheckStore.db = nil
			}
			err := clearDBForTest(healthcheckStore.db)
			require.NoError(t, err)
			isHealthy, err := healthcheckStore.IsHealthy()
			errExpected := testutils.TestErrorAgainstCase(t, err, tc.returnErr)
			if errExpected {
				return
			}
			require.Equal(t, isHealthy, tc.returnIsHealthy)
		})
	}
}

func clearDBForTest(db *sql.DB) error {
	statement, err := db.Prepare("DELETE FROM `healthcheck`")
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec()
	return err
}
