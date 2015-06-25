package newrelic

// Request is the container that holds a JSON request
type Request struct {
	Agent      Agent               `json:"agent"`
	Components []ComponentSnapshot `json:"components"`
}

// Agent encapsulates the agent info
type Agent struct {
	Host    string `json:"host"`
	Version string `json:"version"`
	Pid     int    `json:"pid"`
}

// ComponentSnapshot encapsulates the current, unset state of a component
type ComponentSnapshot struct {
	Name        string                 `json:"name"`
	GUID        string                 `json:"guid"`
	DurationSec int                    `json:"duration"`
	Metrics     map[string]MetricValue `json:"metrics"`
}

type MetricValue struct {
	Min          float64 `json:"min"`
	Max          float64 `json:"max"`
	Total        float64 `json:"total"`
	Count        int     `json:"count"`
	SumOfSquares float64 `json:"sum_of_squares"`
}
