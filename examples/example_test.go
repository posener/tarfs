package examples

import (
	"fmt"
	"github.com/kr/fs"
	"github.com/posener/tarfs"
	"io/ioutil"
)

// ExampleWalk shows a basic usage of the tarfs package as FileSystem
// for fs.WalkFS function.
func ExampleWalk() {
	// use tarfs.New to open tar.gz files,
	// it returns an object that implements the FileSystem interface.
	f, err := tarfs.New("./root.tar.gz")
	if err != nil {
		panic(err)
	}
	defer f.Close()

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

func ExampleRead() {

	// use tarfs.New to open tar.gz files,
	f, err := tarfs.New("./root.tar.gz")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// open a file inside the file
	f.Open("a/b/c/d")

	// read the file content
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(buf))

	// Output: hello
}
