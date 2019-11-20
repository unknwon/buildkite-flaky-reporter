package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type buildkiteClient struct {
	token        string
	orgSlug      string
	pipelineSlug string
}

func newBuildkiteClient(token, orgSlug, pipelineSlug string) *buildkiteClient {
	return &buildkiteClient{
		token:        token,
		orgSlug:      orgSlug,
		pipelineSlug: pipelineSlug,
	}
}

type buildkiteJob struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	State  string `json:"state"`
	WebURL string `json:"web_url"`
}

type buildkiteBuild struct {
	ID     string         `json:"id"`
	Number int            `json:"number"`
	Branch string         `json:"branch"`
	Jobs   []buildkiteJob `json:"jobs"`
}

func (c *buildkiteClient) getBuild(number int) (*buildkiteBuild, error) {
	var build buildkiteBuild
	err := c.get(
		fmt.Sprintf("/organizations/%s/pipelines/%s/builds/%d", c.orgSlug, c.pipelineSlug, number),
		&build,
	)
	if err != nil {
		return nil, err
	}
	return &build, nil
}

func (c *buildkiteClient) get(path string, v interface{}) error {
	req, err := http.NewRequest("GET", "https://api.buildkite.com/v2"+path, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+c.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}
	return nil
}
