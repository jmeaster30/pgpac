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
		Deploy(cli.config, cli.project, cli.host, cli.port, cli.user, cli.pass, cli.db)
	default:
		fmt.Printf("Unknown command '%s'. Expected init, package, deploy, or help.", cli.command)
		os.Exit(1)
	}
}
