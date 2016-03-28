package main

import (
	"fmt"
	"net/http"

	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	cli.Addr
	Maps map[string]uint16 `cli:"M" usage:"addr mapping"`
}

func (argv *argT) Validate(ctx *cli.Context) error {
	//TODO: validate something or remove this function.
	return nil
}

func run(ctx *cli.Context, argv *argT) error {
	if argv.Maps == nil {
		argv.Maps = make(map[string]uint16)
	}
	ctx.JSONIndentln(argv.Maps, "", "    ")
	addr := fmt.Sprintf("%s:%d", argv.Host, argv.Port)
	return http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx.String("req: %v", *req)
		port, ok := argv.Maps[req.Host]
		if !ok {
			fmt.Fprintf(w, "could not redirect %q to new address", req.Host)
			return
		}
		urlStr := fmt.Sprintf("%s://%s:%d%s?%s", "http", req.Host, port, req.URL.Path, req.URL.RawQuery)
		ctx.String("redirect from %q to %q\n", req.URL.Path, urlStr)
		http.Redirect(w, req, urlStr, http.StatusUseProxy)
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
