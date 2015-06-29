package newrelic

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewPlugin(t *testing.T) {
	client := New("abc123")

	assert.Equal(t, "abc123", client.License)

	assert.Equal(t, agentVersion, client.agent.Version)
	assert.Equal(t, os.Getpid(), client.agent.PID)
	host, err := os.Hostname()
	assert.Nil(t, err) // sanity
	assert.Equal(t, host, client.agent.Host)
}

func Test_AppendPlugin(t *testing.T) {
	client := New("abc123")

	assert.Equal(t, 0, len(client.Plugins))

	client.AppendPlugin(&Plugin{Name: "foo", GUID: "com.example.foo"})
	assert.Equal(t, 1, len(client.Plugins))

	client.AppendPlugin(&Plugin{Name: "bar", GUID: "com.example.bar"})
	assert.Equal(t, 2, len(client.Plugins))
	assert.Equal(t, "foo", client.Plugins[0].Name)
	assert.Equal(t, "com.example.foo", client.Plugins[0].GUID)
	assert.Equal(t, "bar", client.Plugins[1].Name)
	assert.Equal(t, "com.example.bar", client.Plugins[1].GUID)
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
	c := &Plugin{Name: "foo"}

	assert.Equal(t, 0, len(c.metrics))

	m := NewMetric("bar", "ducks/furlong", func() (float64, error) { return 1, nil })

	c.AddMetric(m)
	assert.Equal(t, 1, len(c.metrics))
	assert.Equal(t, "bar", c.metrics[generateMetricKey(m)].metric.Name())
}
