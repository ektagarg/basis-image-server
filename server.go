package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	host      = "localhost"
	port      = ":8062"
	filesLoc  = "./basis-files"
	basisuCmd = "./basisu"
)

func main() {
	router := gin.Default()
	router.Static("/files", filesLoc)
	router.POST("/upload", convertAndUpload)
	router.Run(port)
}

func convertAndUpload(c *gin.Context) {
	// Upload a file
	file, err := c.FormFile("file")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(file.Filename)

	err = c.SaveUploadedFile(file, "saved/"+file.Filename)
	if err != nil {
		log.Fatal(err)
	}

	// Execute basisu command to convert png to basis
	cmd := exec.Command(basisuCmd, "saved/"+file.Filename, "-output_path", filesLoc)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Get filename without extension
	filename := strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename))

	c.String(http.StatusOK, fmt.Sprintf("http://%s%s/files/%s.basis", host, port, filename))
}
