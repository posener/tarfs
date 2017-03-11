package tarfs

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type tarFS struct {
	index *node
}

// New returns a new tar filesystem object from a tar.Reader object.
// The returned object implements the FileSystem interface in https://godoc.org/github.com/kr/fs#FileSystem.
// It can be used by the fs.WalkFS function.
func New(r *tar.Reader) *tarFS {
	t := &tarFS{}
	t.createIndex(r)
	return t
}

// ReadDir implements the FileSysyem ReadDir method,
// It returns a list of fileinfos in a given path
func (f *tarFS) ReadDir(dirname string) ([]os.FileInfo, error) {
	cursor, err := f.findNode(dirname)
	if err != nil {
		return nil, err
	}
	if !cursor.IsDir() {
		return nil, os.ErrInvalid
	}

	content := make([]os.FileInfo, len(cursor.next))
	i := 0
	for _, f := range cursor.next {
		content[i] = f
		i++
	}
	sort.Slice(content, func(i, j int) bool { return content[i].Name() < content[j].Name() })
	return content, nil
}

// Lstat implements the FileSystem Lstat method,
// it returns fileinfo for a given path
func (f *tarFS) Lstat(name string) (os.FileInfo, error) {
	cursor, err := f.findNode(name)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

// Join implements the FileSystem Join method,
func (f *tarFS) Join(elem ...string) string {
	return filepath.Join(elem...)
}

func (f *tarFS) createIndex(r *tar.Reader) error {
	f.index = newFakeDirNode("/")
	for {
		h, err := r.Next()
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

		cursor.next[parts[len(parts)-1]] = newNode(h.FileInfo())
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

// splitPath splits a filesystem path to directories and ending file/directory along it's path.
func splitPath(path string) []string {
	parts := strings.Split(strings.Trim(filepath.Clean(path), "/"), "/")
	ret := make([]string, 0, len(parts))
	for _, part := range parts {
		switch part {
		case "", ".":
			// skip empty or current directory
		default:
			ret = append(ret, part)
		}
	}
	return ret
}
