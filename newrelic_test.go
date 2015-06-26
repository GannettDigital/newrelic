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
