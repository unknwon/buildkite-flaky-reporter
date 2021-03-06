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
		ReportCancel      bool
	}
	Slack struct {
		URL             string `ini:"url"`
		CancelNotifyURL string `ini:"cancel_notify_url"`
	}
	Log struct {
		Path string
	}
}

func loadConfig(path string) (*config, error) {
	var c config
	return &c, ini.MapToWithMapper(&c, ini.TitleUnderscore, path)
}
