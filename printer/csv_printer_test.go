package printer

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"


	"github.com/DennisDenuto/uaa-prod-version-viewer/pcfnotes"
	"io"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("Printer", func() {
	Describe("CSV", func() {
		var csvPrinter CSVPrinter
		var uaaReleaseLineItem1, uaaReleaseLineItem2 LineItem
		var writer io.Writer

		BeforeEach(func() {
			writer = NewBuffer()
			uaaReleaseLineItem1 = LineItem{
				PcfVersion: pcfnotes.Version{2, 2},
				BoshRelease: BoshRelease{
					Name:    "uaa-release",
					Version: "uaa-version",
					Packages: []BoshPackage{
						{
							Name:    "uaa",
							Version: "standalone-version",
						},
					},
				},
			}

			uaaReleaseLineItem2 = LineItem{
				PcfVersion: pcfnotes.Version{2, 3},
				BoshRelease: BoshRelease{
					Name:    "uaa-release",
					Version: "uaa-version2",
					Packages: []BoshPackage{
						{
							Name:    "uaa",
							Version: "standalone-version2",
						},
						{
							Name:    "uaa-key-rotator",
							Version: "rotator-version2",
						},
					},
				},
			}

			csvPrinter = CSVPrinter{
				Writer:           writer,
			}
		})

		It("should print a bitbar menu", func() {
			err := csvPrinter.Print([]LineItem{uaaReleaseLineItem1, uaaReleaseLineItem2})
			Expect(err).NotTo(HaveOccurred())

			Expect(writer).Should(Say(`PCF Version,uaa-release,uaa,uaa-key-rotator`))
			Expect(writer).Should(Say(`2.2,uaa-version,standalone-version,`))
			Expect(writer).Should(Say(`2.3,uaa-version2,standalone-version2,rotator-version2`))
		})
	})

})
