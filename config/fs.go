package config

import "github.com/spf13/afero"

// Fs is the global file system instance, can be backed by real FS or MemMapFs for testing.
var Fs afero.Fs

// SetupFS initialize the file system instance. When used in testing, the Fs variable is swapped out to afero.MemMapFs.
func SetupFS() {
	Fs = afero.NewOsFs()
}
