package main

import (
	"regexp"
	"time"

	"github.com/mkideal/pkg/debug"
)

func init() {
	registerParser("io", newIOParser())
}

type ioParser struct {
	expiryRegexp *regexp.Regexp
}

func newIOParser() *ioParser {
	parser := new(ioParser)
	parser.expiryRegexp = regexp.MustCompile(`Expiry : ([0-9\-]+)`)
	return parser
}

func (parser *ioParser) Suffix() string {
	return ".io"
}

func (parser *ioParser) Parse(domain, text string) Result {
	debug.Debugf("domain: %s, text: `%s`", domain, text)
	res := Result{}
	res.Domain = domain
	findRes := parser.expiryRegexp.FindStringSubmatch(text)
	debug.Debugf("findRes: %q\n", findRes)
	if findRes != nil && len(findRes) == 2 {
		res.Expiry, _ = time.Parse("2006-01-02", findRes[1])
	}
	return res
}
