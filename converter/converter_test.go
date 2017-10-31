package converter

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func test(t *testing.T) {
	inputFile, _ := filepath.Abs(filepath.Join("../testdata", "example.docx"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	converter, err := New(inputFile)
	if err != nil {
		t.Fatal(err)
	}
	file, err := converter.Run(ctx)
	defer os.Remove(file)
	if ctx.Err() == context.DeadlineExceeded {
		t.Fatal(ctx.Err())
	}
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
