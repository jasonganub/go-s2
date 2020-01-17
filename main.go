package main

import (
    "fmt"
    "github.com/urfave/cli"
    "go-s2/diffplace"
    "go-s2/g2s2"
    "log"
    "os"
)


func main() {
    app := cli.NewApp()
    app.Name = "go-s2"
    app.Usage = "Spherical Geometry utility collections"
    app.Version = "0.1"
    app.Action = func(c *cli.Context) error {
        fmt.Println("Try 'help'")
        return nil
    }

    app.Commands = []cli.Command{
        {
            Name:    "diffplace",
            Aliases: []string{"d"},
            Usage:   "Diff two S2ID files and move commonalities from first to the second file",
            Action: func(c *cli.Context) error {
                diffplace.Run(c)
                return nil
            },
        },
        {
            Name:    "geojson2s2ids",
            Aliases: []string{"g2s2"},
            Usage:   "Convert geojson FeatureCollection to set of S2IDs",
            Action:  func (c *cli.Context) error {
                g2s2.Run(c)
                return nil
            },
        },
    }

    err := app.Run(os.Args)
    if err != nil {
        log.Fatal(err)
    }
}
