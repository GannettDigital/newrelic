package newrelic

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Accumulate(t *testing.T) {
	var ce CompositeError

	ce = ce.Accumulate(nil)
	assert.Nil(t, ce)

	ce = ce.Accumulate(errors.New("herp"))
	assert.NotNil(t, ce)
	assert.Equal(t, 1, len(ce))
	assert.Equal(t, "herp", ce.Error())

	ce = ce.Accumulate(errors.New("derp"))
	assert.NotNil(t, ce)
	assert.Equal(t, 2, len(ce))
	assert.Equal(t, `2 errors: ["herp", "derp"]`, ce.Error())

	var ce2 CompositeError
	ce2 = ce2.Accumulate(ce)
	assert.Equal(t, 1, len(ce2))
	assert.Equal(t, `2 errors: ["herp", "derp"]`, ce2.Error())

	ce2 = ce2.Accumulate(ce)
	assert.Equal(t, 2, len(ce2))
	assert.Equal(t, `2 errors: ["2 errors: ["herp", "derp"]", "2 errors: ["herp", "derp"]"]`, ce2.Error())

	// instantiation edge case
	ce = CompositeError{}
	assert.Equal(t, "", ce.Error())

}
