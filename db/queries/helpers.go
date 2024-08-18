package queries

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func reflectType(data *map[string]interface{}) ([]string, []interface{}, error) {
	var fields []string
	var args []interface{}

	for key, value := range *data {

		if !isValidName(key) {
			return fields, args, fmt.Errorf("one of the provided fields does not exist in the table: %s", key)
		}

		fields = append(fields, key)
		args = append(args, value)
	}

	return fields, args, nil
}

func filter(filter string) (string, []interface{}, error) {
	filterList := strings.Split(filter, ",")

	var filterConditions []string
	var filterArguments []interface{}

	if len(filterList) == 1 && filterList[0] == "" {
		return "", filterArguments, nil
	}

	for i := 0; i < len(filterList); i++ {
		if len(filterList) == 1 && filterList[0] == "" {
			break
		}

		var field string
		var operator string
		var value string

		if strings.Contains(filterList[i], "=") {
			split := strings.Split(filterList[i], "=")

			field = strings.Trim(split[0], " ")
			operator = "="
			value = split[1]
		} else if strings.Contains(filterList[i], "!=") {
			split := strings.Split(filterList[i], "!=")

			field = strings.Trim(split[0], " ")
			operator = "!="
			value = split[1]
		} else if strings.Contains(filterList[i], "<") {
			split := strings.Split(filterList[i], "<")

			field = strings.Trim(split[0], " ")
			operator = "<"
			value = split[1]
		} else if strings.Contains(filterList[i], ">") {
			split := strings.Split(filterList[i], ">")

			field = strings.Trim(split[0], " ")
			operator = ">"
			value = split[1]
		} else if strings.Contains(filterList[i], "<=") {
			split := strings.Split(filterList[i], "<=")

			field = strings.Trim(split[0], " ")
			operator = "<="
			value = split[1]
		} else if strings.Contains(filterList[i], ">=") {
			split := strings.Split(filterList[i], ">=")

			field = strings.Trim(split[0], " ")
			operator = ">="
			value = split[1]
		} else {
			return "", filterArguments, fmt.Errorf("no filter condition (= != < > <= >=) matches: %s", filterList[i])
		}

		if !isValidName(field) {
			return "", filterArguments, fmt.Errorf("one of the provided fields does not exist in the table: %s", field)
		}

		filterConditions = append(filterConditions, fmt.Sprintf("%s %s ?", field, operator))

		var convertedValue interface{}
		convertedValue, err := strconv.Atoi(value)
		if err != nil {
			convertedValue, err = strconv.ParseFloat(value, 32)
			if err != nil {
				convertedValue = value
			}
		}

		filterArguments = append(filterArguments, convertedValue)

	}

	return fmt.Sprintf("WHERE %s", strings.Join(filterConditions, " AND ")), filterArguments, nil
}

func sort(sort string) (string, error) {
	sortList := strings.Split(sort, ",")

	var sortConditions []string

	for i := 0; i < len(sortList); i++ {
		if len(sortList) == 1 && sortList[0] == "" {
			break
		}

		var field string

		firstCharIsAlphaNum := unicode.IsLetter(rune(sortList[i][0])) || unicode.IsNumber(rune(sortList[i][0]))

		if firstCharIsAlphaNum {
			field = sortList[i]
		} else {
			field = string([]rune(sortList[i])[1:])
		}

		if !isValidName(field) {
			return "", fmt.Errorf("one of the provided fields is not valid: %s", field)
		}

		if firstCharIsAlphaNum {
			sortConditions = append(sortConditions, field)
		} else if strings.Trim(sortList[i], " ")[0] == '+' {
			sortConditions = append(sortConditions, fmt.Sprintf("%s ASC", field))
		} else if strings.Trim(sortList[i], " ")[0] == '-' {
			sortConditions = append(sortConditions, fmt.Sprintf("%s DESC", field))
		} else {
			return "", fmt.Errorf("first character is invalid: %s", field)
		}
	}

	sortString := ""
	if sortJoin := strings.Join(sortConditions, ", "); sortJoin != "" {
		sortString = fmt.Sprintf("ORDER BY %s", strings.Join(sortConditions, ", "))
	}

	return sortString, nil
}

func fields(fields string) (string, error) {
	fieldList := strings.Split(fields, ",")

	for i := 0; i < len(fieldList); i++ {
		if len(fieldList) == 1 && fieldList[0] == "" {
			break
		}

		if !isValidName(fieldList[i]) {
			return "", fmt.Errorf("one of the provided fields does not exist in the table: %s", fieldList[i])
		}
	}

	return strings.Join(fieldList, ", "), nil
}

func isValidName(input string) bool {
	// This regular expression matches strings that only contain upper and lower case letters, numbers, and underscores.
	// Its a gross over-simplification but it does the trick.
	validStringPattern := `^[a-zA-Z0-9_]*$`
	match, _ := regexp.MatchString(validStringPattern, input)
	return match
}

func AddFilter(params *Params, column string, value interface{}) {
	if params.Filter == "" {
		params.Filter = fmt.Sprintf("%s=%v", column, value)
	} else {
		params.Filter = fmt.Sprintf("%s,%s=%v", params.Filter, column, value)
	}
}
