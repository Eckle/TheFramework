package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Eckle/TheFramework/db/queries"
)

func TestBuildGetQuery(t *testing.T) {
	params := queries.Params{
		Filter:  "id=1, name=testing",
		Page:    0,
		PerPage: 10,
	}

	query, _ := queries.BuildGetQuery("namespaces", &params)

	if strings.Compare(query, "SELECT * FROM namespaces WHERE id = ? AND name = ?  LIMIT 10 OFFSET 0") == 0 {
		fmt.Printf("query: %v\n", query)
		t.Errorf("query wrong!")
	}
}
