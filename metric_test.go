package newrelic

import (
	"errors"
	"testing"

	"github.com/neocortical/newrelic/model"
	"github.com/stretchr/testify/assert"
)

func Test_generateMetricKey(t *testing.T) {
	m := &simpleMetric{
		name:  "foo",
		units: "bars/baz",
		poll:  func() (float64, error) { return 1.0, nil },
	}

	key := generateMetricKey(m)
	assert.Equal(t, "Component/foo[bars/baz]", key)
}

func Test_NewMetric(t *testing.T) {
	poll := func() (float64, error) { return 3.14, nil }

	m := NewMetric("foo", "barns/cowboy", poll)
	assert.Equal(t, "foo", m.Name())
	assert.Equal(t, "barns/cowboy", m.Units())
	val, _ := m.(*simpleMetric).poll()
	assert.Equal(t, 3.14, val)
}

func Test_updateState(t *testing.T) {
	st := model.MetricValue{}

	st = updateState(st, 3.0)
	assert.Equal(t, 3.0, st.Min)
	assert.Equal(t, 3.0, st.Max)
	assert.Equal(t, 3.0, st.Total)
	assert.Equal(t, 1, st.Count)
	assert.Equal(t, 9.0, st.SumOfSquares)

	st = updateState(st, 7.0)
	assert.Equal(t, 3.0, st.Min)
	assert.Equal(t, 7.0, st.Max)
	assert.Equal(t, 10.0, st.Total)
	assert.Equal(t, 2, st.Count)
	assert.Equal(t, 58.0, st.SumOfSquares)

	st = updateState(st, 2.0)
	assert.Equal(t, 2.0, st.Min)
	assert.Equal(t, 7.0, st.Max)
	assert.Equal(t, 12.0, st.Total)
	assert.Equal(t, 3, st.Count)
	assert.Equal(t, 62.0, st.SumOfSquares)

	st = updateState(st, 5.0)
	assert.Equal(t, 2.0, st.Min)
	assert.Equal(t, 7.0, st.Max)
	assert.Equal(t, 17.0, st.Total)
	assert.Equal(t, 4, st.Count)
	assert.Equal(t, 87.0, st.SumOfSquares)
}

func Test_simpleMetricsGroup_generateMetricsSnapshots_error(t *testing.T) {
	m := NewMetric("foo", "barns/cowboy", func() (float64, error) { return 1.0, errors.New("duh-hoy") })
	sm := &statefulMetric{metric: m, state: model.MetricValue{}}

	_, err := sm.generateMetricSnapshot()
	assert.NotNil(t, err)
}

func Test_simpleMetricsGroup_generateMetricsSnapshots(t *testing.T) {
	i := 0.0
	m := NewMetric("foo", "barns/cowboy", func() (float64, error) {
		i++
		return i, nil
	})
	sm := &statefulMetric{metric: m, state: model.MetricValue{}}

	// first pass
	result, err := sm.generateMetricSnapshot()
	assert.Nil(t, err)
	assert.NotNil(t, result)
	floatVal, ok := result.(float64)
	assert.True(t, ok)
	assert.Equal(t, 1.0, floatVal)

	// second pass (returns aggregated value)
	result, err = sm.generateMetricSnapshot()
	assert.Nil(t, err)
	assert.NotNil(t, result)
	aggVal, ok := result.(model.MetricValue)
	assert.True(t, ok)
	assert.Equal(t, model.MetricValue{Min: 1.0, Max: 2.0, Total: 3.0, Count: 2, SumOfSquares: 5.0}, aggVal)
}
