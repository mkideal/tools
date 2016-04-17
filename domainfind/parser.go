package main

import (
	"math"
	"os/exec"
	"time"

	"github.com/mkideal/cli"
	"github.com/mkideal/pkg/debug"
)

type Result struct {
	Domain string
	Whois  string
	Expiry time.Time
}

type Parser interface {
	Parse(name, text string) Result
	Suffix() string
}

var parserRepos = map[string]Parser{}

func registerParser(typ string, parser Parser) {
	if _, ok := parserRepos[typ]; ok {
		debug.Panicf("domain parser %s has exsited", typ)
	}
	parserRepos[typ] = parser
}

func findParser(typ string) Parser {
	if parser, ok := parserRepos[typ]; ok {
		return parser
	}
	return nil
}

const (
	Chars  = "abcdefghijklmnopqrstuvwxyz"
	Digits = "0123456789"
)

func whois(ctx *cli.Context, cfg Config, parser Parser) {
	set := []byte{}
	if cfg.OnlyChar {
		set = []byte(Chars)
	} else if cfg.OnlyDigit {
		set = []byte(Digits)
	} else {
		set = append([]byte(Chars), []byte(Digits)...)
	}
	length := len(set)

	nSet := make([][]string, cfg.MaxLen)
	for i := 0; i < cfg.MaxLen; i++ {
		l := int(math.Pow(float64(length), float64(i+1))) + 1
		nSet[i] = make([]string, 0, l)
		if i == 0 {
			for _, c := range set {
				nSet[i] = append(nSet[i], string([]byte{c}))
			}
		}
	}
	var (
		now    = time.Now()
		yellow = ctx.Color().Yellow
	)
	for i := 0; i < cfg.MaxLen; i++ {
		seti := nSet[i]
		for _, s := range seti {
			if i+1 != cfg.MaxLen {
				for _, c := range set {
					next := s + string([]byte{c})
					nSet[i+1] = append(nSet[i+1], next)
				}
			}
			domain := s + parser.Suffix()
			cmd := exec.Command("whois", domain)
			out, err := cmd.Output()
			if err != nil {
				ctx.String("%s\n", yellow(err.Error()))
				continue
			}
			res := parser.Parse(domain, string(out))
			if res.Expiry.Before(now) {
				ctx.String("%s - %s\n", res.Domain, res.Expiry.Format(time.RFC3339))
			}
		}
	}
}
