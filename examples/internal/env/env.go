package env

import (
	"github.com/joho/godotenv"
	"golang.org/x/xerrors"
	"os"
)

// LoadEnvFile loads env var from file
func LoadEnvFile(env string) error {
	path, err := os.Getwd()
	if err != nil {
		return xerrors.New("Error finding current directory")
	}

	if env == "" {
		env = "development"
	}

	godotenv.Load(path + "/" + ".env." + env + ".local")
	if env != "test" {
		godotenv.Load(path + "/" + ".env.local")
	}
	godotenv.Load(path + "/" + ".env." + env)
	godotenv.Load(path + "/.env")

	return nil
}
