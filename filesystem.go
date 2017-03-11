package tarfs

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Open returns a new tar filesystem object from path to a tar archive.
// The returned object implements the FileSystem interface in https://godoc.org/github.com/kr/fs#FileSystem.
// It can be used by the fs.WalkFS function.
// It is also enable of reading of a specific file
func Open(name string) (*filesystem, error) {
	var (
		f   = &filesystem{}
		err error
	)

	f.f, err = os.Open(name)
	if err != nil {
		return nil, err
	}

	f.createIndex()
	return f, nil
}

type filesystem struct {
	f     *os.File
	z     *gzip.Reader
	r     *tar.Reader
	index *node
}

// Close closes filesystem
func (f *filesystem) Close() error {
	if f.z != nil {
		f.z.Close()
	}
	return f.f.Close()
}

// Open opens a file inside the filesystem
func (f *filesystem) Open(path string) error {
	path = filepath.Clean(path)
	f.reset()
	for {
		h, err := f.r.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if filepath.Clean(h.Name) == path {

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

// Read reads content currently pointed by the last Open call
func (f *filesystem) Read(b []byte) (int, error) {
	return f.r.Read(b)
}

// ReadDir implements the FileSysyem ReadDir method,
// It returns a list of fileinfos in a given path
func (f *filesystem) ReadDir(dirname string) ([]os.FileInfo, error) {
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
func (f *filesystem) Lstat(name string) (os.FileInfo, error) {
	cursor, err := f.findNode(name)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

// Join implements the FileSystem Join method,
func (f *filesystem) Join(elem ...string) string {
	return filepath.Join(elem...)
}

func (f *filesystem) createIndex() error {
	f.index = newFakeDirNode("/")
	f.reset()

	for {
		h, err := f.r.Next()
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

func (f *filesystem) findNode(path string) (*node, error) {
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

func (f *filesystem) reset() {
	var err error

	if f.z != nil {
		f.z.Close()
	}

	f.f.Seek(0, os.SEEK_SET)
	if f.z, err = gzip.NewReader(f.f); err == nil {
		f.r = tar.NewReader(f.z)
	} else {
		// assuming that the archive is not gzipped
		f.z = nil
		f.f.Seek(0, os.SEEK_SET)
		f.r = tar.NewReader(f.f)
	}
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
