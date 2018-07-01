package main

import (
	"github.com/google/go-github/github"
	"github.com/DennisDenuto/uaa-prod-version-viewer/pcfnotes"
	"code.cloudfoundry.org/lager"
	"fmt"
	"os"
	"net/url"
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
	for _, version := range versions {

		if found, componentDetails := details.ByName("uaa", version); found {
			logger.Debug(fmt.Sprintf("hello? %v", componentDetails))
		}

	}


	//client.Repositories.ListBranches()
	//fileContent, directoryContent, resp, err := client.Repositories.GetContents(context.Background(), "pivotal-cf", "pcf-release-notes", "runtime-rn.html.md.erb", &github.RepositoryContentGetOptions{Ref: "2.2"})

}
