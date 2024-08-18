package tests

import (
	"testing"

	"github.com/Eckle/TheFramework/db/queries"
)

func TestBuildGetQuery(t *testing.T) {
	params := queries.Params{
		Filter:  "name=testing, id=1",
		Page:    0,
		PerPage: 10,
	}

	_, _ = queries.BuildGetQuery("namespaces", &params)
}
