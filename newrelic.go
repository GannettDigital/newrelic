package newrelic

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/neocortical/newrelic/model"
)

// Logger is the logger used by this package. Set to a custom logger if needed.
var Logger = log.New(os.Stderr, "", log.LstdFlags)

// Verbose can be set globally to produce verbose log messages
var Verbose bool

const (
	// DefaultPollInterval is the recommended poll interval for NewRelic plugins
	DefaultPollInterval = time.Second * 60
)

const (
	agentVersion = "0.0.1"
	apiEndpoint  = "https://platform-api.newrelic.com/platform/v1/metrics"
)

type Client struct {
	License      string
	PollInterval time.Duration
	Plugins      []*Plugin

	agent        model.Agent
	lastPollTime time.Time
	url          string
	client       *http.Client
}

func (c *Client) AddPlugin(p *Plugin) {
	c.Plugins = append(c.Plugins, p)
}

func New(license string) *Client {
	result := &Client{
		License:      license,
		PollInterval: DefaultPollInterval,
		url:          apiEndpoint,
		client:       &http.Client{},
	}

	result.agent.Version = agentVersion
	result.agent.PID = os.Getpid()
	var err error
	if result.agent.Host, err = os.Hostname(); err != nil {
		panic(err)
	}

	return result
}

func (c *Client) doSend(t time.Time) bool {
	request, err := c.generateRequest(t)
	if err != nil {
		Logger.Printf("ERROR: encountered error(s) creating request data: %v", err)
	}

	responseCode := doPost(request, c.url, c.License, c.client)
	switch responseCode {
	case http.StatusOK:
		c.clearState(t)
	case http.StatusBadRequest:
		logResponseError(responseCode)
	case http.StatusForbidden:
		logResponseError(responseCode)
		return true
	case http.StatusNotFound:
		logResponseError(responseCode)
	case http.StatusMethodNotAllowed:
		// won't happen
		logResponseError(responseCode)
	case http.StatusRequestEntityTooLarge:
		// TODO: detect and split large responses
		logResponseError(responseCode)
	case http.StatusInternalServerError:
		logResponseError(responseCode)
	case http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		logResponseError(responseCode)
	}

	return false
}

func (c *Client) clearState(t time.Time) {
	c.lastPollTime = t
	for _, p := range c.Plugins {
		p.clearState()
	}
}

func (c *Client) Run() {
	if Verbose {
		Logger.Printf("Starting NewRelic plugin client...")
	}
	go c.run()
}

func (c *Client) run() {
	var fatal bool
	ticks := time.Tick(time.Duration(c.PollInterval))
	for t := range ticks {
		fatal = c.doSend(t)

		if fatal {
			Logger.Printf("ERROR: NewRelic plugin encountered a fatal error and is shutting down.")
			return
		}
	}
}

func doPost(request model.Request, url, license string, client *http.Client) int {
	var jsonBytes []byte
	var err error
	if Verbose {
		jsonBytes, err = json.MarshalIndent(request, "", "   ")
	} else {
		jsonBytes, err = json.Marshal(request)
	}
	if err != nil {
		Logger.Printf("error encoding json request: %v", err)
		return http.StatusBadRequest
	}

	if Verbose {
		Logger.Printf("Posting request:\n%s", string(jsonBytes))
	}

	httpRequest, err := http.NewRequest("POST", url, strings.NewReader(string(jsonBytes)))
	if err != nil {
		Logger.Printf("error creating request: %v", err)
		return http.StatusBadRequest
	}

	httpRequest.Header.Set("X-License-Key", license)
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Accept", "application/json")

	httpResponse, err := client.Do(httpRequest)
	if err != nil {
		Logger.Printf("error posting request: %v", err)
		return http.StatusServiceUnavailable
	}
	defer httpResponse.Body.Close()
	return httpResponse.StatusCode
}

func logResponseError(responseCode int) {
	Logger.Printf("ERROR: newrelic encountered %d response", responseCode)
}

func (c *Client) generateRequest(t time.Time) (request model.Request, err error) {
	request.Agent = c.agent

	var duration time.Duration
	if c.lastPollTime.IsZero() {
		duration = c.PollInterval
	} else {
		duration = t.Sub(c.lastPollTime)
	}

	for _, p := range c.Plugins {
		pluginSnapshot, cerr := p.generatePluginSnapshot(duration)

		// we are tolerant of request generation errors and should be able to recover
		if cerr != nil {
			err = accumulateErrors(err, cerr)
		}
		request.Plugins = append(request.Plugins, pluginSnapshot)
	}

	return request, err
}
