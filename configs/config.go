package configs

import (
	"fmt"

	"github.com/joho/godotenv"
)

func Load() {
	// Load env variables
	err := godotenv.Load(".env", ".env.local")
	if err != nil {
		fmt.Printf("error: ccannot find .env.local file in the project root")
	}
}
