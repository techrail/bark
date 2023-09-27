package resources

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Config() *pgxpool.Config {
	const defaultMaxConns = int32(20)
	const defaultMinConns = int32(5)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute

	fmt.Printf("Database connection string from Environment: %s\n", os.Getenv("BARK_DATABASE_URL"))

	dbConfig, err := pgxpool.ParseConfig(os.Getenv("BARK_DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
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
