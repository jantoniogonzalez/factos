package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

func ConnectDB() {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Error(err)
		fmt.Printf("Unable to create connection pool \n")
		os.Exit(1)
	}

	defer dbpool.Close()
}
