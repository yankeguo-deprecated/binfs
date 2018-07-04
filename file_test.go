package binfs_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/novakit/binfs"
)

func TestFile_All(t *testing.T) {
	n := binfs.Find("dir2", "file22.html")
	if n == nil {
		t.Fatal("not found")
	}
	if n.Chunk == nil {
		t.Fatal("empty chunk")
	}
	if n.Chunk.Date != binfs7fcedb72721e0b107e4bca45d247d1fceaad1e23.Date {
		t.Fatal("date mismatch")
	}
}

func TestFilesystem_All(t *testing.T) {
	fs := binfs.FileSystem()
	f, err := fs.Open("/dir2/dir21/file212.js")
	if err != nil {
		t.Fatal(err)
	}
	b1, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	b2, err := ioutil.ReadFile("testdata/dir2/dir21/file212.js")
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(b1, b2) != 0 {
		t.Error("not equal")
	}
}
