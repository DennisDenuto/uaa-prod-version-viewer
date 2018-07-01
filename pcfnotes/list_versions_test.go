package pcfnotes_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/DennisDenuto/uaa-prod-version-viewer/pcfnotes"
	"github.com/onsi/gomega/ghttp"
	"github.com/google/go-github/github"
	. "github.com/onsi/gomega"
	"fmt"
	"code.cloudfoundry.org/lager/lagertest"
)

type GithubBranch struct {
	Name string `json:"name,omitempty"`
}

var _ = Describe("ListVersions", func() {
	var branchesUrl = fmt.Sprintf("/repos/%v/%v/branches", "pivotal-cf", "pcf-release-notes")
	var pcfVersion pcfnotes.PcfVersion
	var githubServer *ghttp.Server

	BeforeEach(func() {
		githubServer = ghttp.NewServer()

		githubServer.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", branchesUrl),
			ghttp.RespondWithJSONEncoded(200, []GithubBranch{
				{"master"}, {"abranch"}, {"1.0"}, {"1.1"}, {"2.0"}, {"3.0"},
			}),
		))

		githubClient, err := github.NewEnterpriseClient(githubServer.URL(), githubServer.URL(), nil)
		Expect(err).NotTo(HaveOccurred())

		pcfVersion = pcfnotes.PcfVersion{
			Client: githubClient,
			Logger: lagertest.NewTestLogger("list-versions"),
		}
	})

	AfterEach(func() {
		githubServer.Close()
	})

	It("should return the latest pcf version from the pcfnotes repo", func() {
		latestPcfVersion := pcfVersion.Latest()
		Expect(githubServer.ReceivedRequests(), HaveLen(1))
		Expect(latestPcfVersion).To(Equal(pcfnotes.Version(3.0)))
	})

	It("should return the latest 2 pcf version from the pcfnotes repo", func() {
		latestPcfVersion, err := pcfVersion.LatestN(2)
		Expect(err).NotTo(HaveOccurred())
		Expect(githubServer.ReceivedRequests(), HaveLen(1))
		Expect(latestPcfVersion).To(ConsistOf(pcfnotes.Version(2.0), pcfnotes.Version(3.0)))
	})

	Context("github api returns an error", func() {
		BeforeEach(func() {
			githubServer.SetHandler(0, ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", branchesUrl),
				ghttp.RespondWith(500, nil),
			))
		})

		It("should return an error with a message", func() {
			_, err := pcfVersion.LatestN(1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("Unable to fetch branches from github: GET .*: 500"))
		})
	})
})
