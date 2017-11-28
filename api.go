package apidCounters

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/apid/apid-core"
)

var countersBasePath string

func initAPI(services apid.Services) {
	countersBasePath = config.GetString(configCountersBasePath)
	services.API().HandleFunc(countersBasePath, returnCounters).Methods("GET")
	services.API().HandleFunc(countersBasePath+"/{counter}", incrAndReturnCount).Methods("GET", "POST")
}

// called from GET or POST /counters/{counter}
// GET will just returns the counter, POST will also increment the counter
func incrAndReturnCount(w http.ResponseWriter, r *http.Request) {
	log.Debugf("incr and return count request: %s", r.URL)

	vars := apid.API().Vars(r)
	id := vars["counter"]

	// send event to increment counter
	if r.Method == "POST" {
		sendIncrementEvent(id)
	}

	count, err := getDBCounter(id)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(strconv.Itoa(count)))
}

// called from GET /counters, will return all counters
func returnCounters(w http.ResponseWriter, r *http.Request) {
	log.Debugf("return counts request: %s", r.URL)

	counters, err := getDBCounters()
	if err != nil {
		writeError(w, err)
		return
	}

	body, err := json.Marshal(counters)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(body)
}

func writeError(w http.ResponseWriter, err error) {
	log.Errorf("error handling API request: %s", err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}
