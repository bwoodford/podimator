package terminal

import (
    "errors"
    "flag"
    "os"

    "github.com/IveGotNorto/podimator"
)

func AllCommand() *Command {
    gc := &Command{
        sub: flag.NewFlagSet("all", flag.ContinueOnError),
    }
    gc.sub.StringVar(&gc.podcast, "podcast", "", "name of the podcast to be processed")
    return gc
}

func UpdateCommand() *Command {
    gc := &Command{
        sub: flag.NewFlagSet("update", flag.ContinueOnError),
    }
    gc.sub.StringVar(&gc.podcast, "podcast", "", "name of the podcast to be processed")
    return gc
}

type Command struct {
    sub *flag.FlagSet
    podcast string
}

func (g *Command) Name() string {
    return g.sub.Name()
}

func (g *Command) Init(args []string) error {
    return g.sub.Parse(args)
}

func (g *Command) Run() error {
    if g.sub.Name() == "all" {
        podimator.All(g.podcast) 
    } else {
        podimator.Update(g.podcast) 
    }
    return nil
}

type Runner interface {
    Init([]string) error
    Name() string
    Run() error
}

func process(args []string) error {

    var subcommand string

    if len(args) < 1 {
        // Default command
        subcommand = "update"
    }

    subcommand := os.Args[1]

    cmds := []Runner{
        AllCommand(),
        UpdateCommand(),
    }

    for _, cmd := range cmds {
        if cmd.Name() == subcommand {
            cmd.Init(os.Args[2:])
            return cmd.Run()
        }
    }
    return nil
}

func Run() error {
    // TODO: Implement error handling and usage output
    return process(os.Args[1:])
}

