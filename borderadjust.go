package main

import (
	"math/rand"
	"sort"
)

type diffsort struct {
	differ float64
	data   Observation
}

func (c Clusters) borderadjust(A int, B int) (Observations, []diffsort) {
	rand.Seed(20)
	// For each point p in area A
	var diff []diffsort
	var diffB []diffsort
	var diffA []diffsort
	var obsA Observations
	var obsB Observations

	for _, p := range c[A].Observations {
		// r, _ = c.Neighbour(p, A)
		distA := p.Distance(c[A].Center)
		distB := p.Distance(c[B].Center)
		//  Calculate diff(p, B) based on (2);
		diff = append(diff, diffsort{distB - distA, p})
		// fmt.Println(A, B, "then", ci)
	}
	//  Sort all the diff(p, B) ascending;
	sort.SliceStable(diff, func(i, j int) bool {
		return diff[i].differ < diff[j].differ
	})
	//  Move the first m point in area A based on sorted
	// diff(p, B) to area B;
	// n := rand.Intn(103-101) + 101
	chunkSize := (len(d) + 20 - 1) / 20
	if len(c[A].Observations) > chunkSize { // && len(c[A].Observations) > len(c[r].Observations) {
		m := len(c[A].Observations) - chunkSize
		for i := 0; i < m; i++ {
			// move to B
			diffB = append(diffB, diff[i])
		}
		diffA = diff[m:]
		for i := 0; i < len(diffA); i++ {
			obsA = append(obsA, diffA[i].data)
		}
		for i := 0; i < len(diffB); i++ {
			obsB = append(obsB, diffB[i].data)
		}
	}
	return obsA, diffB
}
