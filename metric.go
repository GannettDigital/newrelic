package newrelic

import "github.com/neocortical/newrelic/model"

// Metric defines
type Metric interface {
	Name() string
	Units() string
	Poll() (float64, error)
}

type metricsGroup interface {
	generateMetricsSnapshots() map[string]interface{}
}

type simpleMetricsGroup struct {
	metric Metric
	state  model.MetricValue
}

func (mg *simpleMetricsGroup) generateMetricsSnapshots() (result map[string]interface{}, err error) {
	result = make(map[string]interface{})

	return result, nil
}

type metric struct {
	state  model.MetricValue
	metric Metric
}

func NewMetric(name, units string, pollFn func() (float64, error)) Metric {
	return &simpleMetric{
		name:  name,
		units: units,
		poll:  pollFn,
	}
}

type simpleMetric struct {
	name  string
	units string
	poll  func() (float64, error)
}

func (sm *simpleMetric) Name() string           { return sm.name }
func (sm *simpleMetric) Units() string          { return sm.units }
func (sm *simpleMetric) Poll() (float64, error) { return sm.poll() }
