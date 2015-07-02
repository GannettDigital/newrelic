package newrelic

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/neocortical/newrelic/model"
	"github.com/stretchr/testify/assert"
)

func Test_NewPlugin(t *testing.T) {
	client := New("abc123")

	assert.Equal(t, "abc123", client.License)

	assert.Equal(t, agentVersion, client.agent.Version)
	assert.Equal(t, os.Getpid(), client.agent.PID)
	host, err := os.Hostname()
	assert.Nil(t, err) // sanity
	assert.Equal(t, host, client.agent.Host)
}

func Test_AddPlugin(t *testing.T) {
	client := New("abc123")

	assert.Equal(t, 0, len(client.Plugins))

	client.AddPlugin(&Plugin{Name: "foo", GUID: "com.example.foo"})
	assert.Equal(t, 1, len(client.Plugins))

	client.AddPlugin(&Plugin{Name: "bar", GUID: "com.example.bar"})
	assert.Equal(t, 2, len(client.Plugins))
	assert.Equal(t, "foo", client.Plugins[0].Name)
	assert.Equal(t, "com.example.foo", client.Plugins[0].GUID)
	assert.Equal(t, "bar", client.Plugins[1].Name)
	assert.Equal(t, "com.example.bar", client.Plugins[1].GUID)
}

func Test_generateRequest_agentAndDurationMath(t *testing.T) {
	t1 := time.Now()
	t2 := t1.Add(time.Second * 15)
	pollInterval := time.Minute

	c := &Client{
		License:      "abc123",
		PollInterval: pollInterval,
		Plugins: []*Plugin{
			&Plugin{
				Name: "MyPlugin",
				GUID: "com.example.myplugin",
				metrics: map[string]*statefulMetric{
					"foo": &statefulMetric{
						metric: NewMetric("foo", "bars", func() (float64, error) { return 1.0, nil }),
					},
				},
			},
		},
		agent: model.Agent{
			Host:    "10.0.0.1",
			Version: "0.0.1",
			PID:     123,
		},
	}

	r, err := c.generateRequest(t2)
	assert.Nil(t, err)
	assert.Equal(t, "10.0.0.1", r.Agent.Host)
	assert.Equal(t, "0.0.1", r.Agent.Version)
	assert.Equal(t, 123, r.Agent.PID)
	assert.Equal(t, pollInterval.Seconds(), r.Plugins[0].DurationSec)

	// set up to test custom duration
	c.Plugins[0].clearState()
	c.lastPollTime = t1

	r, err = c.generateRequest(t2)
	assert.Nil(t, err)
	assert.Equal(t, 15, r.Plugins[0].DurationSec)
}

func Test_doSend_lastPollTimeUpdated(t *testing.T) {
	i := 0
	testSvr := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		switch i {
		case 1:
			http.Error(rw, "bad req", http.StatusBadRequest)
		case 2:
			http.Error(rw, "forbidden", http.StatusForbidden)
		case 3:
			http.Error(rw, "not found", http.StatusNotFound)
		case 4:
			http.Error(rw, "not allowed", http.StatusMethodNotAllowed)
		case 5:
			http.Error(rw, "too large", http.StatusRequestEntityTooLarge)
		case 6:
			http.Error(rw, "server error", http.StatusInternalServerError)
		case 7:
			http.Error(rw, "unavailable", http.StatusServiceUnavailable)
		case 8:
			http.Error(rw, "timeout", http.StatusGatewayTimeout)
		case 9:
			http.Error(rw, "I'm a teapot!", http.StatusTeapot)
		case 10:
			// test default
			http.Error(rw, "herp derp", 8000)
		case 11:
			rw.Write([]byte("OK"))
		}
	}))

	t0 := time.Now()

	c := &Client{
		License:      "abc123",
		PollInterval: time.Minute,
		Plugins: []*Plugin{
			&Plugin{
				Name: "MyPlugin",
				GUID: "com.example.myplugin",
				metrics: map[string]*statefulMetric{
					"foo": &statefulMetric{
						metric: NewMetric("foo", "bars", func() (float64, error) { return 1.0, nil }),
					},
				},
			},
		},
		agent: model.Agent{
			Host:    "10.0.0.1",
			Version: "0.0.1",
			PID:     123,
		},
		lastPollTime: t0,
		client:       &http.Client{},
		url:          testSvr.URL,
	}

	t1 := t0
	for i = 1; i <= 11; i++ {
		t1 = t1.Add(time.Second * 10)
		c.doSend(t1)
		assert.Equal(t, t1, c.lastPollTime)
	}
}

func Test_New(t *testing.T) {
	nr := New("abc123")

	assert.Equal(t, "abc123", nr.License)
	assert.Equal(t, DefaultPollInterval, nr.PollInterval)
	assert.Equal(t, apiEndpoint, nr.url)
	assert.Equal(t, agentVersion, nr.agent.Version)
	assert.Equal(t, os.Getpid(), nr.agent.PID)
	host, _ := os.Hostname()
	assert.Equal(t, host, nr.agent.Host)
}
