package main

import (
	"flag"
	"fmt"
	"os"

	pg_query "github.com/pganalyze/pg_query_go/v2"
)

func main() {

	target := flag.String("target", ".", "The target folder or file.")
	server := flag.String("server", "", "The server to deploy the pgpac file to.")
	database := flag.String("db", "", "The database to deploy the pgpac file to.")
	username := flag.String("user", "", "The user to use for authentication with the database server.")
	password := flag.String("pass", "", "The password to use for authentication with the database server.")

	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("You must supply a command: init, package, or deploy")
		os.Exit(1)
	}

	commandName := flag.Args()[0]
	_ = *target
	_ = *server
	_ = *database
	_ = *username
	_ = *password

	switch commandName {
	case "help":
		fmt.Printf(`Commands:
	init    - initialize the folder structure for pgpac.
	package - package the target folder into a .pgpac file to be deployed.
	deploy  - deploy your pgpac folder or the prepackaged file to your database.
`)
		flag.PrintDefaults()
		os.Exit(0)
	case "init":
		break
	case "package":
		break
	case "deploy":
		break
	default:
		fmt.Printf("Unknown command '%s'. Expected init, package, or deploy.", commandName)
		os.Exit(1)
	}

	fmt.Printf("%s\n", commandName)

	tree, err := pg_query.Parse(`
	-- this is a comment >:)
	select * from mytable
	`)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", tree.String())
}
