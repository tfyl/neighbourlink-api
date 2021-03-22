package db

import (
	"neighbourlink-api/alg"
	"neighbourlink-api/types/dijkstra"
)

func (db *DB) CalculatePath (origin, destination string) (cost int,path []string,err error){
	posts,err := db.GetPostAll()
	if err != nil{
		return 0, nil, nil
	}

	areas := make(map[string]int)
	areas["Camden"] = 1
	areas["Islington"] = 1
	areas["Hackney"] = 1
	areas["Kensington and Chelsea"] = 1
	areas["Westminster"] = 1
	areas["City of London"] = 1
	areas["Tower Hamlets"] = 1
	areas["Southwark"] = 1


	for _ ,post := range posts{
		areas[post.LocalArea] += post.Urgency
	}

	g:= dijkstra.NewGraph()

	g.AddEdge(alg.Weight(areas,"Camden","Islington"))
	g.AddEdge(alg.Weight(areas,"Camden","Westminster"))
	g.AddEdge(alg.Weight(areas,"Camden","City of London"))

	g.AddEdge(alg.Weight(areas,"Islington","Camden"))
	g.AddEdge(alg.Weight(areas,"Islington","Hackney"))
	g.AddEdge(alg.Weight(areas,"Islington","City of London"))

	g.AddEdge(alg.Weight(areas,"Hackney","Islington"))
	g.AddEdge(alg.Weight(areas,"Hackney","City of London"))
	g.AddEdge(alg.Weight(areas,"Hackney","Tower Hamlets"))

	g.AddEdge(alg.Weight(areas,"Kensington and Chelsea","Westminster"))

	g.AddEdge(alg.Weight(areas,"Westminster","Camden"))
	g.AddEdge(alg.Weight(areas,"Westminster","Kensington and Chelsea"))
	g.AddEdge(alg.Weight(areas,"Westminster","City of London"))

	g.AddEdge(alg.Weight(areas,"City of London","Camden"))
	g.AddEdge(alg.Weight(areas,"City of London","Islington"))
	g.AddEdge(alg.Weight(areas,"City of London","Hackney"))
	g.AddEdge(alg.Weight(areas,"City of London","Westminster"))
	g.AddEdge(alg.Weight(areas,"City of London","Tower Hamlets"))
	g.AddEdge(alg.Weight(areas,"City of London","Southwark"))

	g.AddEdge(alg.Weight(areas,"Tower Hamlets","Hackney"))
	g.AddEdge(alg.Weight(areas,"Tower Hamlets","City of London"))
	g.AddEdge(alg.Weight(areas,"Tower Hamlets","Southwark"))

	g.AddEdge(alg.Weight(areas,"Southwark","City of London"))
	g.AddEdge(alg.Weight(areas,"Southwark","Tower Hamlets"))

	cost,path = g.GetPath(origin,destination)
	return
}
