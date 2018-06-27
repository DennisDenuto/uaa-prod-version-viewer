package pcfnotes

import (
	"github.com/google/go-github/github"
	"context"
	"strconv"
	"log"
	"sort"
	"github.com/pkg/errors"
)

type Version float64

type Versions interface {
	Latest() Version
	LatestN(head int) []Version
}


type PcfVersion struct {
	Client *github.Client
	Logger *log.Logger
}

func (p PcfVersion) Latest() Version {
	versions, err := p.LatestN(1)
	if err != nil {
		panic(err)
	}
	return Version(versions[len(versions) - 1])
}

func (p PcfVersion) LatestN(head int) ([]Version, error) {
	branches, _, err := p.Client.Repositories.ListBranches(context.Background(), "pivotal-cf", "pcf-release-notes", &github.ListOptions{})

	if err != nil {
		return nil, errors.Wrap(err, "Unable to fetch branches from github")
	}

	var versions []Version
	for _, branch := range branches {
		branchVersion, err := strconv.ParseFloat(*branch.Name, 64)
		if err != nil {
			p.Logger.Printf("branch name (%v) is not a valid. Skipping", *branch.Name)
			continue
		}
		versions = append(versions, Version(branchVersion))
	}
	sort.Sort(versionSlice(versions))

	return versions[len(versions) - head:], nil
}


type versionSlice []Version

func (p versionSlice) Len() int           { return len(p) }
func (p versionSlice) Less(i, j int) bool { return p[i] < p[j] || isNaN(float64(p[i])) && !isNaN(float64(p[j])) }
func (p versionSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// isNaN is a copy of math.IsNaN to avoid a dependency on the math package.
func isNaN(f float64) bool {
	return f != f
}
