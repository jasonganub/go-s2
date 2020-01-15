package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/golang/geo/s2"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli"
	"github.com/paulmach/go.geojson"
)

func readCSV(filePath string) []int {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	result := []int{}

	// If file contains only one row of comma delimited items
	if len(rows) == 1 {
		for _, row := range rows {
			for _, item := range row {
				item_int, err := strconv.Atoi(strings.TrimSpace(item))
				if err != nil {
					log.Fatal(err)
				}
				result = append(result, item_int)
			}
		}
	} else {
		for _, row := range rows {
			row_int, err := strconv.Atoi(row[0])
			if err != nil {
				log.Fatal(err)
			}
			result = append(result, row_int)
		}
	}

	return result
}

func remove(slice []int, s int) []int {
	return append(slice[:s], slice[s+1:]...)
}

func getCoordinatesAsPoints(fc *geojson.FeatureCollection) []s2.Point {
	var points []s2.Point

	for _, coordinate := range fc.Features[0].Geometry.Polygon[0] {
		fmt.Println(coordinate)
		latLong := s2.PointFromLatLng(s2.LatLngFromDegrees(coordinate[1], coordinate[0]))
		points = append(points, latLong)
	}

	fmt.Println("points")
	fmt.Println(points)

	return points
}

func getChildren(cellID s2.CellID, level int) []s2.CellID{
	children := []s2.CellID{}

	if cellID.Level() >= level {
		return []s2.CellID{cellID.Parent(level)}
	}

	i := cellID.ChildBeginAtLevel(level)
	for i != cellID.ChildEndAtLevel(level) {
		children = append(children, i)
	}
	return children
}

func getOuterCovering(loops []s2.Point, level int) []s2.CellID{
	maxCells := 100
	regionCoverer := s2.RegionCoverer{MaxLevel: level, MaxCells: maxCells}
	var coverings []s2.CellUnion

	for _, loop := range loops {
		coverings = append(coverings, regionCoverer.CellUnion(loop))
	}

	var coveringsUpdated []s2.CellUnion

	if level > 0 {
		for _, cells := range coverings {
			for _, cell := range cells {
				coveringsUpdated = append(coveringsUpdated, getChildren(cell, level))
			}
		}
	}

	var s2IDS []s2.CellID
	for _, cells := range coveringsUpdated {
		for _, cell := range cells {
			s2IDS = append(s2IDS, cell)
		}
	}
	return s2IDS
}

func getFeature(s int) *geojson.Feature {
	fmt.Println("\ns is ")
	fmt.Println(s)

	c := s2.CellFromCellID(s2.CellID(s))

	fmt.Println("\ncell is ")
	fmt.Println(c)

	var y []s2.LatLng
	for _, i := range []int{0, 1, 2, 3, 0} {
		vertex := c.Vertex(i)
		y = append(y, s2.LatLngFromPoint(vertex))
	}

	fmt.Println("\ny is ")
	fmt.Println(y)

	x := [][][]float64{}
	for _, ll := range y {
		lat := []float64{ll.Lat.Degrees()}
		lng := []float64{ll.Lng.Degrees()}
		row := [][]float64{lat, lng}
		x = append(x, row)
	}

	fmt.Println("\nx is")
	fmt.Println(x)

	geometry := &geojson.Geometry{Polygon: x}

	fmt.Println("\ngeometry is")
	fmt.Println(geometry)

	fmt.Println("string s and string c level")
	fmt.Println(s, c.Level())
	fmt.Println(string(s))
	fmt.Println(string(c.Level()))

	properties := map[string]interface{}{
		"s2id": s,
		"lvl": c.Level(),
	}

	fmt.Println("\nproperties")
	fmt.Println(properties)

	f := geojson.Feature{Geometry: geometry, Properties: properties}

	fmt.Println("\nfeatures")
	fmt.Println(f)

	return nil
}

func getGeojson(s2IDs []s2.CellID) string {
	s2 := []int{}
	for _, s2ID := range s2IDs {
		s2 = append(s2, int(s2ID))
	}

	features := []*geojson.Feature{}
	for _, s := range s2 {
		features = append(features, getFeature(s))
	}
	fmt.Println("features", features)

	featureCollection := geojson.FeatureCollection{Features: features}
	fcStr := json.Marshaler(featureCollection)

	fmt.Println("fc string", fcStr)

	geojson_url := fmt.Sprintf("http://geojson.io/#data=data:application/json,%q", fcStr)
	return geojson_url
}

func main() {
	app := cli.NewApp()
	app.Name = "go-s2id"
	app.Usage = "make an explosive entrance"
	app.Action = func(c *cli.Context) error {
		fmt.Println("boom! I say!")
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:    "diffplace",
			Aliases: []string{"d"},
			Usage:   "diff two s2 id files and move commonalities from first to the second file",
			Action:  func(c *cli.Context) error {

				file1 := ""
				file2 := ""
				if c.NArg() > 0 {
					file1 = c.Args().Get(0)
					file2 = c.Args().Get(1)
				}

				if file1 == "" || file2 == "" {
					fmt.Println("you are required to pass in two file paths")
					return nil
				}

				// read file
				arr1 := readCSV(file1)
				arr2 := readCSV(file2)
				fmt.Println("count of arr1 ", len(arr1))
				fmt.Println("count of arr2 ", len(arr2))

				// find commonalities in both arrays
				commonalities := []int{}

				// remove commonalities from first array (first file)
				for i := 0; i < len(arr1); i++ {
					for j := 0; j < len(arr2); j++ {
						if arr1[i] == arr2[j] {
							commonalities = append(commonalities, arr1[i])
							arr2 = remove(arr2, j)
							break
						}
					}
				}

				if len(commonalities) == 0 {
					fmt.Println("NO COMMONALITIES!")
					os.Exit(0)
				}

				// create one new files with commonalities removed from the first file
				fmt.Println("count of comm ", len(commonalities))
				fmt.Println("count of arr2 after ", len(arr2))

				newFileName := strings.Replace(file2, ".csv", "_cleaned.csv", 1)
				newFile, err := os.Create(newFileName)
				if err != nil {
					log.Fatal(err)
				}
				defer newFile.Close()

				for i := 0; i < len(arr2); i++ {
					_, err = newFile.WriteString(fmt.Sprintf("%d,\n", arr2[i]))
					if err != nil {
						fmt.Printf("error writing string: %v", err)
					}
				}

				fmt.Println("Commonalities removed and created new file named: ", newFileName)

				return nil
			},
		},
		{
			Name:    "get s2ids",
			Aliases: []string{"s2id"},
			Usage:   "",
			Action:  func(c *cli.Context) error {

				const S2IDLevel = 20

				rawFeatureJSON := []byte(`{
				 "type": "FeatureCollection",
				 "features": [
				   {
					 "type": "Feature",
					 "properties": {},
					 "geometry": {
					   "type": "Polygon",
					   "coordinates": [
						 [
						   [
							  107.60782241821289,
							  -6.893148077890368
							],
							[
							  107.60822474956512,
							  -6.894474160996839
							],
							[
							  107.60981798171997,
							  -6.8941492976071075
							],
							[
							  107.60939955711365,
							  -6.892748654540805
							],
							[
							  107.60782241821289,
							  -6.893148077890368
							]
						 ]
					   ]
					 }
				   }
				 ]
				}`)

				fc, err := geojson.UnmarshalFeatureCollection(rawFeatureJSON)
				if err != nil {
					fmt.Printf("error from unmarshalling: %v",  err)
				}

				coordinates := getCoordinatesAsPoints(fc)

				s2IDs := getOuterCovering(coordinates, S2IDLevel)
				fmt.Println("\ns2 ids")
				fmt.Println(s2IDs)

				fmt.Println("\ns2 ids")
				for _, s2id := range s2IDs {
					fmt.Print(uint64(s2id))
					fmt.Print(" ")
				}

				//
				//fmt.Println("\ngeojsonURL")
				//fmt.Println(getGeojson(s2IDs))
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}