package terminal

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/IveGotNorto/podimator/podimator"
)

func Run() {
	pod := podimator.New()

	app := &cli.App{
		Name:  "Podimator",
		Usage: "Automated podcast downloader.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
				Value:   "~/.config/podimator/podcasts.toml",
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "Debugging program output",
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Verbose program output",
				Value:   false,
			},
		},
		Before: func(c *cli.Context) error {
			pod.Debug = c.Bool("debug")
			pod.Verbose = c.Bool("verbose")
			pod.ConfigPath = c.String("config")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "update",
				Aliases: []string{"u"},
				Usage:   "retrieve most recent episode",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "Select podcast by name",
						Value:   "",
					},
				},
				Action: func(c *cli.Context) error {
					pod.Start(podimator.Update{
						PodcastName: c.String("name"),
					})
					return nil
				},
			},
			{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "retrieve all episodes",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "Select podcast by name",
						Value:   "",
					},
				},
				Action: func(c *cli.Context) error {
					pod.Start(podimator.All{
						PodcastName: c.String("name"),
					})
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}
