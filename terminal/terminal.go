package terminal

import (
    "os"

    "github.com/urfave/cli/v2"

    "github.com/IveGotNorto/podimator/podimator"
)

func Run() {
    podi := podimator.New()

    app := &cli.App{
        Name: "Podimator",
        Usage: "Automated podcast downloader.",
        Flags: []cli.Flag{
            &cli.StringFlag{
                Name: "config",
                Aliases: []string{"c"},
                Usage: "Load configuration from `FILE`",
                Value: "~/.config/podimator/podcasts.json",
            },
            &cli.BoolFlag{
                Name: "debug",
                Aliases: []string{"d"},
                Usage: "Debugging program output",
                Value: false,
            },
            &cli.BoolFlag{
                Name: "verbose",
                Aliases: []string{"v"},
                Usage: "Verbose program output",
                Value: false,
            },
        },
        Before: func(c *cli.Context) error {
            // Update the program context
            podi.Debug = c.Bool("debug")
            podi.Verbose = c.Bool("verbose")
            podi.ConfigPath = c.String("config")
            return nil
        },
        Commands: []*cli.Command{
            {
                Name: "update",
                Aliases: []string{"u"},
                Usage: "retrieve most recent episode for podcasts",
                Flags: []cli.Flag{
                    &cli.StringFlag {
                        Name: "name",
                        Aliases: []string{"n"},
                        Usage: "Select single podcast by name",
                        Value: "",
                    },
                },
                Action: func(c *cli.Context) error {
                    podi.Run(podimator.Update{
                        PodcastName: c.String("name"),
                        EpisodeRange: "",
                        DateRange: "",
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
