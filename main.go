package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

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

	for _, row := range rows {
		result = append(result, row[0])
	}

	return result
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
				fmt.Println(arr1)
				fmt.Println("length is ", len(arr1))

				arr2 := readCSV(file2)
				fmt.Println(arr2)
				fmt.Println("length is ", len(arr2))

				//csv1, err := os.Open(file1)
				//if err != nil {
				//	log.Fatalln("Couldn't open the csv file", err)
				//}
				//
				//csv2, err := os.Open(file2)
				//if err != nil {
				//	log.Fatalln("Couldn't open the csv file", err)
				//}
				//
				//// convert to arrays
				//arr1 := []int{}
				//arr2 := []int{}
				//reader := csv.NewReader(bufio.NewReader(csv1))
				//for {
				//	// Read each record from csv
				//	record, err := reader.Read()
				//	if err == io.EOF {
				//		break
				//	}
				//	if err != nil {
				//		log.Fatal(err)
				//	}
				//	arr1 = append(arr1, record[0])
				//}




				//fmt.Print(arr1)
				//fmt.Print(arr2)

				// find commonalities in both arrays
				//commonalities := []string{}
				//
				//for i := 0; i < len(arr1); i++ {
				//	for j := 0; j < len(arr2); j++ {
				//		if arr1[i] == arr2[j] {
				//			commonalities = append(commonalities, arr1[i])
				//			break
				//		}
				//	}
				//}

				//fmt.Println(len(arr1))
				//fmt.Println(len(arr2))
				//fmt.Print(commonalities)



				//dat2, err := ioutil.ReadFile(file2)
				//if err != nil {
				//	log.Fatal(err)
				//}




				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}