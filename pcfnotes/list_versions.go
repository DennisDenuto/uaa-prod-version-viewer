package pcfnotes

import (
	"github.com/google/go-github/github"
	"context"
	"strconv"
	"sort"
	"github.com/pkg/errors"
	"code.cloudfoundry.org/lager"
	"fmt"
	"strings"
)

type Version struct {
	Major int
	Minor int
}

func (v Version) String() string {
	majorVersion := strconv.Itoa(v.Major)
	minorVersion := strconv.Itoa(v.Minor)
	return majorVersion + "." + minorVersion
}

type Versions interface {
	Latest() Version
	LatestN(head int) []Version
}

type PcfVersion struct {
	Client *github.Client
	Logger lager.Logger
}

func (p PcfVersion) Latest() Version {
	versions, err := p.LatestN(1)
	if err != nil {
		panic(err)
	}
	return Version(versions[len(versions)-1])
}

func (p PcfVersion) LatestN(head int) ([]Version, error) {
	branches, _, err := p.Client.Repositories.ListBranches(context.Background(), "pivotal-cf", "pcf-release-notes", &github.ListOptions{})

	if err != nil {
		return nil, errors.Wrap(err, "Unable to fetch branches from github")
	}

	var versions []Version
	for _, branch := range branches {
		uaaReleaseVersion, err := toVersion(*branch.Name)
		if err != nil {
			p.Logger.Info(fmt.Sprintf("branch name (%v) is not a valid. Skipping", *branch.Name))
			continue
		}
		versions = append(versions, uaaReleaseVersion)
	}
	sort.Sort(versionSlice(versions))

	return versions[len(versions)-head:], nil
}

func toVersion(branch string) (Version, error) {
	branchVersionSplit := strings.Split(branch, ".")

	majorVersion, err := strconv.Atoi(branchVersionSplit[0])
	if err != nil {
		return Version{}, err
	}

	minorVersion, err := strconv.Atoi(branchVersionSplit[1])
	if err != nil {
		return Version{}, err
	}

	return Version{
		Major: majorVersion,
		Minor: minorVersion,
	}, nil
}

type versionSlice []Version

func (p versionSlice) Len() int           { return len(p) }
func (p versionSlice) Less(i, j int) bool { return p[i].Major < p[j].Major && p[i].Minor < p[j].Minor }
func (p versionSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
