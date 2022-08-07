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
		err     error
		session *gexec.Session
		command *exec.Cmd
	)

	JustBeforeEach(func() {
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
	})

	Context("no fly available", func() {
		Context("no args", func() {
			BeforeEach(func() {
				command = exec.Command(pathToApronBus, []string{}...)
			})

			It("succeeds", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("exits with code zero", func() {
				Eventually(session).Should(gexec.Exit(0))
			})

			It("prints usage to STDOUT", func() {
				Expect(session.Wait().Out.Contents()).To(ContainSubstring("fly"))
			})
		})
	})

	Describe("calling of fly", func() {
		var (
			homeDirectory string
			server        *ghttp.Server
			pathToFakeFly string
			serverVersion string
			flyVersion    string
			args          []string
		)

		BeforeEach(func() {
			serverVersion = "47.1.1"
			server = newMockServer(serverVersion)

			homeDirectory = GinkgoT().TempDir()
			err = writeFlyRC(homeDirectory, server.URL())
			Expect(err).ToNot(HaveOccurred())

			args = []string{"--target", "mock", "--verbose", "status"}
		})

		Context("NO fly matching the server's version is avilable", func() {
			BeforeEach(func() {
				command = exec.Command(pathToApronBus, args...)
				command.Env = []string{
					fmt.Sprintf("HOME=%v", homeDirectory),
				}
			})

			It("starts", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("exits with non-zero code", func() {
				Eventually(session).ShouldNot(gexec.Exit(0))
			})

			It("prints the error on STDERR", func() {
				Expect(string(session.Wait().Err.Contents())).To(ContainSubstring("could not find fly47.1.1"))
			})
		})

		Context("fly matching the server's version IS avilable", func() {
			BeforeEach(func() {
				flyVersion = serverVersion
				pathToFakeFly, err = gexec.Build(fmt.Sprintf("github.com/suhlig/apron-bus/fakes/fly%v", flyVersion))
				Expect(err).ToNot(HaveOccurred())

				command = exec.Command(pathToApronBus, args...)
				command.Env = []string{
					fmt.Sprintf("PATH=%v", filepath.Dir(pathToFakeFly)),
					fmt.Sprintf("HOME=%v", homeDirectory),
				}
			})

			AfterEach(func() {
				server.Close()
			})

			It("starts", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("exits with code zero", func() {
				Eventually(session).Should(gexec.Exit(0))
			})
		})
	})
})

func writeFlyRC(homeDirectory, serverURL string) error {
	targets := rc.RC{
		Targets: rc.Targets{
			rc.TargetName("mock"): rc.TargetProps{
				API: serverURL,
			},
		},
	}

	data, err := yaml.Marshal(&targets)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(homeDirectory, ".flyrc"), data, 0644)
}

func newMockServer(version string) *ghttp.Server {
	server := ghttp.NewServer()

	server.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/api/v1/info"),
			ghttp.RespondWith(http.StatusOK, fmt.Sprintf(`{"version":"%v"}`, version)),
		))

	return server
}
