package neelkanth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
)

func h(connPool *pgxpool.Pool, batch *pgx.Batch) {
	// db.InsertQuery(connPool)	
	connPool.SendBatch(context.Background(), batch)
}

func main() {
	router := mux.NewRouter()

    router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "You've requested the server")
    })

	//Create database connection
	connPool,err := pgxpool.NewWithConfig(context.Background(), db.Config())
	if err!=nil {
		log.Fatal("Error creating connection to the database!!")
	} 

	numBatches:= 10

	batches := make([]*pgx.Batch, 0, numBatches)
	batchSize := 100

	for test:=0; test<10;test++ {

		batch := &pgx.Batch{}
		for i := 0; i < batchSize; i++ {
			batch.Queue("insert into users(name, email) values($1, $2)", "Neelkanth", "neelkanth@admin.com")
		}
		batches = append(batches, batch)
	}

	x:= time.Now()
    for _, batch := range batches {
		// h(connPool, batch)
		go h(connPool, batch)
	}
	// This would not work in case of goroutines 
	fmt.Println("Total time: ", time.Since(x))

	fmt.Println("Connected to the database!!")
	
	CreateTableQuery(connPool)
	
	InsertQuery(connPool)

	defer connPool.Close()

	http.ListenAndServe(":8080", router)
}
