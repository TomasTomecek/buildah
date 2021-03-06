package main

import (
	"fmt"

	"encoding/json"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

type jsonImage struct {
	ID    string   `json:"id"`
	Names []string `json:"names"`
}

var (
	imagesFlags = []cli.Flag{
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "display only image IDs",
		},
		cli.BoolFlag{
			Name:  "noheading, n",
			Usage: "do not print column headings",
		},
		cli.BoolFlag{
			Name:  "notruncate",
			Usage: "do not truncate output",
		},
		cli.BoolFlag{
			Name:  "json",
			Usage: "output in JSON format",
		},
	}
	imagesDescription = "Lists locally stored images."
	imagesCommand     = cli.Command{
		Name:        "images",
		Usage:       "List images in local storage",
		Description: imagesDescription,
		Flags:       imagesFlags,
		Action:      imagesCmd,
		ArgsUsage:   " ",
	}
)

func imagesCmd(c *cli.Context) error {
	store, err := getStore(c)
	if err != nil {
		return err
	}

	images, err := store.Images()
	if err != nil {
		return errors.Wrapf(err, "error reading images")
	}

	quiet := false
	if c.IsSet("quiet") {
		quiet = c.Bool("quiet")
	}
	noheading := false
	if c.IsSet("noheading") {
		noheading = c.Bool("noheading")
	}
	truncate := true
	if c.IsSet("notruncate") {
		truncate = !c.Bool("notruncate")
	}
	if c.IsSet("json") {
		JSONImages := []jsonImage{}
		for _, image := range images {
			JSONImages = append(JSONImages, jsonImage{ID: image.ID, Names: image.Names})
		}
		data, err := json.MarshalIndent(JSONImages, "", "    ")
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", data)
		return nil
	}
	if len(images) > 0 && !noheading && !quiet {
		if truncate {
			fmt.Printf("%-12s %s\n", "IMAGE ID", "IMAGE NAME")
		} else {
			fmt.Printf("%-64s %s\n", "IMAGE ID", "IMAGE NAME")
		}
	}
	for _, image := range images {
		if quiet {
			fmt.Printf("%s\n", image.ID)
			continue
		}
		names := []string{""}
		if len(image.Names) > 0 {
			names = image.Names
		}
		for _, name := range names {
			if truncate {
				fmt.Printf("%-12.12s %s\n", image.ID, name)
			} else {
				fmt.Printf("%-64s %s\n", image.ID, name)
			}
		}
	}

	return nil
}
