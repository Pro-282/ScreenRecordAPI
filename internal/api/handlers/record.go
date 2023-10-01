package handlers

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pro-282/ScreenRecordAPI/pkg/response"
)

var (
	recordings = make(map[string]*os.File) // to hold Files that are currently opened
	mu         sync.Mutex
)

func StartRecording(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	// uuid for the video file
	uuid := uuid.NewString()
	fmt.Println(uuid)

	// Create a temporary directory to store the chunks
	tempDir, err := os.MkdirTemp("", uuid)
	fmt.Println("Creating temporary directory for chunks: ", tempDir)
	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			fmt.Sprintf("Failed to create temp dir: %s", err.Error()),
		)
		return
	}

	// Create an empty file to hold the video data
	videoFile, err := os.Create(tempDir + "/video.webm")
	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			fmt.Sprintf("Failed to create video file: %s", err.Error()),
		)
		return
	}

	// Store the file handle in the recordings map
	recordings[uuid] = videoFile

	response.Success(
		c,
		http.StatusOK,
		"Recording started",
		map[string]interface{}{"uuid": uuid})
}

// func UploadChunk(c *gin.Context) {

// }

// func StopRecording(c *gin.Context) {

// }
