package handlers

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func ServeVideo(c *gin.Context) {
	uuid := c.Param("uuid")
	videoPath := filepath.Join("./final", uuid, "video.webm")

	// Serve the video file
	c.File(videoPath)
}

func ServeSubtitle(c *gin.Context) {
	uuid := c.Param("uuid")
	srtPath := filepath.Join("./final", uuid, "subtitle.srt")

	// Serve the SRT file
	c.File(srtPath)
}
