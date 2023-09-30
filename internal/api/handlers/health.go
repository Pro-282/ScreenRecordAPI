package handlers

import (
	"net/http"

	"github.com/pro-282/ScreenRecordAPI/pkg/response"

	"github.com/gin-gonic/gin"
)

func HealthHandler(c *gin.Context) {
	response.Success(c, http.StatusOK, "Screen recorder API", nil)
}
