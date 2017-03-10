package tarfs

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type tarFS struct {
	tar.Reader
	index *node
}

// New returns a new tar filesystem object from a tar.Reader object.
// This object implements the FileSystem interface in https://godoc.org/github.com/kr/fs#FileSystem.
// It can be used by the Walker object
func New(reader *tar.Reader) *tarFS {
	t := &tarFS{Reader: *reader}
	t.createIndex()
	return t
}

// ReadDir implements the FileSysyem ReadDir method,
// It returns a list of fileinfos in a given path
func (f *tarFS) ReadDir(dirname string) ([]os.FileInfo, error) {
	cursor, err := f.findNode(dirname)
	if err != nil {
		return nil, err
	}
	if !cursor.FileInfo().IsDir() {
		return nil, os.ErrInvalid
	}

	content := make([]os.FileInfo, len(cursor.next))
	i := 0
	for _, h := range cursor.next {
		content[i] = h.FileInfo()
		i++
	}
	return content, nil
}

// Lstat implements the FileSystem Lstat method,
// it returns fileinfo for a given path
func (f *tarFS) Lstat(name string) (os.FileInfo, error) {
	cursor, err := f.findNode(name)
	if err != nil {
		return nil, err
	}
	return cursor.FileInfo(), nil
}

// Join implements the FileSystem Join method,
func (f *tarFS) Join(elem ...string) string {
	return filepath.Join(elem...)
}

func (f *tarFS) createIndex() error {
	f.index = newFakeDirNode("/")
	for {
		h, err := f.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		parts := splitPath(h.Name)

		cursor := f.index
		for _, part := range parts[:len(parts)-1] {
			_, ok := cursor.next[part]
			if !ok {
				cursor.next[part] = newFakeDirNode(part)
			}
			cursor = cursor.next[part]
		}

		cursor.next[parts[len(parts)-1]] = newNode(h)
	}
	return nil
}

func (f *tarFS) findNode(path string) (*node, error) {
	var (
		parts  = splitPath(path)
		cursor = f.index
		ok     bool
	)

	for _, part := range parts {
		cursor, ok = cursor.next[part]
		if !ok {
			return nil, os.ErrNotExist
		}
	}

	return cursor, nil
}

func splitPath(path string) []string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	ret := make([]string, 0, len(parts))
	for _, part := range parts {
		if part != "" {
			ret = append(ret, part)
		}
	}
	return ret
}
