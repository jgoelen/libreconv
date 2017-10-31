package converter

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jgoelen/libreconv/lock"
)

// Libreoffice doesn't handle parallel conversions very well
var sofficeLock = lock.New()

// Converter converts a document to another format
type Converter interface {
	Run(ctx context.Context) (string, error)
}

type converter struct {
	outputDir string
	inputFile string
}

// New creates a new Converter instance
func New(inputFile string) (Converter, error) {
	if _, err := os.Stat(inputFile); err != nil {
		return nil, err
	}
	outputDir := filepath.Dir(inputFile)
	return &converter{outputDir: outputDir, inputFile: inputFile}, nil
}

func (conv *converter) Run(ctx context.Context) (string, error) {
	log.Printf("Convert: %s, out=%s", conv.inputFile, conv.outputDir)
	if !sofficeLock.Try(ctx) {
		return "", ctx.Err()
	}
	defer sofficeLock.Unlock()
	output, err := exec.CommandContext(ctx, "soffice", "--headless", "--convert-to", "pdf", "--outdir", conv.outputDir, conv.inputFile).Output()
	if err != nil {
		return "", err
	}
	log.Printf("Output: %s", string(output))
	pdf, err := findPdf(conv.outputDir)
	if err != nil {
		return "", err
	}
	return filepath.Join(conv.outputDir, pdf), nil
}

func findPdf(dir string) (string, error) {
	files, _ := ioutil.ReadDir(dir)
	log.Printf("Files: %s", files)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), "pdf") {
			return file.Name(), nil
		}
	}
	return "", errors.New("No pdf file in directory")
}
