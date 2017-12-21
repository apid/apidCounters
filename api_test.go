package apidCounters

import (
	"encoding/json"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// BeforeSuite setup and AfterSuite cleanup is in apidCounters_suite_test.go

var _ = Describe("api", func() {

	It("should get a counter", func() {

		counter := "api_test_1"
		count, err := getDBCounter(counter)
		Expect(err).NotTo(HaveOccurred())
		Expect(count).To(Equal(0))

		uri, err := url.Parse(testServer.URL)
		Expect(err).NotTo(HaveOccurred())
		uri.Path = fmt.Sprintf(countersBasePath+"/%s", counter)

		resp, err := http.Get(uri.String())
		Expect(err).ShouldNot(HaveOccurred())
		defer resp.Body.Close()
		Expect(resp.StatusCode).Should(Equal(http.StatusOK))

		body, _ := ioutil.ReadAll(resp.Body)
		Expect(string(body)).Should(Equal(strconv.Itoa(count)))
	})

	It("should increment a counter", func() {

		counter := "api_test_2"
		count, err := getDBCounter(counter)
		Expect(err).NotTo(HaveOccurred())
		Expect(count).To(Equal(0))

		uri, err := url.Parse(testServer.URL)
		Expect(err).NotTo(HaveOccurred())
		uri.Path = fmt.Sprintf(countersBasePath+"/%s", counter)

		req, err := http.NewRequest("POST", uri.String(), nil)
		req.Header.Add("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		Expect(err).ShouldNot(HaveOccurred())
		defer resp.Body.Close()
		Expect(resp.StatusCode).Should(Equal(http.StatusOK))

		body, err := ioutil.ReadAll(resp.Body)
		Expect(string(body)).Should(Equal(strconv.Itoa(count)))

		time.Sleep(50 * time.Millisecond) // todo: hack - make this a listen event on incr event finished
		count, err = getDBCounter(counter)
		Expect(err).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("should retrieve a list of counters", func() {

		counter1 := "api_test_3"
		counter2 := "api_test_4"
		err := incrementDBCounter(counter1)
		Expect(err).NotTo(HaveOccurred())
		err = incrementDBCounter(counter2)
		Expect(err).NotTo(HaveOccurred())
		err = incrementDBCounter(counter2)
		Expect(err).NotTo(HaveOccurred())

		uri, err := url.Parse(testServer.URL)
		Expect(err).NotTo(HaveOccurred())
		uri.Path = countersBasePath

		resp, err := http.Get(uri.String())
		Expect(err).ShouldNot(HaveOccurred())
		defer resp.Body.Close()
		Expect(resp.StatusCode).Should(Equal(http.StatusOK))

		body, err := ioutil.ReadAll(resp.Body)
		Expect(err).ShouldNot(HaveOccurred())
		var counterMap map[string]int
		json.Unmarshal(body, &counterMap)

		Expect(counterMap[counter1]).Should(Equal(1))
		Expect(counterMap[counter2]).Should(Equal(2))
	})
})
