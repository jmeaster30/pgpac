package main

import (
	"flag"
	"fmt"
)

type CliFlags struct {
	command string
	config  string
	server  string
	project string
	host    string
	port    string
	db      string
	user    string
	pass    string
}

// config_param := flag.String("config", "pgpac.yaml", "The file to pull the configuration from.")
// target_param := flag.String("target", ".", "The target folder or file.")
// server_param := flag.String("server", "", "The server from the config file to deploy the pgpac file to.")
// project_param := flag.String("project", "", "The project from the config file to package into the pgpac file.")
// host_param := flag.String("host", "", "The host server to deploy to.")
// port_param := flag.String("port", "", "The port of the server to deploy to")
// database_param := flag.String("db", "", "The database to deploy to.")
// username_param := flag.String("user", "", "The user to use for authentication with the database server.")
// password_param := flag.String("pass", "", "The password to use for authentication with the database server.")

func PrintHelp() {
	fmt.Println("Commands:")
	fmt.Println("\tinit     initialize the folder structure for pgpac.")
	fmt.Println("\tpackage  package the target folder into a .pgpac file to be deployed.")
	fmt.Println("\tdeploy   deploy your pgpac folder or the prepackaged file to your database.")
	fmt.Println("Arguments:")
	fmt.Println("\tconfig <filename>          The configuration file for projects and servers")
	fmt.Println("\t\tDefault: 'pgpac.yaml'")
	fmt.Println("\tserver <server>            The name of the server in the config file")
	fmt.Println("\tproject <project name>     The name of the project in the config file")
	fmt.Println("\thost <db server hostname>  Hostname of the database server")
	fmt.Println("\tport <db server port>      Port of the database server")
	fmt.Println("\tdb <db name>               Database name on the server")
	fmt.Println("\tuser <db username>         Username to use for authentication with the database")
	fmt.Println("\tpass <db password>         Password to use for authentication with the database")
}

func ParseCli() CliFlags {
	flag.Parse()
	args := flag.Args()
	flags := CliFlags{
		command: "",
		config:  "pgpac.yaml",
		server:  "",
		project: "",
		host:    "",
		port:    "",
		db:      "",
		user:    "",
		pass:    "",
	}

	for i := 0; i < len(args); i++ {
		flag := args[i]
		if i == 0 {
			flags.command = flag
		} else {
			switch flag {
			case "-c":
				fallthrough
			case "-config":
				fallthrough
			case "--config":
				fallthrough
			case "config":
				flags.config = args[i+1]
			case "-s":
				fallthrough
			case "-server":
				fallthrough
			case "--server":
				fallthrough
			case "server":
				flags.server = args[i+1]
			case "-t":
				fallthrough
			case "-project":
				fallthrough
			case "--project":
				fallthrough
			case "project":
				flags.project = args[i+1]
			case "-h":
				fallthrough
			case "-host":
				fallthrough
			case "--host":
				fallthrough
			case "host":
				flags.host = args[i+1]
			case "-p":
				fallthrough
			case "-port":
				fallthrough
			case "--port":
				fallthrough
			case "port":
				flags.port = args[i+1]
			case "-d":
				fallthrough
			case "-database":
				fallthrough
			case "--database":
				fallthrough
			case "database":
				flags.db = args[i+1]
			case "-u":
				fallthrough
			case "-user":
				fallthrough
			case "--user":
				fallthrough
			case "user":
				flags.user = args[i+1]
			case "-pass":
				fallthrough
			case "-password":
				fallthrough
			case "--password":
				fallthrough
			case "password":
				flags.pass = args[i+1]
			}
			i += 1
		}
	}

	return flags
}
