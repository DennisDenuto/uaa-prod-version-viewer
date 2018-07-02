package printer_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/DennisDenuto/uaa-prod-version-viewer/printer"
	. "github.com/onsi/gomega/gbytes"
	"github.com/DennisDenuto/uaa-prod-version-viewer/pcfnotes"
	"io"
)

var _ = Describe("Printer", func() {

	Describe("bitbar", func() {
		var bitbarPrinter printer.BitBarPrinter
		var uaaReleaseLineItem1, uaaReleaseLineItem2 printer.LineItem
		var writer io.Writer

		BeforeEach(func() {
			writer = NewBuffer()
			uaaReleaseLineItem1 = printer.LineItem{
				PcfVersion: pcfnotes.Version(2.2),
				BoshRelease: printer.BoshRelease{
					Name:    "uaa-release",
					Version: "uaa-version",
					Packages: []printer.BoshPackage{
						{
							Name:    "uaa",
							Version: "standalone-version",
						},
					},
				},
			}

			uaaReleaseLineItem2 = printer.LineItem{
				PcfVersion: pcfnotes.Version(2.3),
				BoshRelease: printer.BoshRelease{
					Name:    "uaa-release",
					Version: "uaa-version2",
					Packages: []printer.BoshPackage{
						{
							Name:    "uaa",
							Version: "standalone-version2",
						},
					},
				},
			}

			bitbarPrinter = printer.BitBarPrinter{
				Writer: writer,
				StatusIconBase64: "image_base64",
			}
		})

		It("should print a bitbar menu", func() {
			err := bitbarPrinter.Print([]printer.LineItem{uaaReleaseLineItem1, uaaReleaseLineItem2})
			Expect(err).NotTo(HaveOccurred())

			Expect(writer).Should(Say(`image=image_base64`))
			Expect(writer).Should(Say(`---`))
			Expect(writer).Should(Say(`PCF 2.2`))
			Expect(writer).Should(Say(`---`))
			Expect(writer).Should(Say(`uaa-release: uaa-version \| color=black font=UbuntuMono-Bold size=12 href=https://github.com/cloudfoundry/uaa-release/tree/uaa-version`))
			Expect(writer).Should(Say(`-- uaa: version standalone-version \| href=https://github.com/cloudfoundry/uaa/tree/standalone-version`))

			Expect(writer).Should(Say(`---`))
			Expect(writer).Should(Say(`PCF 2.3`))
			Expect(writer).Should(Say(`---`))
			Expect(writer).Should(Say(`uaa-release: uaa-version2 \| color=black font=UbuntuMono-Bold size=12 href=https://github.com/cloudfoundry/uaa-release/tree/uaa-version2`))
			Expect(writer).Should(Say(`-- uaa: version standalone-version2 \| href=https://github.com/cloudfoundry/uaa/tree/standalone-version2`))
		})
	})
})
