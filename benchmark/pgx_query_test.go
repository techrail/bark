package pgx_query

import (
	"testing"
	//"fmt"
)

func BenchmarkConnect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var p PGX_S
		_ =p.PGX_Connect("postgres://bhanu:bhanu@http://localhost:5432/bhanu_db?")
		
	}
}


// func BenchmarkQuery(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		var p PGX_S
// 		err :=p.PGX_Connect("postgres://bhanu:bhanu@http://localhost:5432/bhanu_db?")
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		_, err  =p.PGX_Query("SELECT * FROM department")
// 		if err != nil {
// 			fmt.Println(err)
// 		}
		
// 	}
// }
