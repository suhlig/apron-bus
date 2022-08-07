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
		args    []string
	)

	BeforeEach(func() {
		command = exec.Command(pathToApronBus, args...)
	})

	JustBeforeEach(func() {
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
	})

	Context("no fly available", func() {
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

			It("prints usage", func() {
				Expect(session.Wait().Out.Contents()).To(ContainSubstring("fly"))
			})
		})
	})

	XContext("no fly matching the server version available", func() {
		XIt("prints instructions to fetch the right binary", func() {})
	})

	Context("47.1.1", func() {
		const serverVersion = "47.1.1"

		var (
			homeDirectory string
			server        *ghttp.Server
			pathToFakeFly string
		)

		BeforeEach(func() {
			server = newMockServer(serverVersion)

			homeDirectory = GinkgoT().TempDir()
			err = writeFlyRC(homeDirectory, server.URL())
			Expect(err).ToNot(HaveOccurred())

			pathToFakeFly, err = gexec.Build(fmt.Sprintf("github.com/suhlig/apron-bus/fakes/fly%v", serverVersion))
			Expect(err).ToNot(HaveOccurred())

			command.Env = []string{
				fmt.Sprintf("PATH=%v", filepath.Dir(pathToFakeFly)),
				fmt.Sprintf("HOME=%v", homeDirectory),
			}

			args = []string{"--target", "mock", "--verbose", "status"}
		})

		AfterEach(func() {
			server.Close()
		})

		It("succeeds", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("exits with code zero", func() {
			Eventually(session).Should(gexec.Exit(0))
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
