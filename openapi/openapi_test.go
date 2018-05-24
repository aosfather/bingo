package openapi

import "testing"

func TestQueryFromYoudaoAsString(t *testing.T) {
	t.Log(QueryFromYoudaoAsString("sex"))
	t.Log(QueryFromYoudaoAsString("建设"))
}
