package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

// Config 는 환경변수에서 값을 읽어줌
// .env 파일이 있을 때는 해당 값을 환경변수에서 읽은 거 처럼 해줌
func Config(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Error().Msg("load .env file failed")
	}

	return os.Getenv(key)
}
