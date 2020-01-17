package diffplace

import (
    "encoding/csv"
    "fmt"
    "github.com/urfave/cli"
    "log"
    "os"
    "strconv"
    "strings"
)

func readCSV(filePath string) []int {
    file, err := os.Open(filePath)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    rows, err := reader.ReadAll()
    var result []int

    // If file contains only one row of comma delimited items
    if len(rows) == 1 {
        for _, row := range rows {
            for _, item := range row {
                itemInt, err := strconv.Atoi(strings.TrimSpace(item))
                if err != nil {
                    log.Fatal(err)
                }
                result = append(result, itemInt)
            }
        }
    } else {
        for _, row := range rows {
            rowInt, err := strconv.Atoi(row[0])
            if err != nil {
                log.Fatal(err)
            }
            result = append(result, rowInt)
        }
    }

    return result
}

func remove(slice []int, s int) []int {
    return append(slice[:s], slice[s+1:]...)
}

// Run as entry point from CLI
func Run(c *cli.Context) {

    file1 := ""
    file2 := ""
    if c.NArg() > 0 {
        file1 = c.Args().Get(0)
        file2 = c.Args().Get(1)
    }

    if file1 == "" || file2 == "" {
        fmt.Println("you are required to pass in two file paths")
    }

    // read file
    arr1 := readCSV(file1)
    arr2 := readCSV(file2)
    fmt.Println("count of arr1 ", len(arr1))
    fmt.Println("count of arr2 ", len(arr2))

    // find commonalities in both arrays
    var commonalities []int

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

}
