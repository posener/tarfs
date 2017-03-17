package tarfs

import (
	"os"
)

type node struct {
	os.FileInfo
	next map[string]*node
}

func newNode(i os.FileInfo) *node {
	return &node{FileInfo: i, next: map[string]*node{}}
}

func newFakeDirNode(name string) *node {
	return newNode(&fakeFile{name: name, isDir: true})
}

type fakeFile struct {
	name  string
	isDir bool
	os.FileInfo
}

func (f *fakeFile) Name() string { return f.name }
func (f *fakeFile) IsDir() bool  { return f.isDir }
