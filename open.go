package tarfs

import (
	"os"
	"compress/gzip"
	"archive/tar"
)

type CloseFunc func() error

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
