package cloud

import (
	"net/http"
	"strings"
)

func ValidateQueryDescAsc(r *http.Request) string {
  sort := strings.ToUpper(r.URL.Query().Get("sort"))
  if sort != "ASC" && sort != "DESC" {
    sort = "DESC"
  }

  return sort
}
