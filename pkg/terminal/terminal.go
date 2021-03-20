package terminal

import (
    "errors"
    "flag"
    "os"

    "github.com/IveGotNorto/podimator/pkg/podimator"
)

func AllCommand() *Command {
    gc := &Command{
        sub: flag.NewFlagSet("all", flag.ContinueOnError),
    }
    gc.sub.StringVar(&gc.podcast, "name", "", "name of the podcast to be processed")
    gc.sub.StringVar(&gc.config, "config", "podcast.json", "configuration path location")
    gc.sub.BoolVar(&gc.verbose, "v", false, "change program output level")
    return gc
}

func UpdateCommand() *Command {
    gc := &Command{
        sub: flag.NewFlagSet("update", flag.ContinueOnError),
    }
    gc.sub.StringVar(&gc.podcast, "name", "", "name of the podcast to be processed")
    gc.sub.StringVar(&gc.config, "config", "podcast.json", "configuration path location")
    gc.sub.BoolVar(&gc.verbose, "v", false, "change program output level")
    return gc
}

type Command struct {
    sub *flag.FlagSet
    podcast string
    config string
    verbose bool
}

func (g *Command) Name() string {
    return g.sub.Name()
}

func (g *Command) Init(args []string) error {
    return g.sub.Parse(args)
}

func (g *Command) Run() error {
    if g.sub.Name() == "all" {
        podimator.All(g.name) 
    } else {
        podimator.Update(g.name) 
    }
    return nil
}

type Runner interface {
    Init([]string) error
    Name() string
    Run() error
}

func process(args []string) error {

    if len(args) < 1 {
        return errors.New("You must pass a sub-command")
    }

    cmds := []Runner{
        AllCommand(),
        UpdateCommand(),
    }

    subcommand := os.Args[1]

    for _, cmd := range cmds {
        if cmd.Name() == subcommand {
            cmd.Init(os.Args[2:])
            return cmd.Run()
        }
    }
    return nil
}

func Get() error {
    // TODO: Implement error handling and usage output
    return process(os.Args[1:])
}

