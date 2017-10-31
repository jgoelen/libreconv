package main

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPostFileForConversion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testRouter := SetupRouter()
	resp, err := convertFile(testRouter, "testdata", "example.docx")
	if err != nil {
		t.Errorf("Create request failed %d", err)
	}
	if resp.Code != 200 {
		t.Errorf("Expected 200, received %d.", resp.Code)
	}
	cdh := resp.Header().Get("Content-Disposition")
	if !strings.HasSuffix(cdh, "example.pdf") {
		t.Errorf("Wrong Content-Disposition header %s.", cdh)
	}
	if resp.Body.Len() < 1 {
		t.Errorf("Body should not be empty")
	}
}

func TestPostFileForConversionWithTimeout(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testRouter := SetupRouter()
	resp, err := convertFileWitTimeout(testRouter, "testdata", "example.docx", "10ms")
	if err != nil {
		t.Errorf("Create request failed %d", err)
	}
	if resp.Code != 504 {
		t.Errorf("Expected 504, received %d.", resp.Code)
	}
}

func convertFile(router *gin.Engine, dir string, file string) (*httptest.ResponseRecorder, error) {
	return convertFileWitTimeout(router, dir, file, "10s")
}

func convertFileWitTimeout(router *gin.Engine, dir string, file string, timeout string) (*httptest.ResponseRecorder, error) {
	inputFile, _ := filepath.Abs(filepath.Join("testdata", "example.docx"))
	req, err := multipartRequest("/to/pdf?timeout="+timeout, inputFile)
	if err != nil {
		return nil, err
	}
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	return resp, nil
}

func multipartRequest(uri string, filePath string) (*http.Request, error) {
	log.Printf("Create multipart request uri:%s, file:%s", uri, filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
