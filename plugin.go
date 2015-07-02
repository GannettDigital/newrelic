package newrelic

import (
	"time"

	"github.com/neocortical/newrelic/model"
)

// Plugin encapsulates all data and state for a plug-in (AKA Component)
type Plugin struct {
	Name string
	GUID string

	duration time.Duration
	metrics  map[string]*statefulMetric
}

// AddMetric adds a new metric definition to the plugin/component
func (p *Plugin) AddMetric(metric Metric) {
	if p.metrics == nil {
		p.metrics = make(map[string]*statefulMetric)
	}
	p.metrics[generateMetricKey(metric)] = &statefulMetric{metric: metric}
}

func (p *Plugin) generatePluginSnapshot(duration time.Duration) (result model.PluginSnapshot, err CompositeError) {
	p.duration += duration
	result.Name = p.Name
	result.GUID = p.GUID
	result.DurationSec = int(p.duration / time.Second)
	result.Metrics = make(map[string]interface{})

	for k, m := range p.metrics {
		value, cerr := m.generateMetricSnapshot()

		// we are tolerant of request generation errors. metrics that error out are not sent
		if cerr != nil {
			err = err.Accumulate(cerr)
		}
		result.Metrics[k] = value
	}

	return result, nil
}

func (p *Plugin) clearState() {
	p.duration = 0
	for _, m := range p.metrics {
		m.clearState()
	}
}
