package examples

import (
	"fmt"
	"github.com/kr/fs"
	"github.com/posener/tarfs"
)

func ExampleNew() {
	f, close, err := tarfs.Open("./root.tar.gz")
	if err != nil {
		panic(err)
	}
	defer close()

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
