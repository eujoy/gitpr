package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/eujoy/gitpr/internal/app/infra/pullrequests"
	"github.com/eujoy/gitpr/internal/app/infra/repository"
	"github.com/eujoy/gitpr/internal/app/infra/userrepos"
	"github.com/eujoy/gitpr/internal/config"
	"github.com/eujoy/gitpr/internal/infra/command"
	internalHttp "github.com/eujoy/gitpr/internal/infra/route/http"
	"github.com/eujoy/gitpr/pkg/client"
	"github.com/eujoy/gitpr/pkg/printer"
	"github.com/eujoy/gitpr/pkg/publish"
	"github.com/eujoy/gitpr/pkg/utils"
	"github.com/urfave/cli/v2"
)

const (
	// define the configuration file for the system.
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

	gitRepoFactory, err := client.NewFactory(cfg.Settings.DefaultClient, cfg)
	if err != nil {
		fmt.Printf("Error setting up the service : %v\n", err)
		os.Exit(1)
	}

	urSrv := userrepos.NewService(gitRepoFactory.GetClient())
	repoSrv := repository.NewService(gitRepoFactory.GetClient())
	prSrv := pullrequests.NewService(gitRepoFactory.GetClient())

	switch cfg.Service.Mode {
	case "cli":
		startUpCliService(app, cfg, urSrv, prSrv, repoSrv)
	case "http":
		startUpHTTPServer(cfg, urSrv, prSrv)
	default:
		startUpCliService(app, cfg, urSrv, prSrv, repoSrv)
	}
}

// info sets up the information of the tool.
func info(app *cli.App, cfg config.Config) {
	app.Authors = []*cli.Author{
		{
			Name:  cfg.Application.Author,
			Email: "",
		},
	}
	app.Name = cfg.Application.Name
	app.Usage = cfg.Application.Usage
	app.Version = cfg.Application.Version
}

// startUpCliService runs the service as a cli tool.
func startUpCliService(app *cli.App, cfg config.Config, urSrv *userrepos.Service, prSrv *pullrequests.Service, repoSrv *repository.Service) {
	u := utils.New(cfg)
	tp := printer.NewTablePrinter()

	googleSheetsService, err := publish.NewGoogleSheetsService()
	if err != nil {
		fmt.Printf("Failed to prepare google sheets service with error: %v\n", err)
		os.Exit(1)
	}

	b := command.NewBuilder(cfg, urSrv, prSrv, repoSrv, tp, u, googleSheetsService)

	app.Commands = b.
		Find().
		PullRequests().
		UserRepos().
		Widget().
		CommitList().
		CreateRelease().
		CreatedPullRequests().
		ReleaseReport().
		PublishPullRequestMetrics().
		GetCommands()

	err = app.Run(os.Args)
	if err != nil {
		fmt.Printf("Error running the service : %v\n", err)
		os.Exit(1)
	}
}

// startUpHttpServer runs the service in http rest mode.
func startUpHTTPServer(cfg config.Config, urSrv *userrepos.Service, prSrv *pullrequests.Service) {
	rh := internalHttp.NewHandler(cfg, urSrv, prSrv)

	http.HandleFunc("/settings", rh.GetSettings)
	http.HandleFunc("/userRepos", rh.GetUserRepos)
	http.HandleFunc("/pullRequests", rh.GetPullRequestsOfRepository)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", cfg.Service.Port), nil))
}
