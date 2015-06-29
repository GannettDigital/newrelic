package model

// Request is the container that holds a JSON request
type Request struct {
	Agent   Agent            `json:"agent"`
	Plugins []PluginSnapshot `json:"components"`
}

// Agent encapsulates the agent info
type Agent struct {
	Host    string `json:"host"`
	Version string `json:"version"`
	PID     int    `json:"pid"`
}

// PluginSnapshot encapsulates the current, unset state of a component
type PluginSnapshot struct {
	Name        string                 `json:"name"`
	GUID        string                 `json:"guid"`
	DurationSec int                    `json:"duration"`
	Metrics     map[string]interface{} `json:"metrics"`
}

type MetricValue struct {
	Min          float64 `json:"min"`
	Max          float64 `json:"max"`
	Total        float64 `json:"total"`
	Count        int     `json:"count"`
	SumOfSquares float64 `json:"sum_of_squares"`
}
