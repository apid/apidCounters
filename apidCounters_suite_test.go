package apidCounters

import (
	"github.com/30x/apid"
	"github.com/30x/apid/factory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	testDir    string
	testServer *httptest.Server
)

var _ = BeforeSuite(func() {
	apid.Initialize(factory.DefaultServicesFactory())

	config := apid.Config()

	// create a temp directory for our database
	var err error
	testDir, err = ioutil.TempDir("", "api_test")
	Expect(err).NotTo(HaveOccurred())

	// set the database data path to our temp dir
	config.Set("data_path", testDir)

	// initialize apid
	// note: this will call initPlugin()
	apid.InitializePlugins()

	// get the router - this will already have the plugin routes registered
	router := apid.API().Router()
	// you could add any additional test or mock routes at this point if needed
	// create our test server
	testServer = httptest.NewServer(router)
})

var _ = AfterSuite(func() {
	apid.Events().Close()
	if testServer != nil {
		testServer.Close()
	}
	os.RemoveAll(testDir)
})

func TestApidCounters(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ApidCounters Suite")
}
