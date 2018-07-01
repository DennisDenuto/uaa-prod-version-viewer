package main

import (
	"github.com/google/go-github/github"
	"github.com/DennisDenuto/uaa-prod-version-viewer/pcfnotes"
	"code.cloudfoundry.org/lager"
	"fmt"
	"os"
	"net/url"
	"github.com/DennisDenuto/uaa-prod-version-viewer/uaa_release"
)

func main() {
	client := github.NewClient(nil)
	logger := lager.NewLogger("prod-version-viewer")
	logger.RegisterSink(lager.NewPrettySink(os.Stdout, lager.DEBUG))

	version := pcfnotes.PcfVersion{
		Client: client,
		Logger: logger,
	}

	versions, err := version.LatestN(3)
	if err != nil {
		panic(err)
	}

	relengURL, err := url.Parse("https://releng.ci.cf-app.com/")
	if err != nil {
		panic(err)
	}

	details := pcfnotes.PcfNotesComponentDetails{
		BaseURL: relengURL,
	}

	uaaReleaseRepo := uaa_release.UAAReleaseRepo{
		Client: client,
	}

	for _, version := range versions {

		if found, pcfComponentDetails := details.ByName("uaa", version); found {
			logger.Debug(fmt.Sprintf("%v", pcfComponentDetails))
			uaaVersion := uaaReleaseRepo.GetUAAVersion(pcfComponentDetails.Version)
			println(uaaVersion)
		}

	}



}
