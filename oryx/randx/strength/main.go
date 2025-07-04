// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"sort"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"

	"github.com/ory/x/randx"
)

const iterations = 1000 * 100

type generate func(int, []rune) ([]rune, error)

func main() {
	draw(measureDistribution(iterations, randx.AlphaNum, randx.RuneSequence), "AlphaNum Distribution", "docs/alpha_num.png")
	draw(measureDistribution(iterations, randx.Numeric, randx.RuneSequence), "Num Distribution", "docs/num.png")
	draw(measureResultDistribution(100, 6, randx.Numeric, randx.RuneSequence), "Num Distribution", "docs/result_num.png")
}

func measureResultDistribution(iterations int, length int, characters []rune, fn generate) map[string]int {
	dist := make(map[string]int)
	for index := 1; index <= iterations; index++ {
		// status output to cli
		if index%1000 == 0 {
			fmt.Printf("\r%d / %d", index, iterations)
		}
		raw, err := fn(length, characters)
		if err != nil {
			panic(err)
		}
		dist[string(raw)] = dist[string(raw)] + 1
	}
	return dist
}

func measureDistribution(iterations int, characters []rune, fn generate) map[string]int {
	dist := make(map[string]int)
	for index := 1; index <= iterations; index++ {
		// status output to cli
		if index%1000 == 0 {
			fmt.Printf("\r%d / %d", index, iterations)
		}
		raw, err := fn(100, characters)
		if err != nil {
			panic(err)
		}
		for _, s := range raw {
			c := string(s)
			i := dist[c]
			dist[c] = i + 1
		}
	}
	return dist
}

func draw(distribution map[string]int, title, filename string) {
	keys, values := orderMap(distribution)
	group := plotter.Values{}
	for _, v := range values {
		group = append(group, float64(v))
	}

	p := plot.New()
	p.Title.Text = title
	p.Y.Label.Text = "N"

	bars, err := plotter.NewBarChart(group, vg.Points(4))
	if err != nil {
		panic(err)
	}
	bars.LineStyle.Width = vg.Length(0)
	bars.Color = plotutil.Color(0)

	p.Add(bars)
	p.NominalX(keys...)

	if err := p.Save(300*vg.Millimeter, 150*vg.Millimeter, filename); err != nil {
		panic(err)
	}
}

func orderMap(m map[string]int) (keys []string, values []int) {
	keys = []string{}
	values = []int{}
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		values = append(values, m[key])
	}
	return keys, values
}
