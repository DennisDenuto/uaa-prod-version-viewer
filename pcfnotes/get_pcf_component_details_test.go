package pcfnotes_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/DennisDenuto/uaa-prod-version-viewer/pcfnotes"
	"io/ioutil"
	"net/url"
)

var _ = Describe("GetPcfComponentDetails", func() {

	var getComponentDetails pcfnotes.PcfNotesComponentDetails
	var server *ghttp.Server

	BeforeEach(func() {
		server = ghttp.NewServer()
		pcfPipelineResponse, err := ioutil.ReadFile("test_data/pipeline_resources_response_2_1.json")
		Expect(err).NotTo(HaveOccurred())
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/v1/teams/main/pipelines/build::2.1/resources"),
				ghttp.RespondWith(200, pcfPipelineResponse),
			),
		)

		testServerURI, err := url.Parse(server.URL())
		Expect(err).NotTo(HaveOccurred())
		getComponentDetails = pcfnotes.PcfNotesComponentDetails{
			BaseURL: testServerURI,
		}
	})

	It("should get uaa release by pcf version", func() {
		found, uaaPcfComponent := getComponentDetails.ByName("uaa", pcfnotes.Version(2.1))

		Expect(found).To(Equal(true))
		Expect(uaaPcfComponent.Name).To(Equal("uaa"))
		Expect(uaaPcfComponent.Version).To(Equal("v55"))
	})

	XContext("when pcf pipeline returns an error", func() {
		It("should return a meaningful error", func() {
			Fail("todo")
		})
	})
})
