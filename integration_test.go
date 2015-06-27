// +build integration

package newrelic

// Ex: go test -tags=integration -license=abc123

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

var license = flag.String("license", "", "Valid API license key")

func Test_doSend(t *testing.T) {
	plugin := NewPlugin("Integration Test Plugin", "net.neocortical", *license)
	component := &Component{
		Name: "Test Component",
	}
	component.AddMetric(NewMetric("Test Metric", "rps", func() (float64, error) { return 1.0, nil }))
	plugin.AppendComponent(component)

	result := plugin.doSend()
	assert.False(t, result)
}
