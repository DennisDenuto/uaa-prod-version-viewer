package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"code.cloudfoundry.org/lager"
	"github.com/DennisDenuto/uaa-prod-version-viewer/printer"
	"golang.org/x/oauth2"
	"github.com/DennisDenuto/uaa-prod-version-viewer/pcfnotes"
	"github.com/DennisDenuto/uaa-prod-version-viewer/uaa_release"
	"os"
	"github.com/google/go-github/github"
	"context"
	"encoding/json"
	"io/ioutil"
)

var githubPAT string
var colorText string
var numPCFVersions int
var LineItemsLocalCacheLocation = "/tmp/local-bitbar-lineitems.json"

// bitbarCmd represents the bitbar command
var bitbarCmd = &cobra.Command{
	Use:   "bitbar",
	Short: "generate bitbar output",
	Run: func(cmd *cobra.Command, args []string) {
		logFile, err := os.OpenFile("/tmp/pcf-viewer.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
		if err != nil {
			panic(err)
		}

		logger := lager.NewLogger("pcf-version-viewer")
		logger.RegisterSink(lager.NewPrettySink(logFile, lager.DEBUG))

		pcfBitbarPrinter := printer.BitBarPrinter{
			Writer:           os.Stdout,
			StatusIconBase64: `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`,
		}

		lineItems := getLineItems(logger, githubPAT, numPCFVersions)
		cacheLineItems(logger, lineItems)

		err = pcfBitbarPrinter.Print(lineItems)
		if err != nil {
			logger.Fatal("Unable to print line items to bitbar", err)
		}

	},
}

func cacheLineItems(logger lager.Logger, lineItems []printer.LineItem) {
	lineItemsJson, err := json.Marshal(lineItems)
	if err != nil {
		logger.Fatal("Unable to marshal line items to json", err)
	}
	err = ioutil.WriteFile(LineItemsLocalCacheLocation, lineItemsJson, os.ModePerm)
	if err != nil {
		logger.Fatal("Unable to save json to file", err)
	}
}

func init() {
	bitbarCmd.Flags().StringVarP(&githubPAT, "token", "t", "", "Github PAT")
	bitbarCmd.Flags().StringVarP(&colorText, "color", "c", "black", "bitbar text color")
	bitbarCmd.Flags().IntVarP(&numPCFVersions, "num", "n", 3, "number of pcf versions to fetch")

	rootCmd.AddCommand(bitbarCmd)
}

func getLineItems(logger lager.Logger, githubPAT string, NPcfVersions int) []printer.LineItem {
	githubClient := github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: githubPAT,
	})))

	pcfVersionClient := pcfnotes.PcfVersion{
		Client: githubClient,
		Logger: logger,
	}
	uaaReleaseRepoClient := uaa_release.UAAReleaseRepo{
		Client: githubClient,
	}

	pcfNVersions, err := pcfVersionClient.LatestN(NPcfVersions)
	if err != nil {
		logger.Fatal("Failed getting latest N pcf versions", err)
	}

	pcfNotesComponentDetailsClient := pcfnotes.MustNewPcfNotesComponentDetails("https://releng.ci.cf-app.com/")
	uaaBoshReleasePCFComponentName := "uaa"

	var lineItems = make([]printer.LineItem, 0)
	for _, pcfVersion := range pcfNVersions {
		if found, pcfComponentDetails := pcfNotesComponentDetailsClient.ByName(uaaBoshReleasePCFComponentName, pcfVersion); found {
			logger.Debug(fmt.Sprintf("%v", pcfComponentDetails))
			uaaVersion := uaaReleaseRepoClient.GetUAAVersion(pcfComponentDetails.Version)
			uaaKeyRotatorVersion := uaaReleaseRepoClient.GetUAAKeyRotatorVersion(pcfComponentDetails.Version)

			lineItems = append(lineItems, toLineItem(pcfVersion, pcfComponentDetails, uaaVersion, uaaKeyRotatorVersion))
		}

	}
	return lineItems
}

func toLineItem(version pcfnotes.Version, pcfComponentDetails pcfnotes.ComponentDetails, uaaGitVersion, uaaKeyRotatorVersion string) printer.LineItem {
	return printer.LineItem{
		BoshRelease: printer.BoshRelease{
			Name:    pcfComponentDetails.Name + "-release",
			Version: pcfComponentDetails.Version,
			Packages: []printer.BoshPackage{
				{Name: "uaa", Version: uaaGitVersion},
				{Name: "uaa-key-rotator", Version: uaaKeyRotatorVersion},
			},
		},
		PcfVersion: version,
	}
}
