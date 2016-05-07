package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/mkideal/cli"
	"github.com/mkideal/pkg/expr"
)

type argT struct {
	cli.Helper
	Variables map[string]float64 `cli:"D" usage:"define variables, e.g. -Dx=3 -Dy=4"`
	OutExpr   bool               `cli:"e" usage:"whther output native expression" dft:"false"`
	File      string             `cli:"f,file" usage:"read expr from file"`
	Stdin     bool               `cli:"i,stdin" usage:"read expr from stdin" sdt:"false"`

	args []string `cli:"-"`
}

func (argv *argT) Validate(ctx *cli.Context) error {
	argv.args = ctx.Args()

	dataList := make([]string, 0)
	if argv.File != "" {
		data, err := ioutil.ReadFile(argv.File)
		if err != nil {
			return err
		}
		dataList = append(dataList, string(data))
	}
	if argv.Stdin {
		if data, err := ioutil.ReadAll(os.Stdin); err != nil {
			return err
		} else if len(data) > 0 {
			dataList = append(dataList, string(data))
		}
	}
	for _, data := range dataList {
		args := strings.Split(strings.TrimSpace(data), "\n")
		for _, arg := range args {
			arg = strings.TrimFunc(arg, func(r rune) bool {
				return unicode.IsSpace(r) || r == '"' || r == '\''
			})
			if arg != "" {
				argv.args = append(argv.args, arg)
			}
		}
	}
	return nil
}

func run(ctx *cli.Context, argv *argT) error {
	rand.Seed(time.Now().UnixNano())
	if argv.Variables == nil {
		argv.Variables = make(map[string]float64)
	}
	getter := expr.Getter(argv.Variables)
	yellow := ctx.Color().Yellow

	for k, v := range reservedWords {
		if _, ok := getter[k]; ok {
			return fmt.Errorf("%s is reserved word", yellow(k))
		}
		getter[k] = v
	}

	for _, s := range argv.args {
		e, err := expr.New(s, pool)
		if err != nil {
			return err
		}
		ret, err := e.Eval(getter)
		if err != nil {
			return err
		}
		if argv.OutExpr {
			ctx.String("%s: ", s)
		}
		ctx.String("%.9G\n", ret)
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
	}, `exp - evaluate expressions
examples:
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
	exp x y x+y x-y x*y x/y x%%y x^y -Dx=7 -Dy=2
	exp -e x y x+y x-y x*y x/y x%%y x^y -Dx=7 -Dy=2`)
}

var reservedWords = map[string]float64{
	"e":  math.E,
	"E":  math.E,
	"pi": math.Pi,
	"PI": math.Pi,
}

var pool = func() *expr.Pool {
	p, err := expr.NewPool(map[string]expr.Func{
		"sum": func(args ...float64) (float64, error) {
			sum := float64(0)
			for _, arg := range args {
				sum += arg
			}
			return sum, nil
		},
		"aver": func(args ...float64) (float64, error) {
			n := len(args)
			if n == 0 {
				return 0, fmt.Errorf("missing arguments for function `%s`", "aver")
			}
			sum := float64(0)
			for _, arg := range args {
				sum += arg
			}
			return sum / float64(n), nil
		},
		"abs":   wrapOneArgumentFunc("abs", math.Abs),
		"acos":  wrapOneArgumentFunc("acos", math.Acos),
		"acosh": wrapOneArgumentFunc("acosh", math.Acosh),
		"asin":  wrapOneArgumentFunc("asin", math.Asin),
		"asinh": wrapOneArgumentFunc("asinh", math.Asinh),
		"atan":  wrapOneArgumentFunc("atan", math.Atan),
		"atanh": wrapOneArgumentFunc("atanh", math.Atanh),
		"cbrt":  wrapOneArgumentFunc("cbrt", math.Cbrt),
		"ceil":  wrapOneArgumentFunc("ceil", math.Ceil),
		"cos":   wrapOneArgumentFunc("cos", math.Cos),
		"cosh":  wrapOneArgumentFunc("cosh", math.Cosh),
		"e":     wrapOneArgumentFunc("e", math.Exp),
		"erf":   wrapOneArgumentFunc("erf", math.Erf),
		"erfc":  wrapOneArgumentFunc("erfc", math.Erfc),
		"floor": wrapOneArgumentFunc("floor", math.Floor),
		"gamma": wrapOneArgumentFunc("gamma", math.Gamma),
		"j0":    wrapOneArgumentFunc("j0", math.J0),
		"j1":    wrapOneArgumentFunc("j1", math.J1),
		"sin":   wrapOneArgumentFunc("sin", math.Sin),
		"sinh":  wrapOneArgumentFunc("sinh", math.Sinh),
		"sqrt":  wrapOneArgumentFunc("sqrt", math.Sqrt),
		"tan":   wrapOneArgumentFunc("tan", math.Tan),
		"tanh":  wrapOneArgumentFunc("tanh", math.Tanh),
		"trunc": wrapOneArgumentFunc("trunc", math.Trunc),
		"y0":    wrapOneArgumentFunc("y0", math.Y0),
		"y1":    wrapOneArgumentFunc("y1", math.Y1),
		"sgn": wrapOneArgumentFunc("sgn", func(x float64) float64 {
			if x > 0 {
				return 1
			}
			if x < 0 {
				return -1
			}
			return 0
		}),
		"ln": wrapOneArgumentFunc("ln", math.Log),

		"dim":   wrapTwoArgumentsFunc("dim", math.Dim),
		"log":   wrapTwoArgumentsFunc("log", func(x, y float64) float64 { return math.Log(x) / math.Log(y) }),
		"hypot": wrapTwoArgumentsFunc("hypot", math.Hypot),
		"jn":    wrapTwoArgumentsFunc("jn", func(n, x float64) float64 { return math.Jn(int(n), x) }),
		"yn":    wrapTwoArgumentsFunc("yn", func(n, x float64) float64 { return math.Yn(int(n), x) }),
		"mod":   wrapTwoArgumentsFunc("mod", math.Mod),
	})
	if err != nil {
		panic(err)
	}
	return p
}()

func argumentsSizeOne(name string, fn expr.Func) expr.Func {
	return argumentsSizeN(name, 1, fn)
}

func argumentsSizeN(name string, n int, fn expr.Func) expr.Func {
	return argumentsSizeRange(name, n, n, fn)
}

func argumentsSizeRange(name string, m, n int, fn expr.Func) expr.Func {
	return func(args ...float64) (float64, error) {
		if len(args) < m || len(args) > n {
			return 0, fmt.Errorf("bad arguments size for function `%s`", name)
		}
		return fn(args...)
	}
}

func wrapOneArgumentFunc(name string, fn func(float64) float64) expr.Func {
	return argumentsSizeOne(name, func(args ...float64) (float64, error) {
		return fn(args[0]), nil
	})
}

func wrapTwoArgumentsFunc(name string, fn func(float64, float64) float64) expr.Func {
	return argumentsSizeN(name, 2, func(args ...float64) (float64, error) {
		return fn(args[0], args[1]), nil
	})
}
