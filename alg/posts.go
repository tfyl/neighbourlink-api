package alg

import "neighbourlink-api/types"



func OrderPosts(PS []types.Post)  []types.Post {

	t1 := PS[0].Time // highest value of time - first item
	t2 := PS[len(PS)-1].Time // lowest value of time - last item
	diff := t1.Sub(*t2)

	for _,p := range PS {
		p.Time.Sub(*t1)
	}


}

