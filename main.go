package main

import (
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
	port := flag.String("port", "4222", "The listening port number")
	flag.Parse()

	if err := log.NewConsole(); err != nil {
		panic("error init logger: " + err.Error())
	}
	defer log.Stop()

	log.Info("buildkite-flaky-reporter: %v", Version)

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
				ID     string `json:"id"`
				Name   string `json:"name"`
				State  string `json:"state"`
				WebURL string `json:"web_url"`
			} `json:"job"`
			Build struct {
				ID string `json:"id"`
			} `json:"build"`
		}
		var jobEvent JobEvent

		if err := json.NewDecoder(c.Req.Body().ReadCloser()).Decode(&jobEvent); err != nil {
			log.Error("Failed to decode request body: %v", err)
			return
		}

		// Discard jobs that are not ":chromium:".
		if jobEvent.Job.Name != ":chromium:" {
			log.Trace("Request from %s has 'Job.Name'=%q discarded", c.RemoteAddr(), jobEvent.Job.Name)
			return
		}

		spew.Dump(jobEvent)
		fmt.Println("------------------------------------")

		// TODO: send event info to Slack

		c.Status(http.StatusNoContent)
	})

	addr := "localhost:" + *port
	log.Info("Listening on http://%s...", addr)
	if err := http.ListenAndServe(addr, m); err != nil {
		log.Fatal("Failed to start server: %v", err)
	}
}
