package resources

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

// BarkPostgresDb wraps the *pgxpool.Pool in a custom struct to use it as a receiver for query functions
type BarkPostgresDb struct {
	//Client *sqlx.DB
	Client *pgxpool.Pool
}

var BarkDb *BarkPostgresDb

func InitDB() error {
	// Connect to Postgres DB instance
	var err error
	BarkDb, err = OpenDB()
	if err != nil {
		return fmt.Errorf("E#1KDZOZ - Opening database failed. Error: %v\n", err)
	}
	fmt.Println("E#1KDZG7 - successfully connected to database")
	return nil
}

func OpenDB() (*BarkPostgresDb, error) {
	connPool, err := pgxpool.NewWithConfig(context.Background(), Config())
	if err != nil {
		return &BarkPostgresDb{}, fmt.Errorf("E#1KDW57 - error connecting to db: %w", err)
	}

	return &BarkPostgresDb{Client: connPool}, nil
}

//func OpenDatabase()

func (d *BarkPostgresDb) PingDB(ctx context.Context) error {
	//return d.Client.DB.PingContext(ctx)
	return d.Client.Ping(ctx)
}

// func to close database connection
func (d *BarkPostgresDb) CloseDB() {
	d.Client.Close()
}