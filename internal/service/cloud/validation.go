package cloud

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func ValidateDescAsc(r *http.Request) string {
	sort := strings.ToUpper(mux.Vars(r)["sort"])
	if sort != "ASC" && sort != "DESC" {
		sort = "DESC"
	}

	return sort
}
