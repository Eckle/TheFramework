package queries

import (
	"errors"
	"net/http"
	"strconv"
)

type Params struct {
	Filter  string
	Sort    string
	Fields  string
	Page    int
	PerPage int
}

func ExtractQueryParams(r *http.Request) (*Params, error) {
	query_params := r.URL.Query()
	filter := query_params.Get("filter")
	sort := query_params.Get("sort")
	fields := query_params.Get("fields")

	page := query_params.Get("page")
	per_page := query_params.Get("per_page")

	intPage, err := strconv.Atoi(page)
	if err != nil && page != "" {
		return &Params{Filter: filter, Sort: sort, Fields: fields, Page: 0, PerPage: 0}, errors.New("the page query param must be an integer")
	} else if page == "" {
		intPage = 0
	}

	intPerPage, err := strconv.Atoi(per_page)
	if err != nil && per_page != "" {
		return &Params{Filter: filter, Sort: sort, Fields: fields, Page: intPage, PerPage: 0}, errors.New("the per_page query param must be an integer")
	} else if per_page == "" {
		intPerPage = 20
	}

	return &Params{Filter: filter, Sort: sort, Fields: fields, Page: intPage, PerPage: intPerPage}, nil
}
