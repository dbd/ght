package main

import (
	"fmt"

	"github.com/dbd/ght/components/pullRequestDetail"
	"github.com/dbd/ght/internal/api"
)

func main() {
	prs := api.GetPullRequests("is:pr author:@me")
	//for _, pr := range prs {
	//	fmt.Printf("%v+\n", pr)
	//}
	pr := prs[0]
	fmt.Printf("%s\n", pr.Title)

	// fmt.Printf("%+v\n", pr.TimelineItems.Nodes[2])
	// fmt.Printf("%+v\n", prs[len(prs)-1].TimelineItems.Nodes[4].Format())
	fmt.Println(pullRequestDetail.RenderMergeOverlay(pr, 160))

}
