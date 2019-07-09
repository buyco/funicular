package utils

import (
	"fmt"
	"github.com/joho/godotenv"
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
	}
	return nil
}
