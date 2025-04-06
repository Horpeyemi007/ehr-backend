package config

type DbConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

type ServerConfig struct {
	Addr string
	DB   DbConfig
	Env  string
}
