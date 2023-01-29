package main

import (
	"flag"
	"fmt"
	"os"

	pg_query "github.com/pganalyze/pg_query_go/v2"
)

func main() {

	config_param := flag.String("config", "pgpac.conf", "The file to pull the configuration from.")
	target_param := flag.String("target", ".", "The target folder or file.")
	server_param := flag.String("server", "", "The server from the config file to deploy the pgpac file to.")
	project_param := flag.String("project", "", "The project from the config file to package into the pgpac file.")
	host_param := flag.String("host", "", "The host server to deploy to.")
	port_param := flag.String("port", "", "The port of the server to deploy to")
	database_param := flag.String("db", "", "The database to deploy to.")
	username_param := flag.String("user", "", "The user to use for authentication with the database server.")
	password_param := flag.String("pass", "", "The password to use for authentication with the database server.")

	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("You must supply a command: init, package, or deploy")
		os.Exit(1)
	}

	commandName := flag.Args()[0]
	config = *config_param
	target = *target_param
	server = *server_param
	project = *project_param
	host = *host_param
	port = *port_param
	database = *database_param
	username = *username_param
	password = *password_param

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
		fmt.Printf("Unknown command '%s'. Expected init, package, deploy, or help.", commandName)
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
