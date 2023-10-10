package resources

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

// BarkPostgresDb wraps the *pgxpool.Pool in a custom struct to use it as a receiver for query functions
type BarkPostgresDb struct {
	Client *pgxpool.Pool
}

var BarkDb *BarkPostgresDb

func InitDb() error {
	// Connect to Postgres DB instance
	var err error
	BarkDb, err = OpenDb()
	if err != nil {
		return fmt.Errorf("E#1KDZOZ - Opening database failed. Error: %v\n", err)
	}

	// NOTE: The caller must check the connection being returned
	return nil
}

func OpenDb() (*BarkPostgresDb, error) {
	connPool, err := pgxpool.NewWithConfig(context.Background(), Config())
	if err != nil {
		return &BarkPostgresDb{}, fmt.Errorf("E#1KDW57 - error connecting to db: %w", err)
	}

	return &BarkPostgresDb{Client: connPool}, nil
}

func (d *BarkPostgresDb) CloseDb() {
	d.Client.Close()
}
