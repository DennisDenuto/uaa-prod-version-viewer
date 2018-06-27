package main

import (
	"github.com/google/go-github/github"
	"context"
)

func main() {
	client := github.NewClient(nil)
	client.Repositories.ListBranches()
	closer, err := client.Repositories.DownloadContents(context.Background(), "pivotal-cf", "pcf-release-notes", "runtime-rn.html.md.erb", &github.RepositoryContentGetOptions{Ref: "2.2"})

}
