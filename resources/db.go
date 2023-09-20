package resources

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

// BarkPostgresDb wraps the sqlx.DB in a custom struct to use it as a receiver for query functions
type BarkPostgresDb struct {
	//Client *sqlx.DB
	Client *pgxpool.Pool
}

var BarkDb *BarkPostgresDb

func InitDatabase() error {
	// Connect to Postgres DB instance
	var err error
	BarkDb, err = OpenDatabase()
	if err != nil {
		return fmt.Errorf("E#1KDZOZ - Opening database failed. Error: %v\n", err)
	}

	// Ping DB
	// if err = BarkDb.Ping(context.Background()); err != nil {
	// 	return fmt.Errorf("E#1KDZPY - Opening database failed. Error: %v\n", err)
	// }

	fmt.Println("E#1KDZG7 - successfully connected to database")
	return nil
}

func OpenDatabase() (*BarkPostgresDb, error) {

	// databaseURL := os.Getenv("DATABASE_URL")
	// if strings.TrimSpace(os.Getenv("DATABASE_URL")) == "" {
	// 	return &BarkPostgresDb{}, fmt.Errorf("No env found or empty")
	// }

	//dbConn, err := sqlx.Open("postgres", databaseURL)

	connPool, err := pgxpool.NewWithConfig(context.Background(), Config())
	if err != nil {
		return &BarkPostgresDb{}, fmt.Errorf("E#1KDW57 - error connecting to db: %w", err)
	}

	return &BarkPostgresDb{Client: connPool}, nil
}

//func OpenDatabase()

func (d *BarkPostgresDb) Ping(ctx context.Context) error {
	//return d.Client.DB.PingContext(ctx)
	return d.Client.Ping(ctx)
}
