package logging

import (
	"io"
	"time"

	"github.com/rs/zerolog"
)

var Log = zerolog.Logger{}

// Setup Pass writer. Pass in ioutil.Discard to silence logs.
func Setup(logWriter io.Writer, debug bool) {

	output := zerolog.ConsoleWriter{Out: logWriter, TimeFormat: time.RFC822}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	Log = zerolog.New(output).With().Timestamp().Logger()
}
