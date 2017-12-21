# apidCounters

[![Build Status](https://travis-ci.org/apid/apidCounters.svg)](https://travis-ci.org/apid/apidCounters) [![GoDoc](https://godoc.org/github.com/apid/apidCounters?status.svg)](https://godoc.org/github.com/apid/apidCounters) [![Go Report Card](https://goreportcard.com/badge/github.com/apid/apidCounters)](https://goreportcard.com/report/github.com/apid/apidCounters)

This is an EXAMPLE plugin for [apid](http://github.com/apid/apid). It shows how to register as a plugin and 
use apid's API, Config, Data, Events, and Log Services. In addition, it shows how to write integration tests
that use these services. This is not necessary intended to be a production-ready apid component. 


## Functional description

This plugin simply tracks counters based on called URIs:
 
* `GET /counters` - retrieve all counters
* `GET /counters/{name}` - retrieve named counter
* `POST /counters/{name}` - increment named counter (returns prior value) 


## Apid Services Used

* Config Service
* Log Service
* API Service
* Event Service: to publish counter increment events 
* Data Service: to store counters 


## Building and running

First, install prerequisites:
 
    glide install

To run an apid test instance with just this plugin installed, change to the `cmd/apidExample` folder. 
From here, you may create an executable with: 

    go build 
  
Alternatively, you may run without creating an executable with:

    go run main.go 
 
Either way, once the process is running, you should be able to manually give the plugin's API a whirl...

    curl -i -X POST localhost:9000/counters/counter1 
    curl -i localhost:9000/counters/counter1
    curl -i localhost:9000/counters

Note that both port value and the /counters base path have been exposed as configurable. Thus, the following 
env variables will affect the values used for these, respectively:

* APID_API_PORT
* APID_COUNTERS_BASE_PATH

For example, if you set the `APID_COUNTERS_BASE_PATH` env var to `/test`, all the paths above will 
be based at `/test` instead of `/counters`.


## Running tests

To run the tests, just run:

    go test
    
To generate coverage, you may run:

    ./cover.sh

Then open `cover.html` with your browser to see details.
