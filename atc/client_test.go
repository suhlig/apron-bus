package atc_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/suhlig/apron-bus/atc"
)

var _ = Describe("ATC version client", func() {
	var (
		server *ghttp.Server
		client *atc.VersionClient
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		client = atc.NewVersionClient(server.URL())
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("fetching the server version", func() {
		var (
			err     error
			version string
		)

		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/info"),
					ghttp.RespondWith(http.StatusOK, `{"version":"0.8.15"}`),
				))
		})

		JustBeforeEach(func() {
			version, err = client.GetServerVersion()
		})

		It("makes a request to fetch ATC version", func() {
			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		It("returns no error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns the expected version", func() {
			Expect(version).To(Equal("0.8.15"))
		})
	})
})
