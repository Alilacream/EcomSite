package config

// PSQL CONFIG
type DBConfig struct {
	DSN                string
	RedisURL           string
	Password           string
	MaxOpenConnections int
	MaxIdleConnections int
	MaxIdleTime        string
}
