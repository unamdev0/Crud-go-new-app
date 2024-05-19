package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/unamdev0/go-crud-app/models"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"

	_ "github.com/lib/pq"
)

// Contains all the basic app info

type application struct {
	appName string
	server  server
	debug   bool
	errLog  *log.Logger
	infoLog *log.Logger
	view    *jet.Set
	session *scs.SessionManager
	Models  models.Models
}

type server struct {
	host string
	port string
	url  string
}

func main() {

	isNotMigrated := flag.Bool("migrate", false, "Should migrate - drop all tables")

	flag.Parse()

	server := server{
		host: "localhost",
		port: "8080",
		url:  "http://localhost:8080",
	}

	host := "localhost"
	port := 5432
	user := "udit"
	password := "postgres"
	dbname := "go_lang"

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	database, err := openDB(psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	upper, err := postgresql.New(database)
	if err != nil {

		log.Fatal(err)
	}

	defer func(upper db.Session) {
		err := upper.Close()
		if err != nil {

			log.Fatal(err)
		}
	}(upper)

	if *isNotMigrated {
		fmt.Println("Running Migration")
		err := migrate(upper)
		if err != nil {
			log.Fatal(err)
		}

	}

	app := &application{
		appName: "News Room",
		server:  server,
		debug:   true,
		errLog:  log.New(os.Stderr, "ERROR-LOG\t", log.Ltime|log.Ldate|log.Llongfile),
		infoLog: log.New(os.Stdout, "INFO-LOG\t", log.Ltime|log.Ldate|log.Lshortfile),
		view:    &jet.Set{},
		session: &scs.SessionManager{},
		Models:  models.New(upper),
	}

	//init session
	app.session = scs.New()
	app.session.Lifetime = 24 * time.Hour
	app.session.Cookie.Persist = true
	app.session.Cookie.Domain = server.url
	app.session.Cookie.SameSite = http.SameSiteStrictMode
	app.session.Store = postgresstore.New(database)

	if app.debug {
		app.view = jet.NewSet(jet.NewOSFileSystemLoader("./views"), jet.InDevelopmentMode())
	} else {
		app.view = jet.NewSet(jet.NewOSFileSystemLoader("./views"))
	}

	if err := app.listenAndServer(); err != nil {

		log.Fatal(err)
	}

	fmt.Println("HEMLOOO")
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {

		return nil, err
	}

	return db, nil
}

func migrate(db db.Session) error {

	script, err := os.ReadFile("./migrations/tables.sql")
	if err != nil {
		log.Fatal(err)

	}
	_, err = db.SQL().Exec(string(script))

	return err

}
