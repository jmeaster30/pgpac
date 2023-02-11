package main

import (
	"fmt"
	"os"
)

func main() {
	PrintBanner()
	cli := ParseCli()

	switch cli.command {
	case "help":
		PrintHelp()
		os.Exit(0)
	case "init":
		InitializeProject(cli.config)
	case "package":
		break
	case "deploy":
		break
	default:
		fmt.Printf("Unknown command '%s'. Expected init, package, deploy, or help.", cli.command)
		os.Exit(1)
	}

	/*
		fmt.Printf("%s\n", commandName)

		tree, err := pg_query.Parse(`
		-- this is a comment >:)
		select * from mytable
		`)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", tree.String())*/
}
