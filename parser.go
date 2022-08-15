package main

import (
	"errors"
	"fmt"
	"strings"
)

func generateFilter(filters interface{}) string {
	filter := ""
	switch f := filters.(type) {
	case map[string]interface{}:
		filter = "where "
		for key, value := range f {
			switch val := value.(type) {
			case int:
				filter += fmt.Sprintf(`%s=%d`, key, val)
			case float64:
				filter += fmt.Sprintf(`%s=%f`, key, val)
			case string:
				filter += fmt.Sprintf(`%s='%s'`, key, strings.ReplaceAll(val, "'", "''"))
			case map[string]interface{}:

			}
		}
	default:
		return ""
	}
	fmt.Println(filter)
	return filter
}

func Parse(data map[string]interface{}) error {
	rd, ok := data["read"]
	if !ok {
		return errors.New("not found")
	}
	var queries []map[string]string
	for table, val := range rd.(map[string]interface{}) {
		switch v := val.(type) {
		case string:
			queries = append(queries, map[string]string{table: fmt.Sprintf(`select "%s" from "%s" limit 1`, v, table)})
		case []interface{}:
			var columns string
			for _, col := range v {
				switch c := col.(type) {
				case string:
					if c == "*" {
						columns += "*, "
					} else {
						columns += fmt.Sprintf(`"%s", `, c)
					}
				default:
					return errors.New("not acceptable data type")
				}
			}
			columns = strings.TrimRight(columns, ", ")
			queries = append(queries, map[string]string{table: fmt.Sprintf(`select %s from "%s" limit 1`, columns, table)})
		case map[string]interface{}:
			var query string
			if cols, ok := v["cols"]; ok {
				var columns string
				for _, col := range cols.([]interface{}) {
					if col == "*" {
						columns += "*, "
					} else {
						columns += fmt.Sprintf(`"%s", `, col.(string))
					}
				}
				columns = strings.TrimRight(columns, ", ")
				query = fmt.Sprintf(`select %s from "%s"`, columns, table)
			}
			fmt.Println(query)
			fl, ok := v["filter"]
			if ok {
				generateFilter(fl)
			}
		}
	}
	fmt.Println(queries)
	return nil
}
