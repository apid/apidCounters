package apidCounters

import (
	"fmt"
	"github.com/30x/apid"
)

const (
	configCountersBasePath  = "counters_base_path"
	countersBasePathDefault = "/counters"
)

// keep track of the services that this plugin will use
// note: services would also be available directly via the package global "apid" (eg. `apid.Log()`)
var (
	log    apid.LogService
	config apid.ConfigService
)

// apid.RegisterPlugin() is required to be called in init()
func init() {
	apid.RegisterPlugin(initPlugin)
}

// initPlugin will be called by apid to initialize
func initPlugin(services apid.Services) error {

	// set a logger that is annotated for this plugin
	log = services.Log().ForModule("counters")
	log.Debug("start init")

	// set configuration
	config = services.Config()

	// set plugin config defaults
	config.SetDefault(configCountersBasePath, countersBasePathDefault)

	// check for any missing required configuration values
	// in this example we check for someConfigurationKey, but normally we wouldn't check defaulted values
	for _, key := range []string{configCountersBasePath} {
		if !config.IsSet(key) {
			return fmt.Errorf("Missing required config value: %s", key)
		}
	}

	// register APIs (see api.go)
	initAPI(services)

	// register for events (see events.go)
	initEvents(services)

	// set data service (see data.go)
	err := initDB(services)
	if err != nil {
		return err
	}

	log.Debug("end init")
	return nil
}
