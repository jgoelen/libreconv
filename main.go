package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jgoelen/libreconv/converter"
)

// DefaultTimeout for converting a document
const DefaultTimeout = 30 * time.Second

// SetupRouter configures the Gin engine
func SetupRouter() *gin.Engine {
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	//router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.POST("/to/pdf", convertToPdf)
	return router
}

func main() {
	router := SetupRouter()
	router.Run(":8080")
}

func convertToPdf(c *gin.Context) {
	tmp, err := ioutil.TempDir("/tmp/", "libreconv")
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	defer os.Remove(tmp)
	filePath, err := saveFileToDir(c, tmp)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	conv, err := converter.New(filePath)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	ctx, cancel := createContext(c.Request)
	defer cancel() // Cancel ctx as soon as handleSearch returns.
	convertedFile, err := conv.Run(ctx)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			c.AbortWithError(504, err)
			return
		}
		c.AbortWithError(500, err)
		return
	}
	writeFileToResponse(convertedFile, c)
}

func saveFileToDir(c *gin.Context, dir string) (string, error) {
	file, err := c.FormFile("file")
	if err != nil {
		return "", err
	}
	filePath, _ := filepath.Abs(filepath.Join(dir, file.Filename))
	log.Printf("Save file[name=%s, size=%d bytes] to %s", file.Filename, file.Size, filePath)
	if c.SaveUploadedFile(file, filePath) != nil {
		return "", err
	}
	return filePath, nil
}

func createContext(request *http.Request) (context.Context, context.CancelFunc) {
	duration, err := time.ParseDuration(request.FormValue("timeout"))
	if err == nil {
		return context.WithTimeout(request.Context(), duration)
	}
	return context.WithTimeout(request.Context(), DefaultTimeout)
}

func writeFileToResponse(filePath string, c *gin.Context) {
	name := filepath.Base(filePath)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+name)
	c.Header("Content-Type", "application/octet-stream")
	c.File(filePath)
}
