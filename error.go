package newrelic

import (
	"bytes"
	"fmt"
)

type CompositeError []error

func (ce CompositeError) Accumulate(err error) CompositeError {
	if err == nil {
		return ce
	}

	return append(ce, err)
}

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
