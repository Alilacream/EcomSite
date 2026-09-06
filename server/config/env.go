package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Env struct {
	MainDB             string
	RedisDB            string
	S3_Bucket          string
	JWTSecret          string
	MaxOpenConnections string
	MaxIdleConnections string
	MaxIdleTime        string
}

func AllEnvs() (*Env, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Env{
		MainDB:             os.Getenv("DATABASE_URL"),
		RedisDB:            os.Getenv("REDIS_URL"),
		S3_Bucket:          os.Getenv("S3_BUCKET"),
		JWTSecret:          os.Getenv("SECRET_KEY"),
		MaxIdleTime:        os.Getenv("MAX_IDLE_TIME"),
		MaxOpenConnections: os.Getenv("MAX_OPEN_CONNECTION"),
		MaxIdleConnections: os.Getenv("MAX_IDLE_CONNECTION"),
	}, nil
}

// returns 0 on error
func ParseInt(env string) int {
	Numenv, err := strconv.Atoi(env)
	if err != nil {
		return 0
	}
	return Numenv
}
