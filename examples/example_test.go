package examples

import (
	"fmt"
	"github.com/kr/fs"
	"github.com/posener/tarfs"
)

// ExampleTest shows a basic usage of the tarfs package as FileSystem
// for fs.WalkFS function.
func ExampleTest() {
	// tarfs.Open is an helper function to open tar.gz files,
	// it returns an object that implements the FileSystem interface.
	f, err := tarfs.Open("./root.tar.gz")
	if err != nil {
		panic(err)
	}

	// WalkFS accepts an object that implements the FileSystem interface,
	// give it the tarFS object created above to walk over the files in the
	// tar file.
	walker := fs.WalkFS("/", f)
	for walker.Step() {
		fmt.Println(walker.Path()) // prints all the paths in the tar file
	}
	if walker.Err() != nil {
		panic(walker.Err())
	}

	// Output: /
	// /a
	// /a/b
	// /a/b/c
	// /a/b/c/d
	// /a/b/c/e
}
