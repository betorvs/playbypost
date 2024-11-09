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

func LoadLibraryFiles() map[string]string {
	libraryFiles := make(map[string]string)
	definitionD10HM := GetEnv("D10HM_DEFINITION", "./library/definitions-d10HM.json")
	definitionPFD20 := GetEnv("PFD20_DEFINITION", "./library/definitions-pfd20.json")
	definitionPFD20Ancestries := GetEnv("PFD20_ANCESTRIES_DEFINITION", "./library/definitions-pfd20-ancestries.json")
	definitionPFD20Backgrounds := GetEnv("PFD20_BACKGROUNDS_DEFINITION", "./library/definitions-pfd20-backgrounds.json")
	definitionPFD20Classes := GetEnv("PFD20_CLASSES_DEFINITION", "./library/definitions-pfd20-classes.json")
	libraryFiles["D10HomeMade"] = definitionD10HM
	libraryFiles["Pathfinder"] = definitionPFD20
	libraryFiles["PFD20Ancestries"] = definitionPFD20Ancestries
	libraryFiles["PFD20Backgrounds"] = definitionPFD20Backgrounds
	libraryFiles["PFD20Classes"] = definitionPFD20Classes

	return libraryFiles
}
