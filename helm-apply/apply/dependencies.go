package apply

import (
	"fmt"
)

func ResolveDependencies(releases map[string]*Release) ([]string, error) {
	graph := map[string][]string{}
	for name, rel := range releases {
		deps := make([]string, len(rel.Dependencies))
		rel.DepReleases = map[string]*Release{}
		i := 0
		for chartName, releaseName := range rel.Dependencies {
			deps[i] = releaseName
			rel.DepReleases[chartName] = releases[releaseName]
			i++
		}
		graph[name] = deps
	}
	order, cycle := sort(graph)
	if cycle != nil {
		return nil, fmt.Errorf("dependency cycle detected %s", cycle)
	}

	for i := len(order)/2 - 1; i >= 0; i-- {
		opp := len(order) - 1 - i
		order[i], order[opp] = order[opp], order[i]
	}

	return order, nil
}

func sort(graph map[string][]string) (order, cycle []string) {
	inDegree := map[string]int{}
	for u, n := range graph {
		if _, ok := inDegree[u]; !ok {
			inDegree[u] = 0
		}
		for _, m := range n {
			inDegree[m]++
		}
	}
	var L, S []string
	rem := map[string]int{}
	for n, d := range inDegree {
		if d == 0 {
			S = append(S, n)
		} else {
			rem[n] = d
		}
	}
	for len(S) > 0 {
		last := len(S) - 1
		n := S[last]
		S = S[:last]
		L = append(L, n)
		for _, m := range graph[n] {
			if rem[m] > 0 {
				rem[m]--
				if rem[m] == 0 {
					S = append(S, m)
				}
			}
		}
	}
	for c, in := range rem {
		if in > 0 {
			for _, nb := range graph[c] {
				if rem[nb] > 0 {
					cycle = append(cycle, c)
					break
				}
			}
		}
	}
	if len(cycle) > 0 {
		return nil, cycle
	}
	return L, nil
}
