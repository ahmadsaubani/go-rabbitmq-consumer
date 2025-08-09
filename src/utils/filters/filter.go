package filters

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func parseFilterParam(param string) (field string, operator string) {
	parts := strings.Split(param, "[")
	if len(parts) != 2 {
		return "", ""
	}

	field = parts[0]
	operator = strings.TrimRight(parts[1], "]")
	return field, operator
}

func BuildFilters(ctx *gin.Context) (string, []interface{}, error) {
	var filters []string
	var args []interface{}

	// Iterasi semua query parameter
	for param, values := range ctx.Request.URL.Query() {
		if !strings.Contains(param, "[") || len(values) == 0 {
			continue
		}

		field, operator := parseFilterParam(param)
		if field == "" {
			return "", nil, fmt.Errorf("invalid filter parameter: %s", param)
		}

		value := values[0]

		switch operator {
		case "like", "ilike":
			valueLower := strings.ToLower(value)
			driver := os.Getenv("DB_DRIVER")
			if driver == "postgres" && operator == "ilike" {
				filters = append(filters, fmt.Sprintf("%s ILIKE ?", field))
				args = append(args, "%"+valueLower+"%")
			} else {
				filters = append(filters, fmt.Sprintf("LOWER(%s) LIKE ?", field))
				args = append(args, "%"+valueLower+"%")
			}
		case "moreThan":
			val, err := strconv.Atoi(value)
			if err != nil {
				return "", nil, fmt.Errorf("invalid value for 'moreThan': %s", value)
			}
			filters = append(filters, fmt.Sprintf("%s > ?", field))
			args = append(args, val)

		case "lessThan":
			val, err := strconv.Atoi(value)
			if err != nil {
				return "", nil, fmt.Errorf("invalid value for 'lessThan': %s", value)
			}
			filters = append(filters, fmt.Sprintf("%s < ?", field))
			args = append(args, val)

		case "equals":
			filters = append(filters, fmt.Sprintf("%s = ?", field))
			args = append(args, value)

		case "notEquals":
			filters = append(filters, fmt.Sprintf("%s != ?", field))
			args = append(args, value)

		case "greaterThanOrEqual":
			filters = append(filters, fmt.Sprintf("%s >= ?", field))
			args = append(args, value)

		case "lessThanOrEqual":
			filters = append(filters, fmt.Sprintf("%s <= ?", field))
			args = append(args, value)

		case "in":
			valList := strings.Split(value, ",")
			filters = append(filters, fmt.Sprintf("%s IN (?)", field))
			args = append(args, valList)

		case "notIn":
			valList := strings.Split(value, ",")
			filters = append(filters, fmt.Sprintf("%s NOT IN (?)", field))
			args = append(args, valList)

		default:
			return "", nil, fmt.Errorf("unsupported operator: %s", operator)
		}
	}

	whereClause := strings.Join(filters, " AND ")
	return whereClause, args, nil
}
