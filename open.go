package tarfs

import (
	"archive/tar"
	"compress/gzip"
	"os"
)

// Open is a helper function for opening a tar gz files as file systems
func Open(name string) (tfs *tarFS, closeFunc func() error, err error) {
	f, err := os.Open(name)
	if err != nil {
		return
	}

	var t *tar.Reader

	if z, zipErr := gzip.NewReader(f); zipErr == nil {
		// that archive is zipped
		t = tar.NewReader(z)
		closeFunc = func() error {
			z.Close()
			return f.Close()
		}
	} else {
		// assuming that the archive is not gzipped
		f.Seek(0, os.SEEK_SET)
		t = tar.NewReader(f)
		closeFunc = f.Close
	}

	tfs = New(t)
	return
}
