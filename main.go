package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli"
)

func readCSV(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	result := []string{}

	// If file contains only one row of comma delimited items
	if len(rows) == 1 {
		for _, row := range rows {
			for _, item := range row {
				result = append(result, item)
			}
		}
	} else {
		for _, row := range rows {
			result = append(result, row[0])
		}
	}

	return result
}

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
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

				// find commonalities in both arrays
				commonalities := []string{}

				// remove commonalities from first array (first file)
				for i := 0; i < len(arr1); i++ {
					for j := 0; j < len(arr2); j++ {
						if arr1[i] == arr2[j] {
							commonalities = append(commonalities, arr1[i])
							arr1 = remove(arr1, i)
							break
						}
					}
				}

				// create one new files with commonalities removed from the first file
				fmt.Println("count of comm ", len(commonalities))
				fmt.Println("count of arr1 after ", len(arr1))

				newFileName := strings.Replace(file2, ".csv", "_cleaned.csv", 1)
				newFile, err := os.Create(newFileName)
				if err != nil {
					log.Fatal(err)
				}
				defer newFile.Close()

				writer := csv.NewWriter(newFile)
				defer writer.Flush()

				err = writer.Write(arr1)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println("Commonalities removed and created new file named: ", newFileName)

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}