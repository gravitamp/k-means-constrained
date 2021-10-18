// begin
//    specify the number k of clustering to assign.
//    randomly initialize k centroids.
//    repeat
//       expectation: Assign each point to its closest centroid.
//       maximization: Compute the new centroid (mean) of each cluster.
//    until The centroid position do not change.
// end
// Clustering with Constrained Problem for cluster result to have an equal number of member cluster.
// must learn weighted clustering

package main

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"os"
	"sort"
	"strconv"

	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
)

var d Observations
var count []int

func main() {
	//setup data
	setupData("Traffic4.csv")
	// Partition the data points into 20 clusters
	km, _ := NewWithOptions(0.07, SimplePlotter{})

	clusters, _ := km.Partition(d, 20)

	for _, c := range clusters {
		count = append(count, len(c.Observations))
	}
	for i, c := range clusters {
		fmt.Printf("Centered at x: %.2f y: %.2f\n", c.Center[0], c.Center[1])
		// fmt.Printf("Matching data points: %+v\n", c.Observations)
		fmt.Printf("total %d: %d\n", i, len(c.Observations))
	}

	fmt.Println(sum(count))
	min, max := MinMax(count)
	fmt.Println(min, max)
	iter := 0
	for max-min > 10 {
		var count3 []int
		type sorted struct {
			cl int
			dt float64
		}
		var sortcd []sorted
		for i, _ := range clusters {
			// search nearest neighbour

			distA := clusters[0].Center.Distance(clusters[i].Center)
			sortcd = append(sortcd, sorted{
				i,
				distA,
			})
			// }
		}
		sort.SliceStable(sortcd, func(i, j int) bool {
			return sortcd[i].dt < sortcd[j].dt
		})
		//  Plan the steps of adjustment among clusters;
		// 19 step (0, 1) (1,2), dst

		for i := 1; i < len(clusters); i++ {

			var diffA Observations
			var diffB Observations

			//call borderadjust, get new cluster A & B
			if i < len(clusters) {
				// B, _ = clusters.Neighbour(clusters[i].Observations[i], i)
				diffA, diffB = clusters.borderadjust(sortcd[i-1].cl, sortcd[i].cl)
				// fmt.Println(diffA, diffB)
				// } else {
				// 	diffA, diffB = clusters.borderadjust(sortcd[i].cl, sortcd[0].cl)
			}

			if len(diffA) == 0 && len(diffB) == 0 {
				continue
			} else if len(diffA) != 0 && len(diffB) != 0 && i > 0 {
				clusters[sortcd[i-1].cl].Observations = diffA
				if i > 0 && i < len(clusters) {
					for j := 0; j < len(diffB); j++ {
						clusters[sortcd[i].cl].Observations = append(clusters[sortcd[i].cl].Observations, diffB[j])
					}
				}
			}
			clusters.Recenter()
		}
		//recenter

		iter++
		fmt.Println("iterasi ke-", iter)

		for i := 0; i < len(clusters); i++ {
			count3 = append(count3, len(clusters[i].Observations))
		}
		min, max = MinMax(count3)
		fmt.Println(min, max)
		fmt.Println("jarak", max-min)

		//plot
		if km.plotter != nil {
			err := km.plotter.Plot2(clusters, iter)
			if err != nil {
				return //nil, fmt.Errorf("failed to plot chart: %s", err)
			}
		}
		//max iter or no changes (?)
		if iter == 10 {
			break
		}
	}

	var count2 = 0
	//get balanced cluster
	for i, c := range clusters {
		fmt.Printf("Centered at x: %.2f y: %.2f\n", c.Center[0], c.Center[1])
		// fmt.Printf("Matching data points: %+v\n", c.Observations)
		fmt.Printf("total %d: %d\n", i, len(c.Observations))
		count2 += len(c.Observations)
	}
	fmt.Println(count2)

	//plot
	ctx := sm.NewContext()
	ctx.SetSize(500, 500)
	for _, c := range clusters {
		for _, p := range c.Observations {
			ctx.AddObject(
				sm.NewCircle(
					s2.LatLngFromDegrees(p.Coordinates().Coordinates()[1], p.Coordinates().Coordinates()[0]),
					color.RGBA{0, 0xff, 0, 0xff},
					color.RGBA{0, 0x0f, 0, 0xff},
					80.0,
					10.0,
				),
			)
		}
		ctx.AddObject(
			sm.NewMarker(
				s2.LatLngFromDegrees(c.Center[1], c.Center[0]),
				color.RGBA{0, 0, 0xff, 0xff},
				16.0,
			),
		)
	}

	img, err := ctx.Render()
	if err != nil {
		panic(err)
	}

	if err := gg.SavePNG("my-map.png", img); err != nil {
		panic(err)
	}
}

func setupData(file string) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	csvReader := csv.NewReader(f)
	csvData, _ := csvReader.ReadAll()

	//read without header
	for i := 1; i < len(csvData); i++ {
		val, _ := strconv.Atoi(csvData[i][3])
		for j := 0; j < val; j++ {
			lat, _ := strconv.ParseFloat(csvData[i][1], 64)
			lng, _ := strconv.ParseFloat(csvData[i][2], 64)
			d = append(d, Coordinates{
				lng,
				lat,
			})
		}

	}
}

func sum(arr []int) int {
	var res int
	res = 0
	for i := 0; i < len(arr); i++ {
		res += arr[i]
	}
	return res
}

func MinMax(array []int) (int, int) {
	var max int = array[0]
	var min int = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}

type RGBA struct {
	R, G, B, A uint8
}

var colour = []RGBA{
	{0xf9, 0x26, 0x72, 0xff},
	{0x89, 0xbd, 0xff, 0xff},
	{0x66, 0xd9, 0xef, 0xff},
	{0x67, 0x21, 0x0c, 0xff},
	{0x7a, 0xcd, 0x10, 0xff},
	{0xaf, 0x61, 0x9f, 0xff},
	{0xfd, 0x97, 0x1f, 0xff},
	{0xdc, 0xc0, 0x60, 0xff},
	{0x54, 0x52, 0x50, 0xff},
	{0x4b, 0x75, 0x09, 0xff},
	{0xff, 0x00, 0xff, 0xff},
	{0xff, 0x23, 0x54, 0xff},
	{0x45, 0x2c, 0x1a, 0xff},
	{0x22, 0x4f, 0x42, 0xff},
	{0x7a, 0xab, 0x23, 0xff},
	{0xaf, 0x00, 0x2f, 0xff},
	{0xfd, 0x09, 0x3c, 0xff},
	{0xdc, 0x88, 0xf9, 0xff},
	{0x54, 0x92, 0x6e, 0xff},
	{0xc3, 0x21, 0x01, 0xff},
}
