package upgo

import "sort"

type Pair struct {
	Key   string
	Value int64
}

func SortMapByValue(m map[string]int64) []Pair {
	p := make([]Pair, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
		i++
	}

	sort.Slice(p, func(i, j int) bool {
		return p[i].Value > p[j].Value
	})

	return p
}
