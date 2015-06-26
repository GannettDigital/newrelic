package newrelic

import (
	"bytes"
	"fmt"
)

func accumulateErrors(err, newError error) error {
	if newError == nil {
		return err
	}
	compErr, ok := err.(*CompositeError)
	if !ok {
		panic("not a composite error")
	}

	compErr.Append(err)
	return compErr
}

type CompositeError struct {
	errs []error
}

func (ce *CompositeError) Append(err error) {
	if err != nil {
		compErr, ok := err.(*CompositeError)
		if ok {
			for _, e := range compErr.errs {
				ce.errs = append(ce.errs, e)
			}
		} else {
			ce.errs = append(ce.errs, err)
		}
	}
}

func (ce *CompositeError) Error() string {
	if len(ce.errs) == 0 {
		return ""
	}
	if len(ce.errs) == 1 {
		return ce.errs[0].Error()
	}
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "%d errors: [", len(ce.errs))
	for i, err := range ce.errs {
		if i != 0 {
			fmt.Fprint(buf, ", ")
		}
		fmt.Fprintf(buf, `"%v"`, err)
	}
	fmt.Fprint(buf, "]")
	return buf.String()
}
