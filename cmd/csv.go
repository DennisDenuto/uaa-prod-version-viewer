package cmd

import (
	"code.cloudfoundry.org/lager"
	"encoding/json"
	"github.com/DennisDenuto/uaa-prod-version-viewer/printer"
	"github.com/cznic/fileutil"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
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

		var csvFile *os.File
		if stdOut {
			csvFile = os.Stdout
		} else {
			csvFile, err = fileutil.TempFile(os.TempDir(), "pcf_notes", ".csv")
		}
		if err != nil {
			logger.Fatal("Unable to create a temp file to write csv info into", err)
		}

		csvPrinter := printer.CSVPrinter{
			Writer: csvFile,
		}

		var lineItems []printer.LineItem
		if cached {
			lineItems = getCachedLineItemsFromFile(logger)
		} else {
			lineItems = getLineItems(logger, githubPAT, numPCFVersions)
		}
		err = csvPrinter.Print(lineItems)
		if err != nil {
			logger.Fatal("Unable to print line items to csv file", err)
		}

		if !stdOut {
			err = open.Run(csvFile.Name())
			if err != nil {
				logger.Fatal("Unable to open csv file", err)
			}
		}

	},
}

func getCachedLineItemsFromFile(logger lager.Logger) []printer.LineItem {
	lineItemsJson, err := ioutil.ReadFile(LineItemsLocalCacheLocation)
	if err != nil {
		logger.Fatal("Unable to unmarshal line items from json file", err)
	}

	var lineItems []printer.LineItem
	err = json.Unmarshal(lineItemsJson, &lineItems)
	if err != nil {
		logger.Fatal("Unable to unmarshal line items", err)
	}

	return lineItems
}

func init() {
	csvCmd.Flags().StringVarP(&githubPAT, "token", "t", "", "Github PAT")
	csvCmd.Flags().IntVarP(&numPCFVersions, "num", "n", 3, "number of pcf versions to fetch")
	csvCmd.Flags().BoolVarP(&cached, "cached", "c", true, "should pcf info be read from file")
	csvCmd.Flags().BoolVarP(&stdOut, "stdout", "o", false, "file to write csv to")

	rootCmd.AddCommand(csvCmd)
}
