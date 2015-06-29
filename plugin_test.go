package newrelic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AddMetric(t *testing.T) {
	c := &Plugin{Name: "foo"}

	assert.Equal(t, 0, len(c.metrics))

	m := NewMetric("bar", "ducks/furlong", func() (float64, error) { return 1, nil })

	c.AddMetric(m)
	assert.Equal(t, 1, len(c.metrics))
	assert.Equal(t, "bar", c.metrics[generateMetricKey(m)].metric.Name())
}
