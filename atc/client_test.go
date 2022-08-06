package atc_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/suhlig/apron-bus/atc"
)

var _ = Describe("ATC version client", func() {
	var server *ghttp.Server
	var client *atc.VersionClient

	BeforeEach(func() {
		server = ghttp.NewServer()
		client = atc.NewVersionClient(server.URL())
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("fetching ATC version", func() {
		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.VerifyRequest("GET", "/api/v1/info"),
			)
		})

		It("should make a request to fetch ATC version", func() {
			client.GetServerVersion()
			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})
	})
})
