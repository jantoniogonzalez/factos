package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	localDB "github.com/jantoniogonzalez/factos/internal/db"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	factos         *localDB.FactosModel
	fixtures       *localDB.FixturesModel
	users          *localDB.UserModel
	oauthConfig    *oauth2.Config
	sessionManager *scs.SessionManager
	formDecoder    *form.Decoder
	cachedFiles    map[string]*template.Template
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERRROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	err := godotenv.Load()
	if err != nil {
		errorLog.Fatal(err)
	}

	port := flag.String("port", ":4000", "HTTP network access")
	dsn := flag.String("dsn", os.Getenv("DSN"), "MySQL connection string")

	flag.Parse()

	cachedFiles, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	conf := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URI"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	db, err := openDB(*dsn)

	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.NewWithCleanupInterval(db, 60*time.Minute)
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Name = "factos_session"
	sessionManager.Cookie.Path = "/"

	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		oauthConfig:    conf,
		factos:         localDB.NewFactosModel(db),
		fixtures:       localDB.NewFixturesModel(db),
		users:          localDB.NewUserModel(db),
		sessionManager: sessionManager,
		formDecoder:    formDecoder,
		cachedFiles:    cachedFiles,
	}

	srv := &http.Server{
		Addr:         *port,
		Handler:      app.routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	infoLog.Printf("Starting up the server in port %s\n", *port)
	err = srv.ListenAndServe()

	errorLog.Fatal(err)
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
