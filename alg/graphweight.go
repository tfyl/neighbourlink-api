package alg

func Weight(area map[string]int,orig, dest string) (string,string,int){
	// get the weight of an edge by getting the average weight of the area

	from,exist := area[orig]
	if exist != true{ // checks if the origin exists to stop errors
		from = 0
	}

	to,exist := area[dest]
	if exist != true{ // checks if the destination exists to stop errors
		from = 0
	}

	weight := (from + to)/2
	// average both weights and return them
	return orig,dest,weight
}
