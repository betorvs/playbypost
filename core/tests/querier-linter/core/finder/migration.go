package finder

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	reTableName  = regexp.MustCompile(`(?m)(\w+ \()`)
	reColumnName = regexp.MustCompile(`(?m)^\s\s([[:lower:]]\w+)`)
)

type Migration struct {
	Tables map[string]Table `json:"tables"`
}

func parseMigrationContent(content string) (Migration, error) {
	// create a migration struct
	migration := Migration{
		Tables: make(map[string]Table),
	}

	// Create a scanner
	scanner := bufio.NewScanner(strings.NewReader(content))

	// Read and print lines
	table := Table{}
	tableString := ""
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			// check if line is ending migration block
			if line == ");" && table.Name != "" {
				migration.Tables[tableString] = table
				// fmt.Println("tableString: ", tableString)
			}
			// check if it is a create table
			m := checkTableName(line)
			if m != "" {
				// fmt.Println("m: ", m)
				table = Table{Name: m}
				tableString = m
			}

			// check if it is a column
			columnMatch := checkColumnName(line)
			if columnMatch != "" {
				table.Columns = append(table.Columns, strings.TrimSpace(columnMatch))
			}
		}
	}

	// Check for errors
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return Migration{}, err
	}
	return migration, nil
}

func ParserDBMigration(migrationFile string) (Migration, error) {
	// fmt.Println("migrationFile: ", migrationFile)
	content, err := os.ReadFile(migrationFile)
	if err != nil {
		return Migration{}, err
	}
	// fmt.Println("content: ", string(content))
	return parseMigrationContent(string(content))
}

func ParseAllDBMigrations(migrationsDir string) (Migration, error) {
	allMigrations := Migration{
		Tables: make(map[string]Table),
	}

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return allMigrations, err
	}

	for _, file := range files {
		// fmt.Println("file: ", file.Name())
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".up.sql") {
			filePath := fmt.Sprintf("%s/%s", migrationsDir, file.Name())
			migration, err := ParserDBMigration(filePath)
			if err != nil {
				return allMigrations, err
			}
			for k, v := range migration.Tables {
				// fmt.Println("k: ", k)
				allMigrations.Tables[k] = v
			}
		}
	}
	return allMigrations, nil
}

func checkTableName(line string) string {
	match := reTableName.FindString(line)
	// fmt.Println("match: ", match)
	if match != "UNIQUE" {
		return strings.ReplaceAll(match, " (", "")
	}
	return ""
}

func checkColumnName(line string) string {
	match := reColumnName.FindString(line)
	return strings.TrimSpace(match)
}
