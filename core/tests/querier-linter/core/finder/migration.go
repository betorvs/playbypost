package finder

import (
	"bufio"
	"fmt"
	"log/slog"
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

func ParserDBMigration(migrationFile string) (Migration, error) {
	// create a migration struct
	migration := Migration{
		Tables: make(map[string]Table),
	}
	// read the migration file
	// Open the file
	file, err := os.Open(migrationFile)
	if err != nil {
		fmt.Println(err)
		return migration, err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			slog.Error("error closing file", "error", err)
		}
	}()

	// Create a scanner
	scanner := bufio.NewScanner(file)

	// Read and print lines
	table := Table{}
	tableString := ""
	for scanner.Scan() {
		line := scanner.Text()
		// fmt.Println(line)
		if line != "" {
			// check if line is ending migration block
			if line == ");" && table.Name != "" {
				migration.Tables[tableString] = table
			}
			// check if it is a create table
			// match := reTableName.FindString(line)
			m := checkTableName(line)
			if m != "" {
				// if tableString != m && table.Name != "" {
				// 	migration.Tables[tableString] = table
				// }
				// fmt.Println(m, "table found")
				table = Table{Name: m}
				tableString = m
			}

			// check if it is a column
			columnMatch := checkColumnName(line)
			if columnMatch != "" {
				// fmt.Println(columnMatch, "column found")
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

func checkTableName(line string) string {
	match := reTableName.FindString(line)
	return strings.ReplaceAll(match, " (", "")
}

func checkColumnName(line string) string {
	match := reColumnName.FindString(line)
	return strings.TrimSpace(match)
}
