package tarfs

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile_Read(t *testing.T) {
	t.Parallel()

	f, err := NewFile("./examples/root.tar.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	tests := []struct {
		path    string
		content string
		err     error
	}{
		{
			path:    "a/b/c/d",
			content: "hello\n",
		},
		{
			path:    "./a/b/c/d",
			content: "hello\n",
		},
		{
			path:    "a/b/c/e",
			content: "hello\n",
		},
		{
			path: "a",
			err:  os.ErrInvalid,
		},
		{
			path: "b",
			err:  os.ErrNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			err := f.Open(tt.path)
			assert.Equal(t, tt.err, err)
			if tt.err == nil {
				buf, err := ioutil.ReadAll(f)
				if err != nil {
					t.Fatal(err)
				}
				content := []byte(tt.content)
				assert.Equal(t, content, buf)
			}
		})
	}
}
