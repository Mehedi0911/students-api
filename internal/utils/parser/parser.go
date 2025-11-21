package parser

import (
	"net/http"
	"strconv"
)

func ParseIdToInt(r *http.Request) (int64, error) {
	id := r.PathValue("id")
	intId, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		return 0, err
	}
	return intId, nil
}
