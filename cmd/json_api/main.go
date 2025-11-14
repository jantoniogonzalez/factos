package main

import (
	"database/sql"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	localDB "github.com/jantoniogonzalez/factos/internal/db"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type application struct {
	logger          *slog.Logger
	sessionManager  *scs.SessionManager
	googleoauthconf *oauth2.Config
	factos          *localDB.FactosModel
	fixtures        *localDB.FixturesModel
	leagues         *localDB.LeaguesModel
	teams           *localDB.TeamsModel
	users           *localDB.UserModel
}

func main() {
	logLevel := slog.LevelInfo
	if os.Getenv("DEBUG") == "true" {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))

	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	port := flag.String("port", os.Getenv("DEV_PORT"), "Development port")
	dsn := flag.String("dsn", os.Getenv("DSN"), "MySQL Connection DSN")
	dev_domain := flag.String("dev_domain", os.Getenv("DEV_SESSION_DOMAIN"), "Dev Session Domain")

	db, err := openDB(*dsn)

	if err != nil {
		logger.Error("Failed to open database", "error", err)
		os.Exit(1)
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

	factos := localDB.NewFactosModel(db)
	fixtures := localDB.NewFixturesModel(db)
	leagues := localDB.NewLeaguesModel(db)
	teams := localDB.NewTeamsModel(db)
	users := localDB.NewUserModel(db)

	app := &application{
		logger:          logger,
		sessionManager:  sessionManager,
		googleoauthconf: googleoauthconf,
		factos:          factos,
		fixtures:        fixtures,
		leagues:         leagues,
		teams:           teams,
		users:           users,
	}

	srv := &http.Server{
		Addr:         *port,
		Handler:      app.routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("Starting up server",
		"port", *port,
	)
	logger.Debug("Debug mode is on!")
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
