package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Angelos-Giannis/gitpr/internal/app/infra/command"
	"github.com/Angelos-Giannis/gitpr/internal/config"
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
	configurationFile = "configuration.yaml"
)

func main() {
	cfg, err := config.New(configurationFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var app = cli.NewApp()
	info(app, cfg)

	gcl := githttp.NewClient(&http.Client{Timeout: cfg.Github.Timeout * time.Second}, cfg)
	gr := github.NewResource(gcl)

	tp := printer.NewTablePrinter()
	u := utils.NewUtils(cfg)

	urSrv := userrepos.NewService(gr)
	prSrv := pullrequests.NewService(gr)

	b := command.NewBuilder(cfg, urSrv, prSrv, tp, u)

	app.Commands = b.
		UserRepos().
		PullRequests().
		GetCommands()

	err = app.Run(os.Args)
	if err != nil {
		fmt.Printf("Error : %v", err)
		os.Exit(1)
	}
}

// info sets up the information of the tool.
func info(app *cli.App, cfg config.Config) {
	app.Author = cfg.Application.Author
	app.Name = cfg.Application.Name
	app.Usage = cfg.Application.Usage
	app.Version = cfg.Application.Version
}
