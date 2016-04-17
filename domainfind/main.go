package main

import (
	"fmt"

	"github.com/mkideal/cli"
	"github.com/mkideal/pkg/debug"
)

type argT struct {
	cli.Helper
	Debug bool     `cli:"debug" usage:"开启调试模式" dft:"false"`
	Types []string `cli:"t,type" usage:"域名类型,com/io/org/..."`
	Config
}

type Config struct {
	MinLen    int  `cli:"min-len" usage:"最小长度" dft:"1"`
	MaxLen    int  `cli:"max-len" usage:"最大长度" dft:"3"`
	OnlyDigit bool `cli:"only-digit" usage:"仅包含数字" dft:"false"`
	OnlyChar  bool `cli:"only-char" usage:"仅包含字母" dft:"false"`
}

func (argv *argT) Validate(ctx *cli.Context) error {
	if argv.Types == nil || len(argv.Types) == 0 {
		return fmt.Errorf("types is empty")
	}
	if argv.OnlyChar && argv.OnlyDigit {
		return fmt.Errorf("--only-char and --only-digit both are true")
	}
	if argv.MinLen < 1 {
		argv.MinLen = 1
	}
	if argv.MaxLen > 10 {
		argv.MaxLen = 10
	}
	if argv.MinLen > argv.MaxLen {
		return fmt.Errorf("--min-len great than --max-len")
	}
	return nil
}

func run(ctx *cli.Context, argv *argT) error {
	debug.Switch(argv.Debug)

	parsers := make([]Parser, 0, len(argv.Types))
	yellow := ctx.Color().Yellow
	for _, typ := range argv.Types {
		parser := findParser(typ)
		if parser == nil {
			return fmt.Errorf("domain type %s unsupported", yellow(typ))
		}
		parsers = append(parsers, parser)
	}
	for _, parser := range parsers {
		whois(ctx, argv.Config, parser)
	}
	return nil
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		if argv.Help {
			ctx.WriteUsage()
			return nil
		}
		return run(ctx, argv)
	})
}
