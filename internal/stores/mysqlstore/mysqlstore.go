package mysqlstore

import (
	"database/sql"
	"fmt"

	// used to import the "mysql" package
	_ "github.com/go-sql-driver/mysql"

	"github.com/Pergamene/project-spiderweb-service/internal/util/env"
)

const (
	defaultMySQLHost         = "127.0.0.1:3306"
	defaultMySQLProtocol     = "tcp"
	defaultMySQLDatabase     = "spiderweb_dev"
	defaultMySQLUser         = "spiderweb_dev"
	defaultMySQLPassword     = "password"
	defaultMySQLCharset      = "utf8"
	defaultMySQLRootUser     = "root"
	defaultMySQLRootPassword = "rootpassword"
)

// SetupRootMySQL returns a MySQL db with the credentials pulled from the env vars for the root user.
func SetupRootMySQL(database string) (*sql.DB, error) {
	if database == "" {
		database = getMySQLDatabase()
	}
	dsnFormat := fmt.Sprintf("%v:%v@%v(%v)/%v?charset=%v",
		getMySQLRootUser(),
		getMySQLRootPassword(),
		getMySQLProtocol(),
		getMySQLHost(),
		database,
		getMySQLCharset())
	return setupMySQL(dsnFormat)
}

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
	return setupMySQL(dsnFormat)
}

func setupMySQL(dsnFormat string) (*sql.DB, error) {
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

func getMySQLRootUser() string {
	return env.Get("MYSQL_USER", defaultMySQLRootUser)
}

func getMySQLRootPassword() string {
	return env.Get("MYSQL_PASSWORD", defaultMySQLRootPassword)
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
