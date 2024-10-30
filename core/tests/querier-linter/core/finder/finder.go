/*
Copyright © 2024 Roberto Scudeller <beto.rvs@gmail.com>

Inspired by: https://marianogappa.github.io/software/2019/06/05/lets-build-a-sql-parser-in-go/

*/

package finder

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

var reservedWords = []string{
	"(", ")", ">=", "<=", "!=", ",", "=", ">", "<",
	"SELECT", "INSERT INTO", "VALUES", "UPDATE",
	"DELETE FROM", "WHERE", "FROM", "SET",
	"JOIN", "LEFT", "RIGHT", "INNER",
	".", "RETURNING", "AND", "OR",
	"TRUE", "FALSE", "COUNT", "NOW()",
}

var (
	re = regexp.MustCompile(`(?m)(\(\w+.*\))`)
)

type QueryData struct {
	Type   string
	Tables []Table
}

type Table struct {
	Name    string
	Alias   string
	Columns []string
}

func ParseOne(sql string) (QueryData, error) {
	q := QueryData{}
	countSelects := strings.Count(sql, "SELECT")
	if countSelects > 1 {
		v := re.FindString(sql)
		// fmt.Println("find nested select", v)
		sql = strings.ReplaceAll(sql, v, "$1")
		cleanV := strings.Replace(v, "(", "", 1)
		q, err := ParseOne(cleanV)
		if err != nil {
			return q, err
		}
	}
	n, err := parse(sql)
	if q.Type == "SELECT" {
		q.Tables = append(q.Tables, n.Tables...)
	}
	return n, err
}

func parse(sql string) (QueryData, error) {
	// fmt.Println("sql: ", sql)
	values := strings.Split(sql, " ")
	queryData := QueryData{
		Tables: []Table{},
	}
	var err error
	switch strings.ToUpper(values[0]) {
	case "SELECT":
		queryData.Type = "SELECT"
		err = queryData.parseSelect(values, strings.Contains(sql, "JOIN"))
	case "INSERT":
		queryData.Type = "INSERT"
		err = queryData.parseInsert(values)
	case "UPDATE":
		queryData.Type = "UPDATE"
		err = queryData.parseUpdate(values)
	case "DELETE":
		queryData.Type = "DELETE"
		err = queryData.parseDelete(values)
	default:
		queryData.Type = "UNKNOWN"
		err = fmt.Errorf("unknown query type")
	}
	return queryData, err
}

func (q *QueryData) parseSelect(values []string, withJoin bool) error {
	simpleTables := []string{}
	aliasesColumns := make(map[string][]string)
	// whereClause := []string{}
	fromClause := []string{}
	findFrom := false
	// withJoin := false
	findWhere := false
	findAlias := false

	for _, v := range values {
		upper := strings.ToUpper(v)
		// fmt.Println("v: ", v)
		if v == "," {
			continue
		}
		if upper == "SELECT" {
			continue
		}
		if upper == "FROM" {
			findFrom = true
		}
		if upper == "WHERE" {
			findWhere = true
		}
		if upper == "AS" {
			findAlias = true
		}
		if !findFrom {
			if strings.Contains(v, ".") {
				// in case of a simple table select with alias
				withJoin = true
				// fmt.Println("append 1", v)
				name := columnWithoutSpecialChars(v)
				alias, column := splitAlias(name)
				aliasesColumns[alias] = append(aliasesColumns[alias], column)
				continue
			} else {
				name := columnWithoutSpecialChars(v)
				// fmt.Println("append 2", name)
				if !checkReservedWords(name) {
					simpleTables = append(simpleTables, name)
					continue
				}
			}
		}
		if findFrom && !findWhere {
			if upper != "FROM" {
				// fmt.Println("append 3", v)
				fromClause = append(fromClause, v)
				continue
			}
		}
		if findWhere {
			if upper != "WHERE" {
				cleanV := columnWithoutSpecialChars(v)
				if !checkReservedWords(cleanV) {
					if !strings.Contains(v, "$") {
						// fmt.Println("append 4", v)
						alias, column := splitAlias(v)
						aliasesColumns[alias] = append(aliasesColumns[alias], column)
					}
				}
			}
		}
	}
	// fmt.Println("checking from clause")
	aliasMap := make(map[string]int)
	tableMap := make(map[int]string)
	joinMap := make(map[int]string)
	// onMap := make(map[string]int)
	if withJoin {
		findJoin := false
		indexJoin := -1
		findAS := false
		indexAs := -1
		findON := false
		indexON := -1

		for i, v := range fromClause {
			upper := strings.ToUpper(v)
			// fmt.Println("v: ", v, " i: ", i)
			if upper == "LEFT" || upper == "RIGHT" || upper == "INNER" {
				continue
			}
			if upper == "JOIN" {
				indexJoin = i
				findJoin = true
				if i > 2 && findON {
					// fmt.Println("find another join")
					findON = false
					indexON = -1
				}
				continue
			}
			if upper == "ON" {
				indexON = i
				findON = true
				continue
			}
			if upper == "AS" {
				indexAs = i
				findAS = true
				continue
			}
			if !checkReservedWords(v) {
				if findJoin && indexJoin != -1 {
					// fmt.Println("added joinMap", v)
					joinMap[i] = v
					indexJoin = -1
					continue
				}
				if findAS && indexAs != -1 {
					// fmt.Println("added aliasMap", v)
					aliasMap[v] = i
					indexAs = -1
					continue
				}
				if findON && indexON != -1 {
					// fmt.Println("added onMap", v)
					// onMap[v] = i
					alias, column := splitAlias(v)
					aliasesColumns[alias] = append(aliasesColumns[alias], column)
					continue
				}

				if !strings.Contains(v, ".") {
					// fmt.Println("added tableMap", v)
					tableMap[i] = v
				}
			}
		}
		if len(aliasMap) != 0 {
			for k, v := range aliasMap {
				// fmt.Println("alias: ", k, " index: ", v)
				t := v - 2
				value, ok := tableMap[t]
				if ok {
					// fmt.Println("table: ", t, " index: ", value)
					q.Tables = append(q.Tables, Table{
						Name:    value,
						Alias:   k,
						Columns: aliasesColumns[k],
					})
				}
				value2, ok2 := joinMap[t]
				if ok2 {
					// fmt.Println("join: ", t, " index: ", value2)
					q.Tables = append(q.Tables, Table{
						Name:    value2,
						Alias:   k,
						Columns: aliasesColumns[k],
					})
				}
			}
		}
	} else {
		if len(fromClause) == 1 {
			// fmt.Println("1 from clause: ", fromClause)
			columns := []string{}
			if len(aliasesColumns) != 0 {
				for k := range aliasesColumns {
					columns = append(columns, k)
				}
			}
			if len(simpleTables) != 0 {
				if !findAlias {
					columns = append(columns, simpleTables...)
				} else {
					size := len(simpleTables) - 1
					for k, v := range simpleTables {
						upper := strings.ToUpper(v)
						// fmt.Printf("k: %v v: %v\n", k, v)
						if k >= 1 && k < size {
							if strings.ToUpper(simpleTables[k-1]) != "AS" && upper != "AS" {
								columns = append(columns, v)
							}
						}
						if k == 0 {
							columns = append(columns, v)
						}
					}
				}

			}

			q.Tables = append(q.Tables, Table{
				Name:    fromClause[0],
				Columns: columns,
			})
		}
	}

	// fmt.Println("simple tables: ", simpleTables)
	// fmt.Println("aliases columns: ", aliasesColumns)
	// fmt.Println("from clause: ", fromClause)
	// fmt.Println("table map: ", tableMap)
	// fmt.Println("alias map: ", aliasMap)
	// fmt.Println("join map: ", joinMap)
	if len(q.Tables) == 0 {
		return fmt.Errorf("no tables found")
	}

	return nil
}

func splitAlias(s string) (string, string) {
	splited := strings.Split(s, ".")
	if len(splited) == 2 {
		return splited[0], splited[1]
	}
	return s, ""
}

func (q *QueryData) parseInsert(values []string) error {
	simpleTables := []string{}
	columns := []string{}
	findValues := false
	findValuesString := false
	findReturing := false
	for _, v := range values {
		// fmt.Println("v: ", v)
		upper := strings.ToUpper(v)
		if upper == "INSERT" {
			continue
		}
		if !findValues {
			if upper == "INTO" {
				continue
			}
			if upper == "VALUES" {
				findValuesString = true
			}
			if strings.Contains(v, "(") {
				// table name with special chars: INSERT INTO table_name(column1, column2)
				tableWithSpecial := strings.Split(v, "(")
				if len(tableWithSpecial) == 2 {
					findValues = true
					simpleTables = append(simpleTables, tableWithSpecial[0])
					if strings.Contains(tableWithSpecial[1], ")") {
						columnWithoutComma := columnWithoutSpecialChars(tableWithSpecial[1])
						// fmt.Println("append 1", columnWithoutComma)
						columns = append(columns, columnWithoutComma)
						continue
					} else {
						columnWithoutComma := columnWithoutSpecialChars(tableWithSpecial[1])
						// fmt.Println("append 2", columnWithoutComma)
						columns = append(columns, columnWithoutComma)
						continue
					}
				}
			} else {
				// table name without special chars: INSERT INTO table_name VALUES
				simpleTables = append(simpleTables, strings.TrimSpace(v))
			}
		}
		if v == "(" {
			findValues = true
			continue
		}
		if findValues && !findValuesString {
			if !checkReservedWords(v) {
				if !strings.Contains(v, "$") {
					columnWithoutComma := columnWithoutSpecialChars(v)
					// fmt.Println("append 3", columnWithoutComma)
					columns = append(columns, columnWithoutComma)
					continue
				}
			}
		}
		if findValuesString {
			if upper != "RETURNING" {
				findReturing = true
			}
			if findReturing {
				if upper != "RETURNING" {
					columnWithoutComma := columnWithoutSpecialChars(v)
					// fmt.Println("append 4", columnWithoutComma)
					columns = append(columns, columnWithoutComma)
				}
			}
		}
	}
	// fmt.Println("simple tables: ", simpleTables)
	// fmt.Println("columns: ", columns)
	if len(simpleTables) != 0 {
		q.Tables = append(q.Tables, Table{
			Name:    simpleTables[0],
			Columns: columns,
		})
	}
	return nil
}

func columnWithoutSpecialChars(column string) string {
	column = strings.ReplaceAll(column, "(", "")
	column = strings.ReplaceAll(column, ")", "")
	column = strings.ReplaceAll(column, ",", "")
	column = strings.ReplaceAll(column, "*", "")
	return column
}

func (q *QueryData) parseUpdate(values []string) error {
	simpleTables := []string{}
	columns := []string{}
	findSet := false
	findWhere := false
	for _, v := range values {
		upper := strings.ToUpper(v)
		if upper == "UPDATE" {
			continue
		}
		if upper == "SET" {
			findSet = true
			continue
		}
		if upper == "WHERE" {
			findWhere = true
			continue
		}
		if !findSet {
			simpleTables = append(simpleTables, v)
		}
		if findSet && !findWhere {
			if !checkReservedWords(v) {
				if !strings.Contains(v, "$") {
					columnWithoutComma := columnWithoutSpecialChars(v)
					// fmt.Println("append 1", columnWithoutComma)
					columns = append(columns, columnWithoutComma)
				}
			}
		}
		if findWhere {
			if v != "WHERE" {
				if !checkReservedWords(v) {
					if !strings.Contains(v, "$") {
						columnWithoutComma := columnWithoutSpecialChars(v)
						// fmt.Println("append 2", columnWithoutComma)
						columns = append(columns, columnWithoutComma)
					}
				}
			}
		}
	}
	// fmt.Println("simple tables: ", simpleTables)
	// fmt.Println("columns: ", columns)
	if len(simpleTables) != 0 {
		q.Tables = append(q.Tables, Table{
			Name:    simpleTables[0],
			Columns: columns,
		})
	}
	return nil
}

func (q *QueryData) parseDelete(values []string) error {
	simpleTables := []string{}
	// columns := []string{}
	findWhere := false
	for _, v := range values {
		upper := strings.ToUpper(v)
		if upper == "DELETE" {
			continue
		}
		if upper == "FROM" {
			continue
		}
		if upper == "WHERE" {
			findWhere = true
			continue
		}
		if !findWhere {
			if !checkReservedWords(v) {
				simpleTables = append(simpleTables, v)
			}
		}
	}
	// fmt.Println("simple tables: ", simpleTables)
	// fmt.Println("columns: ", columns)
	if len(simpleTables) != 0 {
		q.Tables = append(q.Tables, Table{
			Name:    simpleTables[0],
			Columns: []string{},
		})
	}
	return nil
}

func checkReservedWords(s string) bool {
	// fmt.Println("reserved 4", s)
	return slices.Contains(reservedWords, strings.ToUpper(s))
}
