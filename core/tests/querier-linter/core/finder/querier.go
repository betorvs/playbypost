package finder

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	singleLineNote string = "// dev:finder+query"
	multilineNote  string = "// dev:finder+multiline+query"
)

var (
	// re     = regexp.MustCompile(`(?m)=\s\W(\w.*)\W\s//`)
	reLine = regexp.MustCompile(`(?m)"(\w.*)"`)
	// reMultilineFind    = regexp.MustCompile(`(?m)=\s\W(\w.*)\W`)
	reMultilineFind = regexp.MustCompile(`(?m)\x60(\w.*)\x60`)
	// reMultilineReplace = regexp.MustCompile(`(?m)\W(\s+)//`)
)

func ReadDir(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return []string{}, err
	}
	result := []string{}
	for _, file := range files {
		if !file.IsDir() {
			result = append(result, file.Name())
		}
	}
	return result, nil
}

func QueriesFromDir(migrationFile string) ([]string, error) {
	result := []string{}
	// read the migration file
	// Open the file
	file, err := os.Open(migrationFile)
	if err != nil {
		fmt.Println(err)
		return result, err
	}
	defer file.Close()

	// Create a scanner
	scanner := bufio.NewScanner(file)
	tempLine := ""
	multiline := false
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			if strings.Contains(line, singleLineNote) {
				// fmt.Printf("line >%s<\n", line)
				quotted := reLine.FindString(line)
				clean := cleanLine(quotted)
				result = append(result, clean)
			}
			if multiline {
				if tempLine == "" {
					tempLine = strings.TrimSpace(line)
				} else {
					// fmt.Printf(">%s<\n", trimAllSpace(line))
					if !strings.Contains(line, multilineNote) {
						tempLine = tempLine + " " + trimAllSpace(line)
					}
				}

			}
			if strings.Contains(line, multilineNote) {
				if tempLine != "" {
					// fmt.Printf("tempfile >%s<\n", tempLine)
					quotted := reMultilineFind.FindString(tempLine)
					// fmt.Printf("tempfile 2 >%s<\n", quotted)
					// quotted2 := reMultilineReplace.ReplaceAllString(quotted, "")
					// fmt.Printf("tempfile 3 >%s<\n", quotted2)
					clean := cleanMultiline(quotted)
					result = append(result, clean)
				}
				multiline = !multiline
				tempLine = ""
			}

		}
	}
	return result, nil
}

func cleanLine(line string) string {
	// line = strings.Replace(line, "= \"", "", 1)
	line = strings.ReplaceAll(line, "\"", "")
	return line
}

func cleanMultiline(line string) string {
	// line = strings.Replace(line, "= `", "", 1)
	line = strings.ReplaceAll(line, "`", "")
	return line
}

func trimAllSpace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
