package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/Pergamene/project-spiderweb-service/internal/api"
	"github.com/Pergamene/project-spiderweb-service/internal/api/handlers"
	"github.com/Pergamene/project-spiderweb-service/internal/services/healthcheckservice"
	"github.com/Pergamene/project-spiderweb-service/internal/services/pageservice"
	"github.com/Pergamene/project-spiderweb-service/internal/stores/mysqlstore"
	"github.com/Pergamene/project-spiderweb-service/internal/util/env"
	"github.com/rs/cors"
)

const (
	defaultMySQLHost     = "127.0.0.1:3306"
	defaultMySQLProtocol = "tcp"
	defaultMySQLDatabase = "spiderweb_dev"
	defaultMySQLUser     = "spiderweb_dev"
	defaultMySQLPassword = "password"
	defaultMySQLCharset  = "utf8"
)

const localUIURL = "http://127.0.0.1:8781/"

const (
	defaultAdminAuthSecret = "DEFAULT_SECRET"
	defaultPort            = "8782"
	defaultStaticPath      = "../../static"
	defaultDatacenter      = "LOCAL"
)

// Indexer sets up indices for the appropriate data store
type Indexer interface {
	EnsureIndices() error
}

func getHTTPServerAddr() string {
	port := env.Get("PORT", defaultPort)
	return ":" + port
}

func getHTTPServerReadTimeout() time.Duration {
	return 10 * time.Second
}

func getHTTPServerWriteTimeout() time.Duration {
	return 10 * time.Second
}

func getHTTPServerMaxHeaderBytes() int {
	return 1 << 20
}

func getAPIPath() string {
	return "api"
}

func getStaticPath() string {
	return env.Get("STATIC_PATH", defaultStaticPath)
}

func getDatacenter() string {
	return env.Get("DATACENTER", defaultDatacenter)
}

func main() {
	mysqldb, err := setupMySQL()
	if err != nil {
		fmt.Printf("Failed to connect to MySQL db.\nIf connecting locally, follow https://github.com/Pergamene/project-spiderweb-db/blob/master/README.md to get the local db running.\n")
		log.Fatal(err)
	}
	defer mysqldb.Close()
	apiPath := getAPIPath()
	staticPath := getStaticPath()
	datacenter := getDatacenter()
	handler, err := setupHandler(apiPath, staticPath, datacenter, mysqldb)
	if err != nil {
		log.Fatal(err)
	}
	handler, err = setupCors(datacenter, handler)
	if err != nil {
		log.Fatal(err)
	}
	s := &http.Server{
		Addr:           getHTTPServerAddr(),
		Handler:        handler,
		ReadTimeout:    getHTTPServerReadTimeout(),
		WriteTimeout:   getHTTPServerWriteTimeout(),
		MaxHeaderBytes: getHTTPServerMaxHeaderBytes(),
	}
	fmt.Printf("Starting server at http://localhost%v\nVerify locally by running:\ncurl -X GET http://localhost%v/%v/healthcheck\n", getHTTPServerAddr(), getHTTPServerAddr(), getAPIPath())
	log.Fatal(s.ListenAndServe())
}

func setupMySQL() (*sql.DB, error) {
	dsnFormat := fmt.Sprintf("%v:%v@%v(%v)/%v?charset=%v",
		getMySQLUser(),
		getMySQLPassword(),
		getMySQLProtocol(),
		getMySQLHost(),
		getMySQLDatabase(),
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

func setupHandler(apiPath, staticPath, datacenter string, mysqldb *sql.DB) (http.Handler, error) {
	var handler http.Handler
	pageStore := mysqlstore.NewPageStore(mysqldb)
	healthcheckStore := mysqlstore.NewHealthcheckStore(mysqldb)
	pageService := pageservice.PageService{
		PageStore: pageStore,
	}
	healthcheckService := healthcheckservice.HealthcheckService{
		HealthcheckStore: healthcheckStore,
	}
	routerHandlers := api.RouterHandlers{
		PageHandler: handlers.PageHandler{
			PageService: pageService,
		},
		HealthcheckHandler: handlers.HealthcheckHandler{
			HealthcheckService: healthcheckService,
		},
	}
	router := api.NewRouter(apiPath, staticPath, routerHandlers)
	authN, authZ, err := getAuths(apiPath, datacenter)
	if err != nil {
		return handler, err
	}
	return &api.Handler{
		AuthN:      authN,
		AuthZ:      authZ,
		Router:     router,
		Datacenter: datacenter,
		APIPath:    apiPath,
	}, nil
}

func getAuths(apiPath, datacenter string) (api.AuthN, api.AuthZ, error) {
	adminAuthSecret, err := getAdminAuthSecret(datacenter)
	if err != nil {
		return api.AuthN{}, api.AuthZ{}, err
	}
	authN := api.AuthN{
		Datacenter:      datacenter,
		AdminAuthSecret: adminAuthSecret,
	}
	authZ := api.AuthZ{
		APIPath: apiPath,
	}
	return authN, authZ, nil
}

func getAdminAuthSecret(datacenter string) (string, error) {
	if datacenter != api.LocalEnv {
		return env.Require("ADMIN_AUTH_SECRET")
	}
	return env.Get("ADMIN_AUTH_SECRET", defaultAdminAuthSecret), nil
}

func setupCors(datacenter string, handler http.Handler) (http.Handler, error) {
	if datacenter != api.LocalEnv {
		return handler, nil
	}
	c := cors.New(cors.Options{
		AllowedOrigins: []string{localUIURL},
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{"X-AUTH-TOKEN", "Content-Type"},
	})
	return c.Handler(handler), nil
}
