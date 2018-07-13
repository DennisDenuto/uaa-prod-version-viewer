package pcfnotes

import (
	"net/http"
	"net/url"
	"encoding/json"
	"strings"
	"fmt"
)

type ComponentDetails struct {
	Name       string
	Version    string
	PCFVersion string
}

type GetComponentDetails interface {
	ByName(releaseName string, version Version) ComponentDetails
}

type PcfNotesComponentDetails struct {
	BaseURL *url.URL
}

func MustNewPcfNotesComponentDetails(pcfURL string) PcfNotesComponentDetails {
	relengURL, err := url.Parse(pcfURL)
	if err != nil {
		panic(err)
	}

	return PcfNotesComponentDetails{
		BaseURL: relengURL,
	}
}

func (p PcfNotesComponentDetails) ByName(releaseName string, version Version) (bool, ComponentDetails) {
	pcfPipelineURI := fmt.Sprintf("/api/v1/teams/main/pipelines/build::%s/resources", version.String())
	pcfPipelineURL, err := p.BaseURL.Parse(pcfPipelineURI)
	if err != nil {
		return false, ComponentDetails{}
	}
	resp, err := http.Get(pcfPipelineURL.String())
	if err != nil {
		return false, ComponentDetails{}
	}

	pcfPipelineDecoder := json.NewDecoder(resp.Body)
	details := &PCFPipelineDetails{}

	err = pcfPipelineDecoder.Decode(details)
	if err != nil {
		return false, ComponentDetails{}
	}

	for _, detail := range *details {
		releaseWithDash := releaseName + "-"
		if strings.HasPrefix(detail.Name, releaseWithDash) && strings.HasSuffix(detail.Name, "-release") {
			versionStartIndex := strings.LastIndex(detail.Name, releaseWithDash) + len(releaseWithDash)
			versionEndIndex := strings.Index(detail.Name[versionStartIndex:], "-release")
			return true, ComponentDetails{
				Name:       detail.Name[:versionStartIndex-1],
				Version:    detail.Name[versionStartIndex:(versionStartIndex + versionEndIndex)],
				PCFVersion: version.String(),
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
