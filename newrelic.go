package newrelic

import (
	"encoding/json"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/neocortical/newrelic/model"
)

const (
	// DefaultPollInterval is the recommended poll interval for NewRelic plugins
	DefaultPollInterval = time.Minute
)

const (
	agentVersion = "0.0.1"
	apiEndpoint  = "https://platform-api.newrelic.com/platform/v1/metrics"
)

var netTransport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 5 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 5 * time.Second,
}

var netClient = &http.Client{
	Timeout:   time.Second * 5,
	Transport: netTransport,
}

// Client encapsulates a NewRelic plugin client and all the plugins it reports
type Client struct {
	License      string
	PollInterval time.Duration
	Plugins      []*Plugin

	// HTTPClient is exposed to allow users to configure proxies, etc.
	HTTPClient *http.Client

	agent        model.Agent
	lastPollTime time.Time
	url          string
}

// AddPlugin appends a plugin to a clients list of plugins. A plugin is a "component"
// in the API call and can be configured (with a unique GUID) in the NewRelic UI.
func (c *Client) AddPlugin(p *Plugin) {
	c.Plugins = append(c.Plugins, p)
}

// New creates a new Client with the given license
func New(license string) *Client {
	result := &Client{
		License:      license,
		PollInterval: DefaultPollInterval,
		HTTPClient:   netClient,
		url:          apiEndpoint,
	}

	result.agent.Version = agentVersion
	result.agent.PID = os.Getpid()
	var err error
	if result.agent.Host, err = os.Hostname(); err != nil {
		panic(err)
	}

	return result
}

func (c *Client) doSend(t time.Time) {
	request, err := c.generateRequest(t)
	if err != nil {
		Log(LogError, "ERROR: encountered error(s) creating request data: %v", err)
	}
	c.lastPollTime = t

	responseCode := doPost(request, c.url, c.License, c.HTTPClient)
	switch responseCode {
	case http.StatusOK:
		c.clearState()
	case http.StatusBadRequest:
		logResponseError(responseCode)
	case http.StatusForbidden:
		logResponseError(responseCode)
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
	case http.StatusTeapot:
		Log(LogError, "Server is a teapot!")
	}
}

func (c *Client) clearState() {
	for _, p := range c.Plugins {
		p.clearState()
	}
}

// Run starts the NewRelic client asynchronously. Do not alter the configuration
// of plugins after starting the client, as this creates race conditions.
func (c *Client) Run() {
	Log(LogInfo, "Starting NewRelic plugin client...")
	go c.run()
}

func (c *Client) run() {
	ticks := time.Tick(time.Duration(c.PollInterval))
	for t := range ticks {
		c.doSend(t)
	}
}

func doPost(request model.Request, url, license string, client *http.Client) int {
	var jsonBytes []byte
	var err error
	if LogLevel <= LogDebug {
		jsonBytes, err = json.MarshalIndent(request, "", "   ")
	} else {
		jsonBytes, err = json.Marshal(request)
	}
	if err != nil {
		Log(LogError, "error encoding json request: %v", err)
		return http.StatusBadRequest
	}

	Log(LogDebug, "Posting request:\n%s", string(jsonBytes))

	httpRequest, err := http.NewRequest("POST", url, strings.NewReader(string(jsonBytes)))
	if err != nil {
		Log(LogError, "error creating request: %v", err)
		return http.StatusBadRequest
	}

	httpRequest.Header.Set("X-License-Key", license)
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Accept", "application/json")

	httpResponse, err := client.Do(httpRequest)
	if err != nil {
		Log(LogError, "error posting request: %v", err)
		return http.StatusServiceUnavailable
	}
	defer httpResponse.Body.Close()
	return httpResponse.StatusCode
}

func logResponseError(responseCode int) {
	Log(LogError, "ERROR: newrelic encountered %d response", responseCode)
}

func (c *Client) generateRequest(t time.Time) (request model.Request, err CompositeError) {
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
			err = err.Accumulate(cerr)
		}
		request.Plugins = append(request.Plugins, pluginSnapshot)
	}

	return request, err
}
