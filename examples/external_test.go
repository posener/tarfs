package examples

import (
	"github.com/kr/fs"
	"github.com/posener/tarfs"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestTwoFiles(t *testing.T) {
	t.Parallel()

	for _, tarFile := range []string{"./root.tar.gz", "./two-files.tar.gz"} {
		f, close, err := tarfs.Open(tarFile)
		if err != nil {
			t.Fatal(err)
		}

		walker := fs.WalkFS("/", f)

		assert.True(t, walker.Step())
		assert.Equal(t, "/", walker.Path())

		assert.True(t, walker.Step())
		assert.Equal(t, "/a", walker.Path())

		assert.True(t, walker.Step())
		assert.Equal(t, "/a/b", walker.Path())

		assert.True(t, walker.Step())
		assert.Equal(t, "/a/b/c", walker.Path())

		files := []string{}
		assert.True(t, walker.Step())
		files = append(files, walker.Path())
		assert.True(t, walker.Step())
		files = append(files, walker.Path())
		sort.Slice(files, func(i, j int) bool { return files[i] < files[j] })
		assert.Equal(t, []string{"/a/b/c/d", "/a/b/c/e"}, files)

		assert.False(t, walker.Step())

		close()
	}
}

func TestOneFile(t *testing.T) {
	t.Parallel()

	f, close, err := tarfs.Open("./one-file.tar.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer close()

	walker := fs.WalkFS("/", f)

	assert.True(t, walker.Step())
	assert.Equal(t, "/", walker.Path())

	assert.True(t, walker.Step())
	assert.Equal(t, "/a", walker.Path())

	assert.True(t, walker.Step())
	assert.Equal(t, "/a/b", walker.Path())

	assert.True(t, walker.Step())
	assert.Equal(t, "/a/b/c", walker.Path())

	assert.True(t, walker.Step())
	assert.Equal(t, "/a/b/c/d", walker.Path())

	assert.False(t, walker.Step())
}
