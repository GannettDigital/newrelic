package newrelic

import (
	"bytes"
	l "log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Log(t *testing.T) {
	var b bytes.Buffer

	Logger = l.New(&b, "", 0)
	LogLevel = LogInfo

	Log(LogError, "my error")
	Log(LogInfo, "my info")

	assert.Equal(t, "my error\nmy info\n", b.String())

	LogLevel = LogError

	Log(LogError, "my error 2")
	Log(LogInfo, "my info 2")

	assert.Equal(t, "my error\nmy info\nmy error 2\n", b.String())
}
