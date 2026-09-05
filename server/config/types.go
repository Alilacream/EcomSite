package config

type DBConfig struct {
	DSN                string
	MaxOpenConnections int
	MaxIdleConnections int
	MaxIdleTime        string
}
