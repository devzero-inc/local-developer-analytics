package process

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestPsCollectWithRealOutput(t *testing.T) {
	// Create a no-op logger for testing
	logger := zerolog.Nop()

	ps := NewPs(logger)

	processes, err := ps.Collect()

	assert.NoError(t, err, "Collect method should not return an error")
	assert.NotEmpty(t, processes, "Collect method should return list of processes")
}
