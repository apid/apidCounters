package main

import (
	"github.com/30x/apid-core"
	"github.com/30x/apid-core/factory"
	_ "github.com/30x/apidCounters"
)

func main() {
	// initialize apid using default services
	apid.Initialize(factory.DefaultServicesFactory())

	log := apid.Log()
	log.Debug("initializing...")

	// this will call all initialization functions on all registered plugins
	apid.InitializePlugins()

	// print the base url to the console
	config := apid.Config()
	basePath := config.GetString("counters_base_path")
	port := config.GetString("api_port")
	log.Print()
	log.Printf("Counters API is at: http://localhost:%s%s", port, basePath)
	log.Print()

	// start client API listener
	api := apid.API()
	err := api.Listen() // doesn't return if no error
	log.Fatalf("Error. Is something already running on port %d? %s", port, err)
}
