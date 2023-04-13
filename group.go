package main

import "fmt"

type Group map[string]int64

func (g Group) String() string {
	var s string

	for query, freq := range g {
		s += fmt.Sprintf("%s %d\n", query, freq)
	}

	return s
}

func (g Group) BytesEncode() []byte {
	return []byte(g.String())
}

func (g Group) Clear() {
	for key := range g {
		delete(g, key)
	}
}
func (g Group) Add(addFreq Group, max int) {
	for k, v := range addFreq {
		if val, ok := g[k]; ok {
			g[k] = val + v
			delete(addFreq, k)
		}
	}

	if len(g) == max {
		return
	}

	for k, v := range addFreq {
		if len(g) < max {
			g[k] = v
			delete(addFreq, k)
		}
	}
}
