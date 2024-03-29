package process

import (
	"bytes"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestPsutilCollectWithRealOutput(t *testing.T) {
	// setup logger to use for the test
	var buf bytes.Buffer
	logger := zerolog.New(&buf)

	psutil := NewPsutil(logger)

	processes, err := psutil.Collect()

	assert.NoError(t, err, "Collect method should not return an error")
	assert.NotEmpty(t, processes, "Collect method should return list of processes")
}
