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

// type metric interface {
// 	generateMetricSnapshot() (interface{}, error)
// 	clearState()
// }

type statefulMetric struct {
	metric Metric
	state  model.MetricValue
}

func (sm *statefulMetric) generateMetricSnapshot() (result interface{}, err error) {
	sm.state, err = pollMetric(sm.metric, sm.state)
	if err == nil {
		if sm.state.Count == 1 {
			result = sm.state.Total
		} else {
			result = sm.state
		}
	}

	return result, err
}

func (sm *statefulMetric) clearState() {
	sm.state = model.MetricValue{}
}

// NewMetric creates a new metric definition using a closure
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
	buf.WriteString("Plugin/")
	buf.WriteString(m.Name())
	buf.WriteRune('[')
	buf.WriteString(m.Units())
	buf.WriteRune(']')
	return buf.String()
}
