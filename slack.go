package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type slackClient struct{}

func newSlackClient() *slackClient {
	return &slackClient{}
}

func (c *slackClient) send(url, text string) error {
	resp, err := http.DefaultClient.Post(url, "application/json",

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
