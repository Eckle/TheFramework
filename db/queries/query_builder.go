package queries

import (
	"fmt"
	"log"
	"strings"

	"github.com/Eckle/TheFramework/db"
)

func BuildGetQuery(table string, params *Params) (string, []interface{}) {
	filterString, args, err := filter(params.Filter)
	if err != nil {
		log.Println(err)
	}

	sortString, err := sort(params.Sort)
	if err != nil {
		log.Println(err)
	}

	fieldsString, err := fields(params.Fields)
	if err != nil {
		log.Println(err)
		fieldsString = "*"
	} else if fieldsString == "" {
		fieldsString = "*"
	}

	query := fmt.Sprintf("SELECT %s FROM %s %s %s LIMIT %d OFFSET %d;", fieldsString, table, filterString, sortString, params.PerPage, params.Page)
	return db.Database.Rebind(strings.Join(strings.Fields(query), " ")), args
}

func BuildInsertQuery(table string, data *map[string]interface{}) (string, []interface{}, error) {
	fields, args, err := reflectType(data)
	if err != nil {
		return "", args, err
	}

	var slots []string
	for i := 0; i < len(fields); i++ {
		slots = append(slots, "?")
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id", table, strings.Join(fields, ", "), strings.Join(slots, ", "))

	return db.Database.Rebind(strings.Join(strings.Fields(query), " ")), args, nil
}

func BuildUpdateQuery(table string, id int, params *Params, data *map[string]interface{}) (string, []interface{}, error) {
	fieldsArray, args, err := reflectType(data)
	if err != nil {
		return "", args, err
	}

	args = append(args, id)

	for i := 0; i < len(fieldsArray); i++ {
		fieldsArray[i] = fmt.Sprintf("%s = ?", fieldsArray[i])
	}

	fieldsReturn, err := fields(params.Fields)
	if err != nil {
		log.Println(err)
		fieldsReturn = "RETURNING *"
	} else if fieldsReturn == "" {
		fieldsReturn = "RETURNING *"
	} else {
		fieldsReturn = fmt.Sprintf("RETURNING %s", fieldsReturn)
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ? %s", table, strings.Join(fieldsArray, ", "), fieldsReturn)

	return db.Database.Rebind(strings.Join(strings.Fields(query), " ")), args, nil
}

func BuildDeleteQuery(table string) (string, error) {
	if isValidName(table) {
		return fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, table), nil
	}

	return "", fmt.Errorf("the table name contains invalid characters")
}
