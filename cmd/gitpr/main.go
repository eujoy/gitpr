package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Angelos-Giannis/gitpr/internal/app/infra/command"
	"github.com/Angelos-Giannis/gitpr/internal/infra/pullrequests"
	"github.com/Angelos-Giannis/gitpr/internal/infra/userrepos"
	"github.com/Angelos-Giannis/gitpr/pkg/github"
	githttp "github.com/Angelos-Giannis/gitpr/pkg/github/http"
	"github.com/Angelos-Giannis/gitpr/pkg/printer"
	"github.com/Angelos-Giannis/gitpr/pkg/utils"
	"github.com/urfave/cli"
)

// Define the app details as constants
const (
	appName = "CLI tool to check status of pull requests in github."
	appUsage = ""
	appAuthor = "Angelos Giannis"
	appVersion = "1.0.0"

	requestTimeout = 5
)

func main() {
	var app = cli.NewApp()
	info(app)

	gcl := githttp.NewClient(&http.Client{Timeout: requestTimeout * time.Second})
	gr := github.NewResource(gcl)

	tp := printer.NewTablePrinter()
	u := utils.NewUtils()

	urSrv := userrepos.NewService(gr)
	prSrv := pullrequests.NewService(gr)

	b := command.NewBuilder(urSrv, prSrv, tp, u)

	app.Commands = b.
		UserRepos().
		PullRequests().
		GetCommands()

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
