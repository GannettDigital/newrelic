package newrelic

import (
	"time"

	"github.com/neocortical/newrelic/model"
)

// Component encapsulates a component of a plugin
type Component struct {
	Name string

	guid     string
	duration time.Duration
	metrics  map[string]*statefulMetric
}

// AddMetric adds a new metric definition to the component
func (c *Component) AddMetric(metric Metric) {
	if c.metrics == nil {
		c.metrics = make(map[string]*statefulMetric)
	}
	c.metrics[generateMetricKey(metric)] = &statefulMetric{metric: metric}
}

func (c *Component) generateComponentSnapshot(duration time.Duration) (result model.ComponentSnapshot, err error) {
	c.duration += duration
	result.Name = c.Name
	result.GUID = c.guid
	result.DurationSec = int(c.duration / time.Second)
	result.Metrics = make(map[string]interface{})

	for k, m := range c.metrics {
		value, cerr := m.generateMetricSnapshot()

		// we are tolerant of request generation errors and should be able to recover
		if cerr != nil {
			err = accumulateErrors(err, cerr)
		}
		result.Metrics[k] = value
	}

	return result, nil
}

func (c *Component) clearState() {
	c.duration = 0
	for _, m := range c.metrics {
		m.clearState()
	}
}
