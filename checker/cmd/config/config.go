package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
)

// Cfg - глобальный конфиг, содержащий переменные .env
var Cfg *Config

type Config struct {
	Kafka            []string
	DatabaseHost     string
	DatabasePort     int
	DatabaseUser     string
	Database         string
	DatabasePassword string
}

// Init - инициализация конфигурации
func Init() {
	pathToEnv := ".env"

	err := godotenv.Load(pathToEnv)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	kafka := strings.Split(os.Getenv("KAFKA"), ",")
	DatabasePort, _ := strconv.Atoi(os.Getenv("DATABASE_PORT"))

	Cfg = &Config{
		Kafka:            kafka,
		DatabaseHost:     os.Getenv("DATABASE_HOST"),
		DatabasePort:     DatabasePort,
		Database:         os.Getenv("DATABASE"),
		DatabaseUser:     os.Getenv("DATABASE_USER"),
		DatabasePassword: os.Getenv("DATABASE_PASSWORD"),
	}

	fmt.Printf("Конфигурация:")
}
