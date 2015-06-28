package newrelic

import (
	"bytes"
	"fmt"
	"math"

	"github.com/neocortical/newrelic/model"
)

// Metric defines
type Metric interface {
	Name() string
	Units() string
	Poll() (float64, error)
}

type metricsGroup interface {
	generateMetricsSnapshots() (map[string]interface{}, error)
	clearState()
}

type simpleMetricsGroup struct {
	metric Metric
	state  model.MetricValue
}

func (mg *simpleMetricsGroup) generateMetricsSnapshots() (result map[string]interface{}, err error) {
	result = make(map[string]interface{})

	mg.state, err = pollMetric(mg.metric, mg.state)
	if err == nil {
		key := generateMetricKey(mg.metric)
		if mg.state.Count == 1 {
			result[key] = mg.state.Total
		} else {
			result[key] = mg.state
		}
	}

	return result, err
}

func (mg *simpleMetricsGroup) clearState() {
	mg.state = model.MetricValue{}
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

func pollMetric(metric Metric, state model.MetricValue) (model.MetricValue, error) {
	val, err := metric.Poll()
	if err != nil {
		return state, fmt.Errorf("%s error: %v", metric.Name(), err)
	}

	return updateState(state, val), nil
}

func updateState(state model.MetricValue, val float64) model.MetricValue {
	if state.Count == 0 {
		state.Min = val
		state.Max = val
	} else {
		state.Min = math.Min(val, state.Min)
		state.Max = math.Max(val, state.Max)
	}
	state.Total += val
	state.Count++
	state.SumOfSquares += val * val
	return state
}

func generateMetricKey(m Metric) string {
	var buf bytes.Buffer
	buf.WriteString("Component/")
	buf.WriteString(m.Name())
	buf.WriteRune('[')
	buf.WriteString(m.Units())
	buf.WriteRune(']')
	return buf.String()
}
