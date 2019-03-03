package mysqlstore

import (
	"database/sql"
	"fmt"

	"github.com/Pergamene/project-spiderweb-service/internal/util/env"
)

const (
	defaultMySQLHost     = "127.0.0.1:3306"
	defaultMySQLProtocol = "tcp"
	defaultMySQLDatabase = "spiderweb_dev"
	defaultMySQLUser     = "spiderweb_dev"
	defaultMySQLPassword = "password"
	defaultMySQLCharset  = "utf8"
)

// SetupMySQL returns a MySQL db with the credentials pulled from the env vars.
func SetupMySQL(database string) (*sql.DB, error) {
	if database == "" {
		database = getMySQLDatabase()
	}
	dsnFormat := fmt.Sprintf("%v:%v@%v(%v)/%v?charset=%v",
		getMySQLUser(),
		getMySQLPassword(),
		getMySQLProtocol(),
		getMySQLHost(),
		database,
		getMySQLCharset())
	// see: https://github.com/go-sql-driver/mysql/wiki/Examples#a-word-on-sqlopen
	db, err := sql.Open("mysql", dsnFormat)
	if err != nil {
		return db, err
	}
	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		return db, err
	}
	return db, nil
}

func getMySQLUser() string {
	return env.Get("MYSQL_USER", defaultMySQLUser)
}

func getMySQLPassword() string {
	return env.Get("MYSQL_PASSWORD", defaultMySQLPassword)
}

func getMySQLProtocol() string {
	return env.Get("MYSQL_PROTOCOL", defaultMySQLProtocol)
}

func getMySQLHost() string {
	return env.Get("MYSQL_HOST", defaultMySQLHost)
}

func getMySQLDatabase() string {
	return env.Get("MYSQL_DATABASE", defaultMySQLDatabase)
}

func getMySQLCharset() string {
	return env.Get("MYSQL_CHARSET", defaultMySQLCharset)
}
