package main

import (
	"encoding/csv"
	"fmt"
	"github.com/golang/geo/s1"
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

func s2LoopsFromFC(fc *geojson.FeatureCollection) []s2.Point {
	coordinates := []s2.Point{}

	for _, coordinate := range fc.Features[0].Geometry.Polygon[0] {
		//s2.PointFromLatLng(coordinate[0])
		//fmt.Println(coordinate[0])
		latLong := s2.LatLng{s1.Angle(coordinate[0]), s1.Angle(coordinate[1])}
		coordinates = append(coordinates, s2.PointFromLatLng(latLong))
	}

	return coordinates
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

				const S2IDLevel = "15"

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

				s2Loops := s2LoopsFromFC(fc)
				fmt.Println(s2Loops)

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}