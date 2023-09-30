package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pro-282/ScreenRecordAPI/configs"
	"github.com/pro-282/ScreenRecordAPI/internal/routes"
	"github.com/pro-282/ScreenRecordAPI/internal/server"
)

func main() {
	configs.Load()

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(fmt.Sprintf("Failed to conver PORT to integer: %v", err))
	}

	srv := server.NewServer(uint16(port), routes.BuildRoutesHandler())
	srv.Listen()
}
