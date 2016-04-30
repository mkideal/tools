package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/labstack/gommon/color"
	"github.com/mkideal/cli"
	"github.com/mkideal/pkg/textutil"
)

type argT struct {
	cli.Helper
	TemplateFile string `cli:"t,template" usage:"template file, if TPL is empty, then read template from stdin" name:"TPL"`

	Out string            `cli:"o,out" usage:"write to specified file instead of stdout"`
	M   map[string]string `cli:"M" usage:"map values(e.g. -Ma=A -Mkey=value)"`
}

func (argv *argT) Validate(ctx *cli.Context) error {
	return nil
}

func run(ctx *cli.Context, argv *argT) error {
	if argv.M == nil {
		argv.M = make(map[string]string)
	}
	var (
		t   *template.Template
		err error
	)
	if argv.TemplateFile == "" {
		// read template from stdin
		if bytes, err := ioutil.ReadAll(os.Stdin); err == nil {
			t = template.New("tmp")
			t, err = t.Parse(string(bytes))
		}
	} else {
		// read template from file
		t, err = template.ParseFiles(argv.TemplateFile)
	}
	if err != nil {
		return err
	}
	if t == nil {
		return fmt.Errorf("unknown error")
	}
	t.Option("missingkey=zero")
	// output file
	w := os.Stdout
	if argv.Out != "" {
		w, err = os.Create(argv.Out)
	}
	if err != nil {
		return err
	}
	return t.Execute(w, argv.M)
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
	}, Tpl(`{{.tpler}} is a template generator built by github.com/mkideal/cli

{{.Usage}}: tpler [-h | --help]
       tpler [-o | --out={{.OUT}}] [-M...] <-t | --tpl={{.TPL}}>

{{.Examples}}:
       tpler -t template.txt -Ma=1 -Mb=2
       echo "{{.hello}}" | tpler -Mhello=world
       echo "{{.hello}}" | tpler -Mhello=world -o out.txt`, map[string]string{
		"tpler":    color.Bold("tpler"),
		"Usage":    color.Bold("Usage"),
		"Examples": color.Bold("Examples"),
		"hello":    "{{.hello}}",
		"OUT":      color.Bold("OUT"),
		"TPL":      color.Bold("TPL"),
	}))
}
