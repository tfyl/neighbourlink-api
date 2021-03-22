package httpd

import (
	"encoding/json"
	"neighbourlink-api/db"
	"neighbourlink-api/types"
	"net/http"
)

func RetrieveShortestPath(w http.ResponseWriter, r *http.Request, db *db.DB) {

	origin := r.URL.Query().Get("origin")
	destination := r.URL.Query().Get("destination")

	var path types.ShortestPath
	var err error

	path.Cost, path.Path, err = db.CalculatePath(origin,destination)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	_ = json.NewEncoder(w).Encode(path)
}
