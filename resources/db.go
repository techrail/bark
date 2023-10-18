package resources

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

// BarkPostgresDb wraps the *pgxpool.Pool in a custom struct to use it as a receiver for query functions
type BarkPostgresDb struct {
	Client *pgxpool.Pool
}

var ServerDbSaverWg sync.WaitGroup

var BarkDb *BarkPostgresDb

// InitDb : InitDB returns error (if any) encountered while trying to establish a connection to the postgres DB instance.
func InitDb(dbUrl string) error {
	// Connect to Postgres DB instance
	var err error
	BarkDb, err = OpenDb(dbUrl)
	if err != nil {
		return fmt.Errorf("E#1KDZOZ - Opening database failed. Error: %v\n", err)
	}

	// NOTE: The caller must check the connection being returned
	return nil
}

// OpenDb : OpenDB returns a pointer to the `BarkPostgresDb` object.
func OpenDb(dbUrl string) (*BarkPostgresDb, error) {
	connPool, err := pgxpool.NewWithConfig(context.Background(), Config(dbUrl))
	if err != nil {
		return &BarkPostgresDb{}, fmt.Errorf("E#1KDW57 - error connecting to db: %w", err)
	}

	return &BarkPostgresDb{Client: connPool}, nil
}

// CloseDb closes the connection to the DB.
func (d *BarkPostgresDb) CloseDb() {
	d.Client.Close()
}
