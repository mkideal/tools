package main

import (
	"github.com/mkideal/cli"
	"github.com/mkideal/pkg/debug"
	"github.com/mkideal/pkg/smtpd"
)

type argT struct {
	cli.Helper
	cli.AddrWithShort
	Debug bool `cli:"debug" usage:"enable debug mode" dft:"false"`
}

func run(ctx *cli.Context, argv *argT) error {
	// switch debug mode
	debug.SwitchDebug(argv.Debug)

	// new smtp server
	svr := smtpd.NewServer()
	return svr.Start(argv.AddrWithShort.String(), func(ret string) {
		ctx.String(ret + "\n")
	})
}

func main() {
	cli.SetUsageStyle(cli.ManualStyle)
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		if argv.Help {
			ctx.WriteUsage()
			return nil
		}
		return run(ctx, argv)
	})
}
