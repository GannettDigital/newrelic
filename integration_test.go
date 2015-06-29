// +build integration

package newrelic

// Ex: go test -tags=integration -license=abc123

import (
	"flag"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var license = flag.String("license", "", "Valid API license key")

func Test_doSend(t *testing.T) {
	client := New(*license)
	plugin := &Plugin{
		Name: "Test Plugin",
		GUID: "com.example.newrelic.test",
	}
	plugin.AddMetric(NewMetric("Test Metric", "rps", func() (float64, error) { return 1.0, nil }))
	client.AddPlugin(plugin)

	result := client.doSend(time.Now())
	assert.False(t, result)
}
