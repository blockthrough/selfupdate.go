package cli

import "github.com/urfave/cli/v2"

type App = cli.App
type Command = cli.Command
type Context = cli.Context
type Flag = cli.Flag
type StringFlag = cli.StringFlag
type BoolFlag = cli.BoolFlag

var (
	Exit = cli.Exit
)
