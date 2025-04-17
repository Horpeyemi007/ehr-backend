package config

import "time"

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
	Auth AuthConfig
}

type AuthConfig struct {
	Token TokenConfig
}

type TokenConfig struct {
	Secret string
	Exp    time.Duration
	Iss    string
}
