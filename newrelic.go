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
	Components   []Component

	agent model.Agent
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

type Component struct {
	Name string

	guid     string
	duration int
	metrics  map[string]Metric
}

type Metric interface {
	Poll() float64
}

// type metric struct {
// 	state model.MetricValue
// }
