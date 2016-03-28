package main

import (
	"fmt"
	"net/http"

	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	cli.Addr
	Maps map[string]string `cli:"M" usage:"addr mapping"`
}

func run(ctx *cli.Context, argv *argT) error {
	if argv.Maps == nil {
		argv.Maps = make(map[string]string)
	}
	ctx.JSONIndentln(argv.Maps, "", "    ")
	addr := fmt.Sprintf("%s:%d", argv.Host, argv.Port)
	return http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		readdr, ok := argv.Maps[req.Host]
		if !ok {
			fmt.Fprintf(w, "could not redirect %q to other address", req.Host)
			return
		}
		urlStr := fmt.Sprintf("%s://%s%s?%s", "http", readdr, req.URL.Path, req.URL.RawQuery)
		ctx.String("redirect from %q to %q\n", req.URL.Path, urlStr)
		http.Redirect(w, req, urlStr, http.StatusFound)
	}))
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
