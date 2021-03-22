package alg

func Weight(area map[string]int,orig, dest string) (string,string,int){
	from,exist := area[orig]
	if exist != true{
		from = 0
	}

	to,exist := area[dest]
	if exist != true{
		from = 0
	}

	weight := (from + to)/2
	//fmt.Println(orig,dest,weight)
	return orig,dest,weight
}
