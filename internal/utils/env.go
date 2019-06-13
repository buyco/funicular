package utils

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func LoadEnvFile(envFile string, runningEnv string) error {
	allowedRunningEnv := []string{"development"}
	exists, _ := InArray(runningEnv, allowedRunningEnv)
	if exists {
		path, _ := os.Getwd()
		err := godotenv.Load(fmt.Sprintf("%s/%s", path, envFile))
		if err != nil {
			return ErrorPrint("Error loading environment file")
		}
	} else {
		log.Print("Environment file not loaded for the current env")
	}
	return nil
}
