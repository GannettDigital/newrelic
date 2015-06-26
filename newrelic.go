package newrelic

import (
	"os"
	"time"

	"github.com/neocortical/newrelic/model"
)

const (
	// DefaultPollInterval is the recommended poll interval for NewRelic plugins
	DefaultPollInterval = 60
)

const (
	agentGUID    = "net.neocortical.newrelic"
	agentVersion = "0.0.1"
)

type Plugin struct {
	Name         string
	Company      string
	License      string
	PollInterval int
	Verbose      bool
	Components   []*Component

	agent        model.Agent
	lastPollTime time.Time
}

func NewPlugin(name, company, license string, verbose bool) *Plugin {
	result := &Plugin{
		Name:         name,
		Company:      company,
		License:      license,
		PollInterval: DefaultPollInterval,
		Verbose:      verbose,
	}

	result.agent.Version = agentVersion
	result.agent.PID = os.Getpid()
	var err error
	if result.agent.Host, err = os.Hostname(); err != nil {
		panic(err)
	}

	return result
}

func (p *Plugin) AppendComponent(c *Component) {
	p.Components = append(p.Components, c)
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

func generateRequest(p *Plugin, t time.Time) (request model.Request, err error) {
	request.Agent = p.agent

	var duration int
	if p.lastPollTime.IsZero() {
		duration = p.PollInterval
	} else {
		duration = int(t.Sub(p.lastPollTime).Seconds())
	}

	for _, component := range p.Components {
		componentRequest, cerr := generateComponentSnapshot(component, duration)

		// we are tolerant of request generation errors and should be able to recover
		if cerr != nil {
			err = accumulateErrors(err, cerr)
		}
		request.Components = append(request.Components, componentRequest)
	}

	return request, err
}

func generateComponentSnapshot(component *Component, duration int) (result model.ComponentSnapshot, err error) {
	result.Name = component.Name
	result.GUID = component.guid
	result.DurationSec = component.duration + duration
	result.Metrics = make(map[string]interface{})

	for _, metricsGroup := range component.metrics {
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
