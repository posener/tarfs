package tarfs

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// File is a struct that represent a Reader of a tar file
// it is created by the NewFile function.
type File struct {
	*tar.Reader
	f *os.File
	z *gzip.Reader
}

// NewFile returns a new File object, given a path to a tar file.
func NewFile(path string) (*File, error) {
	f := &File{}
	var err error

	f.f, err = os.Open(path)
	if err != nil {
		return nil, err
	}

	// reset will set the zip and tar readers
	f.reset()

	return f, nil
}

// Close closes the file
func (f *File) Close() error {
	if f.z != nil {
		f.z.Close()
	}
	return f.f.Close()
}

// Open opens a file inside the tar reader. If the returned error
// is a nil, the next Read call will read the requested file inside
// the tar file.
func (f *File) Open(path string) error {
	path = cleanPath(path)

	// cant open "/" for reading
	if path == "" {
		return os.ErrInvalid
	}

	// reset, we need to iterate the tar index from the beginning
	f.reset()

	for {
		h, err := f.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if cleanPath(h.Name) == path {

			// if request path is a directory, we can't open it for reading
			if h.FileInfo().IsDir() {
				return os.ErrInvalid
			}

			// pointing to the right file, stopping the search
			return nil
		}
	}
	return os.ErrNotExist
}

func (f *File) reset() {
	var err error

	if f.z != nil {
		f.z.Close()
	}
	f.f.Seek(0, io.SeekStart)

	if f.z, err = gzip.NewReader(f.f); err == nil {
		f.Reader = tar.NewReader(f.z)
	} else {
		// assuming that the archive is not gzipped
		f.f.Seek(0, io.SeekStart)
		f.Reader = tar.NewReader(f.f)
	}
}

func cleanPath(path string) string {
	return strings.Trim(filepath.Clean(path), "/")
}
