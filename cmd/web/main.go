package main

import (
	"fmt"
	"log"
	"os"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
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
}

type server struct {
	host string
	port string
	url  string
}

func main() {

	server := server{
		host: "localhost",
		port: "8080",
		url:  "http://localhost:8080",
	}

	app := &application{
		appName: "News Room",
		server:  server,
		debug:   true,
		infoLog: log.New(os.Stdout, "INFO-LOG\t", log.Ltime|log.Ldate|log.Lshortfile),
		errLog:  log.New(os.Stderr, "ERROR-LOG\t", log.Ltime|log.Ldate|log.Llongfile),
	}

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
