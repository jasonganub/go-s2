package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/urfave/cli"
)

const (
	text1 = "Lorem ipsum dolor."
	text2 = "Lorem dolor sit amet."
)

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
			Usage:   "diff two s2 id files and move unions from first to the second file",
			Action:  func(c *cli.Context) error {
				fmt.Print("DIFFPLACE!")
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}


	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(text1, text2, false)

	fmt.Println(dmp.DiffPrettyText(diffs))



}