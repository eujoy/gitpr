package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Angelos-Giannis/gitpr/internal/app/infra/pullrequests"
	"github.com/Angelos-Giannis/gitpr/internal/app/infra/userrepos"
	"github.com/Angelos-Giannis/gitpr/internal/config"
	"github.com/Angelos-Giannis/gitpr/internal/infra/command"
	internalHttp "github.com/Angelos-Giannis/gitpr/internal/infra/route/http"
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
		fmt.Printf("Error parsing configuration : %v\n", err)
		os.Exit(1)
	}

	var app = cli.NewApp()
	info(app, cfg)

	gcl := githttp.NewClient(&http.Client{Timeout: cfg.Clients.Github.Timeout * time.Second}, cfg)
	gr := github.NewResource(gcl)

	urSrv := userrepos.NewService(gr)
	prSrv := pullrequests.NewService(gr)

	switch cfg.Service.Mode {
	case "cli":
		startUpCliService(app, cfg, urSrv, prSrv)
	case "http":
		startUpHttpServer(cfg, urSrv, prSrv)
	default:
		fmt.Println(cfg.Service)

		err := errors.New("invalid service mode provided in configuration")
		fmt.Printf("Error starting the service : %v\n", err)
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

// startUpCliService runs the service as a cli tool.
func startUpCliService(app *cli.App, cfg config.Config, urSrv *userrepos.Service, prSrv *pullrequests.Service) {
	tp := printer.NewTablePrinter()
	u := utils.NewUtils(cfg)

	b := command.NewBuilder(cfg, urSrv, prSrv, tp, u)

	app.Commands = b.
		UserRepos().
		PullRequests().
		GetCommands()

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("Error running the service : %v\n", err)
		os.Exit(1)
	}
}

// startUpHttpServer runs the service in http rest mode.
func startUpHttpServer(cfg config.Config, urSrv *userrepos.Service, prSrv *pullrequests.Service) {
	rh := internalHttp.NewHandler(cfg, urSrv, prSrv)

	http.HandleFunc("/defaults", rh.GetDefaultSettings)
	http.HandleFunc("/userRepos", rh.GetUserRepos)
	http.HandleFunc("/pullRequests", rh.GetPullRequestsOfRepository)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", cfg.Service.Port), nil))
}
