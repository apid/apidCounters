package apidCounters

import (
	"github.com/30x/apid-core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// BeforeSuite setup and AfterSuite cleanup is in apidCounters_suite_test.go

var _ = Describe("events", func() {

	It("should increment a counter in the database", func(done Done) {

		counter := "event_counter"

		count, err := getDBCounter(counter)
		Expect(err).ShouldNot(HaveOccurred())

		receivedCount := 0

		h := func(e apid.Event) {

			id := e.(string)

			// ignore any events not created by this test
			if id != counter {
				return
			}

			// ignore the first event (so standard listener will definitely process it)
			receivedCount++
			if receivedCount < 2 {
				return
			}

			// now since all events on a channel are guaranteed to be delivered in order,
			// we know the previous events have been completely handled already
			// however, since the 2nd one may or may not be complete, +1 or +2 could be correct

			newCount, err := getDBCounter(counter)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newCount).Should(SatisfyAny(Equal(count+1), Equal(count+2)))

			close(done)
		}

		// add my test handler as a listener
		apid.Events().ListenFunc(counterIncrementEventSelector, h)

		// send 2 increment events
		sendIncrementEvent(counter) // this will definitely be counted
		sendIncrementEvent(counter) // this might not be counted as it's async
	})
})
