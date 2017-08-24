package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/to/pdf", func(c *gin.Context) {
		file, _ := c.FormFile("file")
		tmpDir, _ := ioutil.TempDir("", "libreconv")
		defer os.Remove(tmpDir)
		filePath, _ := filepath.Abs(filepath.Join(tmpDir, file.Filename))
		err := c.SaveUploadedFile(file, filePath)
		if err != nil {
			log.Fatal(err)
		}
		pdf, err := convert(tmpDir, filePath)
		if err != nil {
			log.Fatal(err)
		}
		c.String(http.StatusOK, fmt.Sprintf("Converted '%s' to '%s'", file.Filename, pdf))
	})
	router.Run(":8080")
}

func convert(outDir, file string) (string, error) {
	log.Printf("Convert: %s, out=%s", file, outDir)
	output, err := exec.Command("soffice", "--headless", "--convert-to", "pdf", "--outdir", outDir, file).Output()
	log.Printf("Output: %s", output)
	if err != nil {
		return "", err
	}
	pdf, err := findPdf(outDir)
	if err != nil {
		return "", err
	}
	return filepath.Join(outDir, pdf), nil
}

func findPdf(dir string) (string, error) {
	files, _ := ioutil.ReadDir(dir)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), "pdf") {
			return file.Name(), nil
		}
	}
	return "", errors.New("No pdf file in directory")
}
