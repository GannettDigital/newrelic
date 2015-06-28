package newrelic

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/neocortical/newrelic/model"
)

// Logger is the logger used by this package. Set to a custom logger if needed.
var Logger = log.New(os.Stderr, "", log.LstdFlags)

// Verbose can be set globally to produce verbose log messages
var Verbose bool

var guidNormalizationRegexp = regexp.MustCompile(`[^a-zA-Z0-9\._]+`)

const (
	// DefaultPollInterval is the recommended poll interval for NewRelic plugins
	DefaultPollInterval = 60
)

const (
	agentGUID    = "net.neocortical.newrelic"
	agentVersion = "0.0.1"
	apiEndpoint  = "https://platform-api.newrelic.com/platform/v1/metrics"
)

type Plugin struct {
	Name         string
	License      string
	PollInterval int
	Components   []*Component

	agent        model.Agent
	lastPollTime time.Time
	url          string
	client       *http.Client
}

func (p *Plugin) AppendComponent(c *Component) {
	c.guid = generateComponentGUID(agentGUID, p.Name)
	p.Components = append(p.Components, c)
}

func NewPlugin(name, license string) *Plugin {
	result := &Plugin{
		Name:         name,
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

// Component encapsulates a component of a plugin
type Component struct {
	Name string

	guid     string
	duration int
	metrics  []metricsGroup
}

func (c *Component) AddMetric(metric Metric) {
	c.metrics = append(c.metrics, &simpleMetricsGroup{metric: metric})
}

func (p *Plugin) doSend() bool {
	t := time.Now()
	request, err := p.generateRequest(t)
	if err != nil {
		Logger.Printf("ERROR: encountered error(s) creating request data: %v", err)
	}

	responseCode := doPost(request, p.url, p.License, p.client)
	switch responseCode {
	case http.StatusOK:
		p.clearState()
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

func (p *Plugin) clearState() {
	for _, c := range p.Components {
		c.clearState()
	}
}

func (p *Plugin) Run() {
	Logger.Printf("Starting NewRelic plugin client %s...", p.Name)
	go p.run()
}

func (p *Plugin) run() {
	var fatal bool
	ticks := time.Tick(time.Duration(p.PollInterval) * time.Second)
	for _ = range ticks {
		fatal = p.doSend()

		if fatal {
			Logger.Printf("ERROR: NewRelic plugin %s encountered a fatal error and is shutting down.", p.Name)
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

func (p *Plugin) generateRequest(t time.Time) (request model.Request, err error) {
	request.Agent = p.agent

	var duration int
	if p.lastPollTime.IsZero() {
		duration = p.PollInterval
	} else {
		duration = int(t.Sub(p.lastPollTime).Seconds())
	}

	p.lastPollTime = t

	for _, component := range p.Components {
		componentRequest, cerr := component.generateComponentSnapshot(duration)

		// we are tolerant of request generation errors and should be able to recover
		if cerr != nil {
			err = accumulateErrors(err, cerr)
		}
		request.Components = append(request.Components, componentRequest)
	}

	return request, err
}

func (c *Component) generateComponentSnapshot(duration int) (result model.ComponentSnapshot, err error) {
	c.duration += duration
	result.Name = c.Name
	result.GUID = c.guid
	result.DurationSec = c.duration
	result.Metrics = make(map[string]interface{})

	for _, metricsGroup := range c.metrics {
		values, cerr := metricsGroup.generateMetricsSnapshots()

		// we are tolerant of request generation errors and should be able to recover
		if cerr != nil {
			err = accumulateErrors(err, cerr)
		}
		for key, value := range values {
			result.Metrics[key] = value
		}
	}

	return result, nil
}

func (c *Component) clearState() {
	c.duration = 0
	for _, mg := range c.metrics {
		mg.clearState()
	}
}

func generateComponentGUID(name, plugin string) string {
	var buf bytes.Buffer
	buf.WriteString(name)
	buf.WriteRune('.')
	buf.WriteString(plugin)
	return buf.String()
}
