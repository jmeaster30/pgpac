package main

import (
	"fmt"

	pg_query "github.com/pganalyze/pg_query_go/v2"
)

func main() {
	tree, err := pg_query.ParseToJSON("select * from mytable")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", tree)
}
