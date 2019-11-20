package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type slackClient struct {
	url string
}

func newSlackClient(url string) *slackClient {
	return &slackClient{
		url: url,
	}
}

func (c *slackClient) send(text string) error {
	if c.url == "" {
		return nil
	}

	resp, err := http.DefaultClient.Post(c.url, "application/json",
		strings.NewReader(
			fmt.Sprintf(`{"text": %s}`, strconv.Quote(text)),
		))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	p, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if string(p) != "ok" {
		return errors.New(string(p))
	}

	return nil
}
