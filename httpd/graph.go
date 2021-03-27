package httpd

import (
	"encoding/json"
	"neighbourlink-api/db"
	"neighbourlink-api/types"
	"net/http"
)

func RetrieveShortestPath(w http.ResponseWriter, r *http.Request, db *db.DB) {
	// get the query parameters : origin and destination
	origin := r.URL.Query().Get("origin")
	destination := r.URL.Query().Get("destination")

	var path types.ShortestPath
	// instantiates the object/struct for the shortest path which stores the path along with the cost
	var err error

	path.Cost, path.Path, err = db.CalculatePath(origin,destination)
	// calculates the shortest path between the origin area and destination
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	// returns the object in json form
	_ = json.NewEncoder(w).Encode(path)
}
