package main

import (
	"fmt"
	"os"

	"github.com/Angelos-Giannis/gitpr/internal/app/infra/command"
	"github.com/urfave/cli"
)

// Define the app details as constants
const (
	appName = "CLI tool to check status of pull requests in github."
	appUsage = ""
	appAuthor = "Angelos Giannis"
	appVersion = "1.0.0"
)

func main() {
	var app = cli.NewApp()
	info(app)

	app.Commands = command.NewBuilder().Hello().GetCommands()

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("Error : %v", err)
		os.Exit(1)
	}
}

// info sets up the information of the tool.
func info(app *cli.App) {
	app.Name = appName
	app.Usage = appUsage
	app.Author = appAuthor
	app.Version = appVersion
}
