package newrelic

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewPlugin(t *testing.T) {
	plugin := NewPlugin("foo", "bar", true)

	assert.Equal(t, "foo", plugin.Name)
	assert.Equal(t, "bar", plugin.License)
	assert.True(t, plugin.Verbose)

	assert.Equal(t, agentVersion, plugin.agent.Version)
	assert.Equal(t, os.Getpid(), plugin.agent.PID)
	host, err := os.Hostname()
	assert.Nil(t, err) // sanity
	assert.Equal(t, host, plugin.agent.Host)
}

func Test_AppendComponent(t *testing.T) {
	plugin := NewPlugin("foo", "bar", true)

	assert.Equal(t, 0, len(plugin.Components))

	plugin.AppendComponent(&Component{Name: "foo"})
	assert.Equal(t, 1, len(plugin.Components))

	plugin.AppendComponent(&Component{Name: "bar"})
	assert.Equal(t, 2, len(plugin.Components))
	assert.Equal(t, "foo", plugin.Components[0].Name)
	assert.Equal(t, "bar", plugin.Components[1].Name)
}

func Test_AddMetric(t *testing.T) {
	c := &Component{Name: "foo"}

	assert.Equal(t, 0, len(c.metrics))

	m := &TestMetric{
		name:  "bar",
		units: "ducks/furlong",
		poll:  func() (float64, error) { return 1, nil },
	}

	c.AddMetric(m)
	assert.Equal(t, 1, len(c.metrics))
	assert.Equal(t, "bar", c.metrics[0].(simpleMetricsGroup).metric.Name())
}

// Helper objects -------

type TestMetric struct {
	name  string
	units string
	poll  func() (float64, error)
}

func (tm *TestMetric) Name() string           { return tm.name }
func (tm *TestMetric) Units() string          { return tm.units }
func (tm *TestMetric) Poll() (float64, error) { return tm.poll() }
