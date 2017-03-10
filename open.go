package tarfs

import (
	"archive/tar"
	"compress/gzip"
	"os"
)

// Open is a helper function for opening a tar gz files as file systems
func Open(name string) (*tarFS, CloseFunc, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, nil, err
	}
	z, err := gzip.NewReader(f)
	if err != nil {
		return nil, nil, err
	}
	tgz := tar.NewReader(z)

	closeFunc := func() error {
		z.Close()
		return f.Close()
	}

	return New(tgz), closeFunc, nil
}

// CloseFunc is a function that closes something
type CloseFunc func() error
