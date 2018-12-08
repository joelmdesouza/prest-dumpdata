package database

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Table struct {
	Database string `json:"database"`
	Schema   string `json:"schema"`
	Name     string `json:"table"`
}

type Register struct {
	Database string                 `json:"database"`
	Schema   string                 `json:"schema"`
	Table    string                 `json:"table"`
	Fields   map[string]interface{} `json:"fields"`
}

const TablesSelect = `
SELECT
	t.table_catalog as "database",
	t.table_schema as "schema",
	t.table_name as "table"
FROM 
	information_schema.tables t 
WHERE 
	t.table_schema NOT IN('information_schema', 'pg_catalog')
ORDER BY 1, 2, 3
	`

// ParseInsertRegister create insert SQL
func ParseInsertRegister(register Register) (colsName string, colsValue string, values []interface{}, err error) {
	fields := make([]string, 0)
	for key, value := range register.Fields {
		if chkInvalidIdentifier(key) {
			err = errors.New("Insert: Invalid identifier")
			return
		}
		fields = append(fields, fmt.Sprintf(`"%s"`, key))

		switch value.(type) {
		case []interface{}:
			values = append(values, formatArray(value))
		default:
			values = append(values, value)
		}
	}

	colsName = strings.Join(fields, ", ")
	colsValue = createPlaceholders(1, len(values))
	return
}

func createPlaceholders(initial, lenValues int) (ret string) {
	for i := initial; i <= lenValues; i++ {
		if ret != "" {
			ret += ","
		}
		ret += fmt.Sprintf("$%d", i)
	}
	ret = fmt.Sprintf("(%s)", ret)
	return
}

// chkInvalidIdentifier return true if identifier is invalid
func chkInvalidIdentifier(identifer ...string) bool {
	for _, ival := range identifer {
		if ival == "" || len(ival) > 63 || unicode.IsDigit([]rune(ival)[0]) {
			return true
		}
		count := 0
		for _, v := range ival {
			if !unicode.IsLetter(v) &&
				!unicode.IsDigit(v) &&
				v != '(' &&
				v != ')' &&
				v != '_' &&
				v != '.' &&
				v != '-' &&
				v != '*' &&
				v != '[' &&
				v != ']' &&
				v != '"' {
				return true
			}
			if unicode.Is(unicode.Quotation_Mark, v) {
				count++
			}
		}
		if count%2 != 0 {
			return true
		}
	}
	return false
}

// formatArray format slice to a postgres array format
// today support a slice of string, int and fmt.Stringer
func formatArray(value interface{}) string {
	var aux string
	var check = func(aux string, value interface{}) (ret string) {
		if aux != "" {
			aux += ","
		}
		ret = aux + formatArray(value)
		return
	}
	switch value.(type) {
	case []fmt.Stringer:
		for _, v := range value.([]fmt.Stringer) {
			aux = check(aux, v)
		}
		return "{" + aux + "}"
	case []interface{}:
		for _, v := range value.([]interface{}) {
			aux = check(aux, v)
		}
		return "{" + aux + "}"
	case []string:
		for _, v := range value.([]string) {
			aux = check(aux, v)
		}
		return "{" + aux + "}"
	case []int:
		for _, v := range value.([]int) {
			aux = check(aux, v)
		}
		return "{" + aux + "}"
	case string:
		aux := value.(string)
		aux = strings.Replace(aux, `\`, `\\`, -1)
		aux = strings.Replace(aux, `"`, `\"`, -1)
		return `"` + aux + `"`
	case int:
		return strconv.Itoa(value.(int))
	case fmt.Stringer:
		v := value.(fmt.Stringer)
		return formatArray(v.String())
	}
	return ""
}
