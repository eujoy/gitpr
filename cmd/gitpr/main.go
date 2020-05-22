package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Angelos-Giannis/gitpr/internal/app/infra/pullrequests"
	"github.com/Angelos-Giannis/gitpr/internal/app/infra/userrepos"
	"github.com/Angelos-Giannis/gitpr/internal/config"
	"github.com/Angelos-Giannis/gitpr/internal/infra/command"
	internalHttp "github.com/Angelos-Giannis/gitpr/internal/infra/route/http"
	"github.com/Angelos-Giannis/gitpr/pkg/client"
	"github.com/Angelos-Giannis/gitpr/pkg/printer"
	"github.com/Angelos-Giannis/gitpr/pkg/utils"
	"github.com/urfave/cli"
)

// Define the app details as constants
const (
	configurationFile = "configuration.yaml"
)

// main sets up the dependencies as well as starts up the service itself.
func main() {
	cfg, err := config.New(configurationFile)
	if err != nil {
		fmt.Printf("Error parsing configuration : %v\n", err)
		os.Exit(1)
	}

	var app = cli.NewApp()
	info(app, cfg)

	gitRepoFactory, err := client.NewFactory(cfg.Clients.Default, cfg)
	if err != nil {
		fmt.Printf("Error setting up the service : %v\n", err)
		os.Exit(1)
	}

	urSrv := userrepos.NewService(gitRepoFactory.GetClient())
	prSrv := pullrequests.NewService(gitRepoFactory.GetClient())

	switch cfg.Service.Mode {
	case "cli":
		startUpCliService(app, cfg, urSrv, prSrv)
	case "http":
		startUpHttpServer(cfg, urSrv, prSrv)
	default:
		startUpCliService(app, cfg, urSrv, prSrv)
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

	http.HandleFunc("/settings", rh.GetSettings)
	http.HandleFunc("/userRepos", rh.GetUserRepos)
	http.HandleFunc("/pullRequests", rh.GetPullRequestsOfRepository)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", cfg.Service.Port), nil))
}
