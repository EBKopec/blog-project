package utils

import (
	"os"
)

func GetEnv(varName string) string {
	return os.Getenv(varName)
}
