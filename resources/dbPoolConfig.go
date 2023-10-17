package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Config retrieves the postgres DB connection url from environment variable named `BARK_DATABASE_URL`.
// It then parses the connection url and checks for empty and invalid urls, if found, it logs the error and connection is not made.
// If it passes the checks, the connection is established and function returns object of pgxpool config.
func Config(dbUrl string) *pgxpool.Config {
	const defaultMaxConns = int32(20)
	const defaultMinConns = int32(5)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute

	dbConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		panic(fmt.Sprintf("Failed to create a config, error: %v", err))
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = 5 * time.Second

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		//	log.Println("Before acquiring the connection pool to the database!!")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		//	log.Println("After releasing the connection pool to the database!!")
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		//	log.Println("Closed the connection pool to the database!!")
	}

	return dbConfig
}
