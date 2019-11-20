package main

import (
	"gopkg.in/ini.v1"
)

type config struct {
	Addr      string
	Buildkite struct {
		Token             string
		OrgSlug           string
		PipelineSlug      string
		BuildBranch       string
		JobName           string
		FailuresThreshold int
	}
	Slack struct {
		URL string `ini:"url"`
	}
}

func loadConfig(path string) (*config, error) {
	var c config
	return &c, ini.MapToWithMapper(&c, ini.TitleUnderscore, path)
}
