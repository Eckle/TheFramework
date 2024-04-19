package resources

import (
	"github.com/Eckle/TheFramework/db"
	"github.com/Eckle/TheFramework/db/queries"
)

type Resource interface {
	GetResource(params *queries.Params) (*[]map[string]interface{}, error)
	CreateResource(resource *map[string]interface{}) error
	UpdateResource(id int, params *queries.Params, resource *map[string]interface{})
	DeleteResource(id int) error
}

type BaseResource struct {
	Table string
}

func (res BaseResource) GetResource(params *queries.Params) (*[]map[string]interface{}, error) {
	db := db.Database
	resourceList := make([]map[string]interface{}, 0)

	query, args := queries.BuildGetQuery(res.Table, params)

	rows, err := db.Queryx(query, args...)
	if err != nil {
		return &resourceList, err
	}

	for rows.Next() {
		var resource map[string]interface{}
		err = rows.MapScan(resource)
		if err != nil {
			return &resourceList, err
		}
		resourceList = append(resourceList, resource)
	}

	return &resourceList, nil
}

func (res BaseResource) CreateResource(resource *map[string]interface{}) error {
	db := db.Database
	query, args, err := queries.BuildInsertQuery(res.Table, resource)
	if err != nil {
		return err
	}

	var lastInsertId int64
	db.QueryRow(query, args...).Scan(&lastInsertId)

	(*resource)["id"] = lastInsertId

	return nil
}

func (res BaseResource) UpdateResource(id int, params *queries.Params, resource *map[string]interface{}) error {
	db := db.Database
	query, args, err := queries.BuildUpdateQuery(res.Table, id, params, resource)
	if err != nil {
		return err
	}
	db.QueryRowx(query, args...).MapScan(*resource)
	return nil
}

func (res BaseResource) DeleteResource(id int) error {
	db := db.Database

	query, err := queries.BuildDeleteQuery(res.Table)
	if err != nil {
		return err
	}

	_, err = db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
