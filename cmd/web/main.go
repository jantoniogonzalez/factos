package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jantoniogonzalez/factos/internal/models"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	factos   *models.Factos
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERRROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	port := flag.String("port", ":4000", "HTTP network access")
	dsn := flag.String("dsn", os.Getenv("DSN"), "MySQL connection string")

	flag.Parse()

	db, err := openDB(*dsn)

	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:         *port,
		Handler:      app.routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	err = srv.ListenAndServe()
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, err
}
