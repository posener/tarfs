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
	return newNode(&file{name: name, isDir: true})
}

type file struct {
	name  string
	isDir bool
	os.FileInfo
}

func (f *file) Name() string { return f.name }
func (f *file) IsDir() bool  { return f.isDir }
