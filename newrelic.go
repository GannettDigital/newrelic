package newrelic

import (
	"os"

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
	License      string
	PollInterval int
	Verbose      bool
	Components   []*Component

	agent        model.Agent
	lastPollTime int64
}

func NewPlugin(name, license string, verbose bool) *Plugin {
	result := &Plugin{
		Name:         name,
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

// func generateRequest(p *Plugin) (request model.Request, err error) {
// 	request.Agent = p.agent

// 	return request, nil
// }

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
	c.metrics = append(c.metrics, simpleMetricsGroup{metric: metric})
}

// Metric defines
type Metric interface {
	Name() string
	Units() string
	Poll() (float64, error)
}

type metricsGroup interface {
}

type simpleMetricsGroup struct {
	metric Metric
}

type metric struct {
	state model.MetricValue
}
