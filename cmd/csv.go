package cmd

import (
	"github.com/spf13/cobra"
	"code.cloudfoundry.org/lager"
	"github.com/DennisDenuto/uaa-prod-version-viewer/printer"
	"os"
	"github.com/skratchdot/open-golang/open"
	"github.com/cznic/fileutil"
)

var csvCmd = &cobra.Command{
	Use:   "csv",
	Short: "generate csv file of pcf versions",
	Run: func(cmd *cobra.Command, args []string) {
		logFile, err := os.OpenFile("/tmp/pcf-viewer.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
		if err != nil {
			panic(err)
		}

		logger := lager.NewLogger("pcf-version-viewer")
		logger.RegisterSink(lager.NewPrettySink(logFile, lager.DEBUG))

		csvFile, err := fileutil.TempFile(os.TempDir(), "pcf_notes", ".csv")
		if err != nil {
			logger.Fatal("Unable to create a temp file to write csv info into", err)
		}

		csvPrinter := printer.CSVPrinter{
			Writer: csvFile,
		}

		err = csvPrinter.Print(getLineItems(logger, githubPAT, numPCFVersions))
		if err != nil {
			logger.Fatal("Unable to print line items to csv file", err)
		}

		err = open.Run(csvFile.Name())
		if err != nil {
			logger.Fatal("Unable to open csv file", err)
		}

	},
}

func init() {
	csvCmd.Flags().StringVarP(&githubPAT, "token", "t", "", "Github PAT")
	csvCmd.Flags().IntVarP(&numPCFVersions, "num", "n", 3, "number of pcf versions to fetch")

	rootCmd.AddCommand(csvCmd)
}
