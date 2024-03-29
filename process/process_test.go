package process

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestFactory_Create(t *testing.T) {
	// Create a no-op logger for testing
	logger := zerolog.Nop()

	// Initialize the factory with the test logger
	factory := NewFactory(logger)

	// Define test cases
	tests := []struct {
		name    string
		pType   string
		wantErr bool
	}{
		{"Create PsutilType", PsutilType, false},
		{"Create PsType", PsType, false},
		{"Create Unsupported", "Unsupported", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp, err := factory.Create(tt.pType)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Additionally check that sp is not nil and possibly that it's the expected type
				assert.NotNil(t, sp, "Created instance should not be nil")

				switch tt.pType {
				case PsutilType:
					_, ok := sp.(*Psutil)
					assert.True(t, ok, "Expected PsutilType instance")
				case PsType:
					_, ok := sp.(*Ps)
					assert.True(t, ok, "Expected PsType instance")
				}
			}
		})
	}
}
