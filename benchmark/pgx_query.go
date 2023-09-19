package pgx_query

import (
	pgx "github.com/jackc/pgx/v5"
	"fmt"
	"context"
)

type PGX_S struct {
 	Con 	   *pgx.Conn	
}

func(p *PGX_S) PGX_Connect(cs string) error {
	var err error 
	p.Con, err = pgx.Connect(context.Background(), cs)
	if err != nil {
		return fmt.Errorf("Unable to connect to Database: %s", err)
	}
	return nil 
}

func(p *PGX_S) PGX_Query(q string)(pgx.Rows , error ){
	res, err := p.Con.Query(context.Background(),q)
	if err != nil {
	 	return res , err 
	}
	return res, nil
}

