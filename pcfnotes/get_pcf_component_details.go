package pcfnotes

import (
	"net/http"
	"net/url"
	"encoding/json"
	"strings"
	"fmt"
)

type ComponentDetails struct {
	Name    string
	Version string
}

type GetComponentDetails interface {
	ByName(releaseName string, version Version) ComponentDetails
}

type PcfNotesComponentDetails struct {
	BaseURL *url.URL
}

func (p PcfNotesComponentDetails) ByName(releaseName string, version Version) (bool, ComponentDetails) {
	pcfPipelineURI := fmt.Sprintf("/api/v1/teams/main/pipelines/build::%s/resources", version.String())
	pcfPipelineURL, err := p.BaseURL.Parse(pcfPipelineURI)
	if err != nil {
		panic(err)
	}
	println(pcfPipelineURL.String())
	resp, err := http.Get(pcfPipelineURL.String())
	if err != nil {
		panic(err)
	}

	pcfPipelineDecoder := json.NewDecoder(resp.Body)
	details := &PCFPipelineDetails{}

	err = pcfPipelineDecoder.Decode(details)
	if err != nil {
		panic(err)
	}

	for _, detail := range *details {
		releaseWithDash := releaseName + "-"
		if strings.HasPrefix(detail.Name, releaseWithDash) && strings.HasSuffix(detail.Name, "-release") {
			versionStartIndex := strings.LastIndex(detail.Name, releaseWithDash) + len(releaseWithDash)
			versionEndIndex := strings.Index(detail.Name[versionStartIndex:], "-release")
			return true, ComponentDetails{
				Name:    detail.Name[:versionStartIndex-1],
				Version: detail.Name[versionStartIndex:(versionStartIndex + versionEndIndex)],
			}
		}
	}

	return false, ComponentDetails{}
}

type PCFPipelineDetails []struct {
	Name           string `json:"name"`
	PipelineName   string `json:"pipeline_name"`
	TeamName       string `json:"team_name"`
	Type           string `json:"type"`
	LastChecked    int    `json:"last_checked"`
	FailingToCheck bool   `json:"failing_to_check,omitempty"`
}
