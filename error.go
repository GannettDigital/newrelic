package newrelic

import (
	"bytes"
	"fmt"
)

// CompositeError accumulates errors from calling metric.Poll(). These errors
// are logged but otherwise ignored so that functioning metrics may still be
// collected.
type CompositeError []error

// Accumulate errors into a single CompositeError
func (ce CompositeError) Accumulate(err error) CompositeError {
	if err == nil {
		return ce
	}

	return append(ce, err)
}

// Error implements the error interface.
func (ce CompositeError) Error() string {
	if len(ce) == 0 {
		return ""
	}
	if len(ce) == 1 {
		return ce[0].Error()
	}
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "%d errors: [", len(ce))
	for i, err := range ce {
		if i != 0 {
			fmt.Fprint(buf, ", ")
		}
		fmt.Fprintf(buf, `"%v"`, err)
	}
	fmt.Fprint(buf, "]")
	return buf.String()
}
