package neelkanth

import (
	"context"
	"log"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTableQuery(p *pgxpool.Pool){
	_, err := p.Exec(context.Background(), "CREATE TABLE users (id SERIAL PRIMARY KEY,name VARCHAR(255) NOT NULL,email VARCHAR(255) NOT NULL);")
	if err!=nil {
		log.Fatal("Error creating table")
	}
}

func InsertQuery(p *pgxpool.Pool) {
	_, err := p.Exec(context.Background(), "insert into users(name, email) values($1, $2)", "Neelkanth", "neelkanth@admin.com")
	if err!=nil {
		log.Fatal("Error inserting values into the table")
	}
}