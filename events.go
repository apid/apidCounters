package apidCounters

import (
	"github.com/apid/apid-core"
)

const (
	counterIncrementEventSelector apid.EventSelector = "counter_increment"
)

var (
	events apid.EventsService
)

func initEvents(services apid.Services) {
	events = services.Events()
	events.ListenFunc(counterIncrementEventSelector, receivedIncrementEvent)
}

// receives an events for counterEventKey selector and requests DB update
func receivedIncrementEvent(e apid.Event) {
	log.Debugf("received event increment: %s", e)
	id, ok := e.(string)
	if ok {
		incrementDBCounter(id)
	}
}

// publishes events to counterEventKey selector for a given ID
func sendIncrementEvent(id string) {
	log.Debugf("sending event increment: %s", id)
	events.Emit(counterIncrementEventSelector, id)
}
