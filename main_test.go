package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_convert(t *testing.T) { //TODO: cleanup diles
	path, _ := filepath.Abs(filepath.Join("testdata", "example.docx"))
	tmpDir, _ := ioutil.TempDir("", "test")
	defer os.Remove(tmpDir)
	file, err := convert(tmpDir, path)
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(file)
	if os.IsNotExist(err) {
		t.Fatal(err)
	}
	if !strings.HasSuffix(info.Name(), "pdf") {
		t.Fatalf("Converted file is not a pdf (%s)", info.Name())
	}
}
