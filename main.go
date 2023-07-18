package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/urfave/cli/v2"

	awsresq "github.com/thaim/awsresq/internal"
)

var (
	version = "main"
	region string
	service string
	resource string
)

func main() {
	app := &cli.App{
		Name:  "awsresq",
		Usage: "search resources on AWS",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "region",
				Usage: "region name",
				Destination: &region,
			},
			&cli.StringFlag{
				Name: "service",
				Usage: "service name",
				Destination: &service,
			},
			&cli.StringFlag{
				Name: "resource",
				Usage: "resource name",
				Destination: &resource,
			},
		},
		Action: func(ctx *cli.Context) error {
			query := ctx.Args().Get(0)
			client, err := awsresq.NewAwsresqClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "initialized failed:%v\n", err)
				os.Exit(1)
			}
			if res, err := client.Search(service, resource, query); err != nil {
				fmt.Fprintf(os.Stderr, res)
			}
			return err
		},
		HideHelpCommand: true,
		Version: getVersion(),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	os.Exit(0)
}

func getVersion() string {
	if version != "" {
		return version
	}
	i, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}

	return i.Main.Version
}

