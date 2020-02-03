package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"gopkg.in/macaron.v1"
	log "unknwon.dev/clog/v2"
)

var Version = "dev"

func main() {
	configPath := flag.String("config", "./app.ini", "The config file path")
	flag.Parse()

	if err := log.NewConsole(); err != nil {
		panic("error init logger: " + err.Error())
	}
	defer log.Stop()

	config, err := loadConfig(*configPath)
	if err != nil {
		log.Fatal("Failed to load config: %v", err)
	}

	if config.Log.Path != "" {
		log.NewFile(log.FileConfig{
			Level:    log.LevelInfo,
			Filename: config.Log.Path,
			FileRotationConfig: log.FileRotationConfig{
				Rotate: true,
				Daily:  true,
			},
		})
	}

	log.Info("buildkite-flaky-reporter: %v", Version)

	slack := newSlackClient()
	buildkite := newBuildkiteClient(
		config.Buildkite.Token,
		config.Buildkite.OrgSlug,
		config.Buildkite.PipelineSlug,
	)
	_ = buildkite

	store := newStore()

	m := macaron.New()
	m.Use(macaron.Renderer())
	m.Post("/", func(c *macaron.Context) {
		// Discard requests that are not from Buildkite.
		ua := c.Req.Header.Get("User-Agent")
		if ua != "Buildkite-Request" {
			log.Trace("Request from %s has 'User-Agent'=%q discarded", c.RemoteAddr(), ua)
			return
		}

		// Discard events that are not "job.finished".
		event := c.Req.Header.Get("X-Buildkite-Event")
		if event != "job.finished" {
			log.Trace("Request from %s has 'X-Buildkite-Event'=%q discarded", c.RemoteAddr(), event)
			return
		}

		type JobEvent struct {
			Job struct {
				ID             string `json:"id"`
				Name           string `json:"name"`
				State          string `json:"state"`
				WebURL         string `json:"web_url"`
				Retried        bool   `json:"retried"`
				RetriedInJobID string `json:"retried_in_job_id"`
				RetriedCount   int    `json:"retries_count"`
			} `json:"job"`
			Build struct {
				ID            string `json:"id"`
				WebURL        string `json:"web_url"`
				Number        int    `json:"number"`
				Branch        string `json:"branch"`
				TriggeredFrom struct {
					BuildID           string `json:"build_id"`
					BuildNumber       int    `json:"build_number"`
					BuildPipelineSlug string `json:"build_pipeline_slug"`
				} `json:"triggered_from"`
			} `json:"build"`
		}
		var jobEvent JobEvent

		if err := json.NewDecoder(c.Req.Body().ReadCloser()).Decode(&jobEvent); err != nil {
			log.Error("Failed to decode request body: %v", err)
			return
		}

		// Discard builds that are not interested.
		if jobEvent.Build.Branch != config.Buildkite.BuildBranch {
			log.Trace("Request from %s has 'Build.Branch'=%q discarded", c.RemoteAddr(), jobEvent.Build.Branch)
			return
		}

		// Discard jobs that are not interested.
		if jobEvent.Job.Name != config.Buildkite.JobName {
			log.Trace("Request from %s has 'Job.Name'=%q discarded", c.RemoteAddr(), jobEvent.Job.Name)
			return
		}

		log.Info(spew.Sdump(jobEvent))
		fmt.Println("------------------------------------")

		c.Status(http.StatusNoContent)

		// Send notification if the job is canceled.
		if config.Buildkite.ReportCancel && jobEvent.Job.State == "canceled" {
			// TODO
			buildLink := fmt.Sprintf("<%s|%s#%d>",
				jobEvent.Build.WebURL,
				config.Buildkite.PipelineSlug,
				jobEvent.Build.Number,
			)
			triggerLink := fmt.Sprintf("<https://buildkite.com/%[1]s/%[2]s/builds/%[3]d|%[2]s (%[4]s) #%[3]d>",
				config.Buildkite.OrgSlug,
				jobEvent.Build.TriggeredFrom.BuildPipelineSlug,
				jobEvent.Build.TriggeredFrom.BuildNumber,
				jobEvent.Build.Branch,
			)
			msg := fmt.Sprintf(":warning: %s was canceled for %s", buildLink, triggerLink)
			if err = slack.send(config.Slack.CancelNotifyURL, msg); err != nil {
				log.Error("Failed to send Slack message: %v", err)
			}
			return
		}

		// Warning when failures exceeds threshold.
		if jobEvent.Job.State == "failed" {
			info := store.push(jobEvent.Build.Number, jobEvent.Job.WebURL)
			if jobEvent.Job.RetriedCount+1 >= config.Buildkite.FailuresThreshold {
				var buf bytes.Buffer
				buf.WriteString(fmt.Sprintf("%s test failed %d times in a row on %q!\n", jobEvent.Job.Name, jobEvent.Job.RetriedCount+1, jobEvent.Build.Branch))
				for _, url := range info.failedURLs {
					buf.WriteString("- " + url + "\n")
				}

				if err = slack.send(config.Slack.URL, buf.String()); err != nil {
					log.Error("Failed to send Slack message: %v", err)
				}
			}
			return
		}

		// Mark as flaky test when passed after some failures.
		if jobEvent.Job.RetriedCount > 0 {
			var buf bytes.Buffer
			buf.WriteString(fmt.Sprintf("%s found flaky tests on %q!\n", jobEvent.Job.Name, jobEvent.Build.Branch))

			info := store.get(jobEvent.Build.Number)
			if info != nil {
				for _, url := range info.failedURLs {
					buf.WriteString("- " + url + "\n")
				}
			} else {
				buf.WriteString("No history found, showing passed: " + jobEvent.Job.WebURL)
			}

			if err = slack.send(config.Slack.URL, buf.String()); err != nil {
				log.Error("Failed to send Slack message: %v", err)
			}
		}
	})

	log.Info("Listening on http://%s...", config.Addr)
	if err := http.ListenAndServe(config.Addr, m); err != nil {
		log.Fatal("Failed to start server: %v", err)
	}
}
