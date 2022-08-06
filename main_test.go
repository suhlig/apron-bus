package main_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"path/filepath"

	"github.com/concourse/concourse/fly/rc"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
	"gopkg.in/yaml.v2"
)

var pathToApronBus string

var _ = BeforeSuite(func() {
	var err error
	pathToApronBus, err = gexec.Build("github.com/suhlig/apron-bus")
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

var _ = Describe("main", func() {
	var (
		err           error
		session       *gexec.Session
		args          []string
		homeDirectory string
		server        *ghttp.Server
		pathToFakeFly string
		serverVersion string
	)

	BeforeEach(func() {
		homeDirectory = GinkgoT().TempDir()
		server = ghttp.NewServer()

		targets := rc.RC{
			Targets: rc.Targets{
				rc.TargetName("mock"): rc.TargetProps{
					API: server.URL(),
				},
			},
		}

		d, err := yaml.Marshal(&targets)
		Expect(err).ToNot(HaveOccurred())

		err = ioutil.WriteFile(filepath.Join(homeDirectory, ".flyrc"), d, 0644)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		server.Close()
	})

	JustBeforeEach(func() {
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/v1/info"),
				ghttp.RespondWith(http.StatusOK, fmt.Sprintf(`{"version":"%v"}`, serverVersion)),
			))

		pathToFakeFly, err = gexec.Build(fmt.Sprintf("github.com/suhlig/apron-bus/fakes/fly%v", serverVersion))
		Expect(err).ToNot(HaveOccurred())

		command := exec.Command(pathToApronBus, args...)
		command.Env = []string{
			fmt.Sprintf("PATH=%v", filepath.Dir(pathToFakeFly)),
			fmt.Sprintf("HOME=%v", homeDirectory),
		}

		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
	})

	// TODO Move setup of fake fly into a separate context so that we can test here without any fly
	XContext("no fly available", func() {
		Context("no args", func() {
			BeforeEach(func() {
				args = []string{}
			})

			It("succeeds", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("exits with code zero", func() {
				Eventually(session).Should(gexec.Exit(0))
			})

			XIt("prints usage", func() {})
		})
	})

	XContext("no fly matching the server version available", func() {
		XIt("prints instructions to fetch the right binary", func() {})
	})

	Context("47.1.1", func() {
		BeforeEach(func() {
			serverVersion = "47.1.1"
			args = []string{"--target", "mock", "--verbose", "status"}
		})

		It("succeeds", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("exits with code zero", func() {
			Eventually(session).Should(gexec.Exit(0))
		})
	})
})
