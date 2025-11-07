package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
}

func Load() *Config {
	_ = godotenv.Load(".env.local")
	//_, filename, _, ok := runtime.Caller(0)
	//if !ok {
	//	log.Fatal("Не удалось определить путь к main.go")
	//}
	//// поднимаемся на 2 уровня вверх до backend
	//dir := filepath.Join(filepath.Dir(filename), "../..")
	//
	//localEnv := filepath.Join(dir, ".env.local")
	//dockerEnv := filepath.Join(dir, ".env.docker")
	//
	//if _, err := os.Stat(localEnv); err == nil {
	//	_ = godotenv.Load(localEnv)
	//} else if _, err := os.Stat(dockerEnv); err == nil {
	//	_ = godotenv.Load(dockerEnv)
	//} else {
	//	log.Fatal("Не найден ни один файл .env")
	//}

	cfg := &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		ServerPort: os.Getenv("SERVER_PORT"),
	}

	if cfg.DBHost == "" {
		log.Fatal("Не найдены переменные окружения в .env.local")
	}

	return cfg
}
