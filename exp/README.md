# exp [![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/mkideal/tools/master/LICENSE)

## License

[The MIT License (MIT)](https://raw.githubusercontent.com/mkideal/tools/master/LICENSE)

## Install

```shell
go get github.com/tools/exp
```

## Usage

using `exp -h` show help information.

```
exp -h
exp [-e] [-i] [-f FILE] [-D...] [EXPR]
```

**Opts**

```
-e
	show native expression before expression result.
	`exp x -Dx=1`		=> `1`
	`exp -e x -Dx=1`	=> `x: 1`

-i
	read expression from stdin

-f FILE
	read expression from FILE

-D
	define variable. e.g. `-Dx=3` `-Da=1 -Db=2` `-D x=1 -D y=2`
```

## Examples

```shell
exp 1+2
exp -e 1+2
exp "1 + 2"
exp x -Dx=2.5
exp "x * y" -Dx=2 -Dy=6
exp "min(x, 4)" -Dx=3
exp "max(x, y, z)" -Dx=2 -Dy=6 -Dz=5
exp "rand() //rand in [0,10000)"
exp 'rand(n)' -Dn=100
exp 'rand(1,to)' -Dto=5
exp 'sum(1,2,3)'
exp 'aver(1,2,3)'
exp x y x+y x-y x*y x/y x%y x^y -Dx=7 -Dy=2
exp -e x y x+y x-y x*y x/y x%y x^y -Dx=7 -Dy=2
exp 'sin(pi)' 'sin(pi/2)'
exp e
exp pi
```
