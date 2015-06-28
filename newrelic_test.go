package newrelic

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewPlugin(t *testing.T) {
	plugin := NewPlugin("foo", "bar")

	assert.Equal(t, "foo", plugin.Name)
	assert.Equal(t, "bar", plugin.License)

	assert.Equal(t, agentVersion, plugin.agent.Version)
	assert.Equal(t, os.Getpid(), plugin.agent.PID)
	host, err := os.Hostname()
	assert.Nil(t, err) // sanity
	assert.Equal(t, host, plugin.agent.Host)
}

func Test_AppendComponent(t *testing.T) {
	plugin := NewPlugin("foo", "bar")

	assert.Equal(t, 0, len(plugin.Components))

	plugin.AppendComponent(&Component{Name: "foo"})
	assert.Equal(t, 1, len(plugin.Components))

	plugin.AppendComponent(&Component{Name: "bar"})
	assert.Equal(t, 2, len(plugin.Components))
	assert.Equal(t, "foo", plugin.Components[0].Name)
	assert.Equal(t, "bar", plugin.Components[1].Name)
}

func Test_NewMetrict(t *testing.T) {
	i := 0.0
	m := NewMetric("foo", "yards/wombat", func() (float64, error) {
		i++
		if i == 1.0 {
			return i, nil
		}
		return i, errors.New("i is not 1")
	})

	assert.Equal(t, "foo", m.Name())
	assert.Equal(t, "yards/wombat", m.Units())

	val, err := m.Poll()
	assert.Nil(t, err)
	assert.Equal(t, 1.0, val)

	val, err = m.Poll()
	assert.NotNil(t, err)
}

func Test_AddMetric(t *testing.T) {
	c := &Component{Name: "foo"}

	assert.Equal(t, 0, len(c.metrics))

	m := NewMetric("bar", "ducks/furlong", func() (float64, error) { return 1, nil })

	c.AddMetric(m)
	assert.Equal(t, 1, len(c.metrics))
	assert.Equal(t, "bar", c.metrics[0].(*simpleMetricsGroup).metric.Name())
}
