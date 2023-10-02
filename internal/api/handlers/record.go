package handlers

import (
	"context"
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
	"github.com/sashabaranov/go-openai"
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

	audioPath := finalVideoDir + "/audio.mp3"
	// FFmpeg to extract audio
	cmd = exec.Command(
		"ffmpeg",
		"-i", finalVideoPath,
		"-q:a", "0", "-map", "a", audioPath)
	fmt.Println(cmd.String())

	err = cmd.Run()
	if err != nil {
		response.Error(
			c, http.StatusInternalServerError, "audio generation failed")
		return
	}

	//Translate audio to .srt file
	whisper := openai.NewClient("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")
	//todo: make this an env variable
	ctx := context.Background()

	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: audioPath,
		Format:   openai.AudioResponseFormatSRT,
	}
	resp, err := whisper.CreateTranscription(ctx, req)
	if err != nil {
		fmt.Printf("Transcription error: %v\n", err)
		return
	}
	fmt.Println(resp.Text)

	f, err := os.Create(finalVideoDir + "/subtitle.srt")
	if err != nil {
		fmt.Printf("Could not open file: %v\n", err)
		response.Error(
			c, http.StatusInternalServerError,
			fmt.Sprintf("Could not open file: %v\n", err))
		return
	}
	defer f.Close()
	if _, err := f.WriteString(resp.Text); err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		response.Error(
			c, http.StatusInternalServerError,
			fmt.Sprintf("Error writing to file: %v\n", err))
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
