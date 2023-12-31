package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/urfave/cli/v2"

	awsresq "github.com/thaim/awsresq/internal"
)

var (
	version  = "main"
	region   string
	service  string
	resource string
)

func main() {
	app := &cli.App{
		Name:  "awsresq",
		Usage: "search resources on AWS",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "region",
				Usage:       "region name",
				Destination: &region,
			},
			&cli.StringFlag{
				Name:        "service",
				Usage:       "service name",
				Destination: &service,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "resource",
				Usage:       "resource name",
				Destination: &resource,
			},
		},
		Action: func(ctx *cli.Context) error {
			client, err := awsresq.NewAwsresqClient(region, service)
			if err != nil {
				fmt.Fprintf(os.Stderr, "initialized failed:%v\n", err)
				os.Exit(1)
			}

			validate := client.Validate(resource)
			if !validate {
				fmt.Fprintf(os.Stderr, "resource '%s' not supported in service '%s'\n", resource, service)
				os.Exit(1)
			}

			res, err := client.Search(service, resource)
			if err == nil {
				fmt.Fprint(os.Stdout, res+"\n")
			}
			return err
		},
		HideHelpCommand: true,
		Version:         getVersion(),
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
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
