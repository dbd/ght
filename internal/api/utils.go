package api

func (c ReviewThreads) GetPRCommentsMap() map[string]map[string]map[LineRange][]Comment {
	var res = map[string]map[string]map[LineRange][]Comment{}
	for _, rt := range c.Nodes {
		lr := LineRange{StartLine: rt.StartLine, EndLine: rt.Line}
		if _, ok := res[rt.Path]; !ok {
			res[rt.Path] = map[string]map[LineRange][]Comment{}
		}
		if _, ok := res[rt.Path][rt.DiffSide]; !ok {
			res[rt.Path][rt.DiffSide] = map[LineRange][]Comment{}
		}
		res[rt.Path][rt.DiffSide][lr] = rt.Comments.Nodes
	}
	return res
}
