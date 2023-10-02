package routes

import (
	"os"

	"github.com/pro-282/ScreenRecordAPI/internal/api/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func BuildRoutesHandler() *gin.Engine {
	r := gin.New()

	if os.Getenv("APP_ENV") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	r.GET("/health", handlers.HealthHandler)

	record := r.Group("/record")

	record.POST("/start", handlers.StartRecording)
	record.POST("/upload/:uuid", handlers.UploadChunk)
	// record.POST("/stop/:uuid", StopRecording)

	return r
}
