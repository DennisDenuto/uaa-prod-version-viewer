package uaa_release

import (
	"github.com/google/go-github/github"
	"context"
	"bufio"
	"strings"
)

type UAAReleaseRepo struct {
	Client *github.Client
}

func (u UAAReleaseRepo) GetUAAVersion(uaaReleaseBranch string) string {
	fileContent, _, _, err := u.Client.Repositories.GetContents(context.Background(), "cloudfoundry", "uaa-release", "src/uaa", &github.RepositoryContentGetOptions{Ref: uaaReleaseBranch})
	if err != nil {
		panic(err)
	}
	uaaSHA := fileContent.GetSHA()

	gradlePropertiesContent, _, _, err := u.Client.Repositories.GetContents(context.Background(), "cloudfoundry", "uaa", "gradle.properties", &github.RepositoryContentGetOptions{Ref: uaaSHA})
	if err != nil {
		panic(err)
	}

	gradleContents, err := gradlePropertiesContent.GetContent()
	if err != nil {
		panic(err)
	}

	gradleProperties, err := ReadProperties(gradleContents)
	if err != nil {
		panic(err)
	}
	return gradleProperties["version"]
}

func (u UAAReleaseRepo) GetUAAKeyRotatorVersion(uaaReleaseBranch string) string {
	fileContent, _, _, err := u.Client.Repositories.GetContents(context.Background(), "cloudfoundry", "uaa-release", "src/github.com/cloudfoundry/uaa-key-rotator", &github.RepositoryContentGetOptions{Ref: uaaReleaseBranch})
	if fileContent == nil {
		return "n/a"
	}

	if err != nil {
		panic(err)
	}
	uaaKeyRotatorSHA := fileContent.GetSHA()
	tags, _, err := u.Client.Repositories.ListTags(context.Background(), "cloudfoundry", "uaa-key-rotator", &github.ListOptions{PerPage: 1000})
	if err != nil {
		panic(err)
	}

	for _, tag := range tags {
		if tag.Commit.GetSHA() == uaaKeyRotatorSHA {
			return tag.GetName()
		}
	}
	return uaaKeyRotatorSHA
}

type AppConfigProperties map[string]string

func ReadProperties(propertiesContent string) (AppConfigProperties, error) {
	config := AppConfigProperties{}

	scanner := bufio.NewScanner(strings.NewReader(propertiesContent))
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				config[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return config, nil
}
