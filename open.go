package tarfs

import (
	"archive/tar"
	"compress/gzip"
	"os"
)

// Open is a helper function for opening a tar gz files as file systems
func Open(name string) (t *tarFS, closeFunc func() error, err error) {
	f, err := os.Open(name)
	if err != nil {
		return
	}
	z, err := gzip.NewReader(f)
	if err != nil {
		return
	}
	tgz := tar.NewReader(z)

	closeFunc = func() error {
		z.Close()
		return f.Close()
	}
	t = New(tgz)
	return
}
