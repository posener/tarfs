package tarfs

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestTarFS_ReadDir(t *testing.T) {
	t.Parallel()

	f, close, err := Open("./examples/root.tar.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer close()

	tests := []struct {
		dir   string
		files []file
		err   error
	}{
		{
			dir:   "/",
			files: []file{{name: "a", isDir: true}},
		},
		{
			dir:   "/a",
			files: []file{{name: "b", isDir: true}},
		},
		{
			dir:   "/a/b",
			files: []file{{name: "c", isDir: true}},
		},
		{
			dir:   "/a/b/c",
			files: []file{{name: "d", isDir: false}, {name: "e", isDir: false}},
		},
		{
			dir: "/a/b/c/d",
			err: os.ErrInvalid,
		},
		{
			dir: "/a/b/c/e",
			err: os.ErrInvalid,
		},
		{
			dir: "/b",
			err: os.ErrNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.dir, func(t *testing.T) {
			files, err := f.ReadDir(tt.dir)
			assert.Equal(t, err, tt.err)
			if tt.err == nil {
				assert.Equal(t, len(files), len(tt.files))
				for i := range files {
					assert.Equal(t, files[i].Name(), tt.files[i].name)
					assert.Equal(t, files[i].IsDir(), tt.files[i].isDir)
				}
			}
		})
	}
}

func TestTarFS_Lstat(t *testing.T) {
	t.Parallel()

	f, close, err := Open("./examples/root.tar.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer close()

	tests := []struct {
		path string
		file file
		err  error
	}{
		{
			path: "/",
			file: file{name: "/", isDir: true},
		},
		{
			path: "/a",
			file: file{name: "a", isDir: true},
		},
		{
			path: "/a/b",
			file: file{name: "b", isDir: true},
		},
		{
			path: "/a/b/c",
			file: file{name: "c", isDir: true},
		},
		{
			path: "/a/b/c/d",
			file: file{name: "d", isDir: false},
		},
		{
			path: "/a/b/c/e",
			file: file{name: "e", isDir: false},
		},
		{
			path: "/b",
			err:  os.ErrNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			file, err := f.Lstat(tt.path)
			assert.Equal(t, err, tt.err)
			if tt.err == nil {
				assert.Equal(t, file.Name(), tt.file.name)
				assert.Equal(t, file.IsDir(), tt.file.isDir)
			}
		})
	}
}

func TestSplitPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		path  string
		parts []string
	}{
		{"/", []string{}},
		{"./", []string{}},
		{"a", []string{"a"}},
		{"/a", []string{"a"}},
		{"./a", []string{"a"}},
		{"/a/", []string{"a"}},
		{"a/", []string{"a"}},
		{"/a/b", []string{"a", "b"}},
		{"a/b", []string{"a", "b"}},
		{"a/b/", []string{"a", "b"}},
		{"/a/b/", []string{"a", "b"}},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			assert.Equal(t, tt.parts, splitPath(tt.path))
		})
	}
}
