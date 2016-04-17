package main

import (
	"fmt"
	"regexp"
	"testing"
)

const result = `
Domain : 9n.io
Status : Client Updt+Delt Lock
Expiry : 2016-12-10

NS 1   : ns21.domaincontrol.com
NS 2   : ns22.domaincontrol.com`

func TestRegexpIO(t *testing.T) {
	expiryRegexp := regexp.MustCompile(`Expiry : ([0-9\-]+)`)
	ret := expiryRegexp.FindStringSubmatch(result)
	fmt.Printf("find result: %q\n", ret)
}
