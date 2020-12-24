package godotenv_test

import (
    "testing"
    "log"

    "github.com/joho/godotenv"
)

func TestReadConfig(t *testing.T) {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
}