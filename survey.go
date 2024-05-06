package main

import (
	//"flag"
	//"fmt"
	//"io/ioutil"
	"log"
	"os"
	"os/signal"
	//"strings"
	"syscall"
	"time"

	"github.com/dchote/survey-engine/config"
	"github.com/dchote/survey-engine/models"
	"github.com/dchote/survey-engine/server"

	"github.com/GeertJohan/go.rice"
	"github.com/docopt/docopt-go"
	"github.com/jinzhu/gorm"
)

// VERSION is probably worthless
const VERSION = "0.0.1"

var (
	db           *gorm.DB
	err          error
	staticAssets *rice.Box
)

func cliArguments() {
	usage := `
Usage: survey [options]

Options:
  -c, --config=<json>           Specify config file [default: ./config.json]
  -h, --help                    Show this screen.
  -v, --version                 Show version.
`
	args, _ := docopt.ParseArgs(usage, os.Args[1:], VERSION)

	config.ConfigFile, _ = args.String("--config")

	_, err = config.LoadConfig(config.ConfigFile)
	if err != nil {
		log.Fatalf("Unable to load "+config.ConfigFile+" ERROR=", err)
	}

	log.Printf("Config: %+v", config.Config)
}

func main() {
	cliArguments()

	staticAssets, err = rice.FindBox("frontend/dist")
	if err != nil {
		log.Fatalf("Static assets not found. Build them with npm first.")
	}

	db = models.OpenDatabase(*config.Config)

	// start the webserver
	go server.StartServer(*config.Config, db, staticAssets)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown

	log.Println("Shutting down")

	// shut down listener, with a hard timeout
	server.StopServer()
	models.CloseDatabase()

	// extra grace time
	time.Sleep(time.Second)

	os.Exit(0)
}
