package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
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

func UploadChunk(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	uuid := c.Param("uuid")
	fmt.Println("uuid gotten: ", uuid)

	// Get the file handle from the recordings map
	videoFile, ok := recordings[uuid]
	if !ok {
		response.Error(
			c, http.StatusBadRequest, "Invalid UUID")
		return
	}

	// Append the uploaded chunk to the video file
	_, err := io.Copy(videoFile, c.Request.Body)
	if err != nil {
		response.Error(
			c, http.StatusInternalServerError, "Failed to write chunk")
		return
	}
}

func StopRecording(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	uuid := c.Param("uuid")
	fmt.Println("uuid gotten: ", uuid)

	videoFile, ok := recordings[uuid]
	if !ok {
		response.Error(
			c, http.StatusBadRequest, "Invalid UUID")
		return
	}

	// Close the video file
	videoFile.Close()
	delete(recordings, uuid)

	// Path to the temporary video file
	tempVideoPath := videoFile.Name()
	fmt.Println("temp video path: ", tempVideoPath)

	// Path to the final video file
	finalVideoDir := "./final/" + uuid
	finalVideoPath := finalVideoDir + "/video.webm"

	// Create the final video directory
	err := os.MkdirAll(finalVideoDir, 0755)
	if err != nil {
		log.Printf("Failed to create directory: %v", err)
		response.Error(
			c, http.StatusInternalServerError, "Failed to create video directory")
		return
	}

	// FFmpeg to reformat the video file
	cmd := exec.Command(
		"ffmpeg",
		"-i", tempVideoPath,
		"-c", "copy", // Copy the video and audio streams without re-encoding
		finalVideoPath,
	)

	err = cmd.Run()
	if err != nil {
		response.Error(
			c, http.StatusInternalServerError, "Failed to process video")
		return
	}

	// Delete the temporary video file and directory
	err = os.RemoveAll(tempVideoPath)
	if err != nil {
		log.Printf("Failed to delete temporary video file: %v", err)
	}

	response.Success(
		c, http.StatusOK, "Recording stopped and video processed successfully", nil)
}
