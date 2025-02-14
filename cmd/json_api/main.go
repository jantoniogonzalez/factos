package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type application struct {
	logger          *log.Logger
	sessionManager  *scs.SessionManager
	googleoauthconf *oauth2.Config
}

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	port := flag.String("port", os.Getenv("DEV_PORT"), "Development port")
	dsn := flag.String("dsn", os.Getenv("DSN"), "MySQL Connection DSN")
	dev_domain := flag.String("dev_domain", os.Getenv("DEV_SESSION_DOMAIN"), "Dev Session Domain")

	db, err := openDB(*dsn)

	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.NewWithCleanupInterval(db, 1*time.Hour)
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.IdleTimeout = 30 * time.Minute
	sessionManager.Cookie.Name = "factos_session"
	sessionManager.Cookie.Domain = *dev_domain
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Path = "/"
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Cookie.Secure = true

	// Add Oauth2 Config
	googleoauthconf := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URI"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	app := &application{
		logger:          logger,
		sessionManager:  sessionManager,
		googleoauthconf: googleoauthconf,
	}

	srv := &http.Server{
		Addr:         *port,
		Handler:      app.routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
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

	return db, nil
}
