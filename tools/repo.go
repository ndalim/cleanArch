package tools

import "redditapp/pkg/post"

func CalculateUpvotePercentage(pst *post.Post) {

	size := len(pst.Votes)
	if size == 0 {
		pst.UpvotePercentage = 0
		return
	}

	plus := 0
	for _, v := range pst.Votes {
		if v.Vote != 1 {
			continue
		}
		plus += 1
	}
	pst.UpvotePercentage = 100 / size * plus
	return
}
