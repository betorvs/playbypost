package utils

import (
	"fmt"
	"strconv"
)

func LoadDBEnvVars() string {
	host := GetEnv("PGHOST", "localhost")
	user := GetEnv("PGUSER", "postgres")
	pass := GetEnv("PGPASSWORD", "mypassword")
	dbname := GetEnv("PDATABASE", "playbypost")
	sslMode := GetEnv("SSLMode", "disable")
	dbEnvPort := GetEnv("PGPORT", "5432")
	port := 5432
	i, err := strconv.Atoi(dbEnvPort)
	if err == nil {
		port = i
	}
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		host, port, user, pass, dbname, sslMode)
}
