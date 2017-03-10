package tarfs

import "archive/tar"

type node struct {
	tar.Header
	next map[string]*node
}

func newNode(h *tar.Header) *node {
	return &node{*h, map[string]*node{}}
}

func newFakeDirNode(name string) *node {
	return newNode(&tar.Header{
		Name:     name,
		Typeflag: tar.TypeDir,
	})
}
