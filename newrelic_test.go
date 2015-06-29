package newrelic

import (
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

func Test_AddPlugin(t *testing.T) {
	client := New("abc123")

	assert.Equal(t, 0, len(client.Plugins))

	client.AddPlugin(&Plugin{Name: "foo", GUID: "com.example.foo"})
	assert.Equal(t, 1, len(client.Plugins))

	client.AddPlugin(&Plugin{Name: "bar", GUID: "com.example.bar"})
	assert.Equal(t, 2, len(client.Plugins))
	assert.Equal(t, "foo", client.Plugins[0].Name)
	assert.Equal(t, "com.example.foo", client.Plugins[0].GUID)
	assert.Equal(t, "bar", client.Plugins[1].Name)
	assert.Equal(t, "com.example.bar", client.Plugins[1].GUID)
}
