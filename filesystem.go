package tarfs

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// NewFS returns a new tar FileSystem object from path to a tar archive.
// The returned object implements the FileSystem interface in https://godoc.org/github.com/kr/fs#FileSystem.
// It can be used by the fs.WalkFS function.
// It also enables reading of a specific fakeFile.
func NewFS(path string) (*FileSystem, error) {
	fs := &FileSystem{}

	f, err := NewFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fs.createIndex(f.Reader)

	return fs, nil
}

// FileSystem is a struct that describes a tar filesystem.
// It should be created with the NewFile function.
type FileSystem struct {
	index *node
}

// ReadDir implements the FileSystem ReadDir method,
// It returns a list of fileinfos in a given path
func (f *FileSystem) ReadDir(dirname string) ([]os.FileInfo, error) {
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
func (f *FileSystem) Lstat(name string) (os.FileInfo, error) {
	cursor, err := f.findNode(name)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

// Join implements the FileSystem Join method,
func (f *FileSystem) Join(elem ...string) string {
	return filepath.Join(elem...)
}

func (f *FileSystem) createIndex(t *tar.Reader) error {
	f.index = newFakeDirNode("/")

	for {
		h, err := t.Next()
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

func (f *FileSystem) findNode(path string) (*node, error) {
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

// splitPath splits a FileSystem path to directories and ending fakeFile/directory along it's path.
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
