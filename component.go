package newrelic

import "github.com/neocortical/newrelic/model"

// Component encapsulates a component of a plugin
type Component struct {
	Name string

	guid     string
	duration int
	metrics  []metricsGroup
}

// AddMetric adds a new metric definition to the component
func (c *Component) AddMetric(metric Metric) {
	c.metrics = append(c.metrics, &simpleMetricsGroup{metric: metric})
}

func (c *Component) generateComponentSnapshot(duration int) (result model.ComponentSnapshot, err error) {
	c.duration += duration
	result.Name = c.Name
	result.GUID = c.guid
	result.DurationSec = c.duration
	result.Metrics = make(map[string]interface{})

	for _, metricsGroup := range c.metrics {
		values, cerr := metricsGroup.generateMetricsSnapshots()

		// we are tolerant of request generation errors and should be able to recover
		if cerr != nil {
			err = accumulateErrors(err, cerr)
		}
		for key, value := range values {
			result.Metrics[key] = value
		}
	}

	return result, nil
}

func (c *Component) clearState() {
	c.duration = 0
	for _, mg := range c.metrics {
		mg.clearState()
	}
}
