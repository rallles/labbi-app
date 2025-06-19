package config

import "os"

type Config struct {
	ServerAddress string
	Neo4jUri      string
	Neo4jUser     string
	Neo4jPassword string
}

func LoadConfig() Config {
	return Config{
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
		Neo4jUri:      os.Getenv("NEO4J_URI"),
		Neo4jUser:     os.Getenv("NEO4J_USER"),
		Neo4jPassword: os.Getenv("NEO4J_PASSWORD"),
	}
}
