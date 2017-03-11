package tarfs

import (
	"archive/tar"
	"compress/gzip"
	"os"
)

// Open is a helper function for opening a tar gz files as file systems
func Open(name string) (tfs *tarFS, err error) {
	f, err := os.Open(name)
	if err != nil {
		return
	}
	defer f.Close()

	var t *tar.Reader

	if z, zipErr := gzip.NewReader(f); zipErr == nil {
		// that archive is zipped
		defer z.Close()
		t = tar.NewReader(z)
	} else {
		// assuming that the archive is not gzipped
		f.Seek(0, os.SEEK_SET)
		t = tar.NewReader(f)
	}

	tfs = New(t)
	return
}
