package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/mkideal/cli"
	"github.com/mkideal/pkg/debug"
)

func main() {
	if err := cli.Root(root,
		cli.Tree(help),
		cli.Tree(daemon),
		cli.Tree(api,
			cli.Tree(ping),
			cli.Tree(website),
		),
	).Run(os.Args[1:]); err != nil {
		fmt.Println(err)
	}
}

//------
// root
//------
var root = &cli.Command{
	Fn: func(ctx *cli.Context) error {
		ctx.WriteUsage()
		return nil
	},
}

//------
// help
//------
var help = &cli.Command{
	Name:        "help",
	Desc:        "display help",
	CanSubRoute: true,
	HTTPRouters: []string{"/v1/help"},
	HTTPMethods: []string{"GET"},

	Fn: cli.HelpCommandFn,
}

//--------
// daemon
//--------
type daemonT struct {
	cli.Helper
	cli.AddrWithShort
	Debug    bool   `cli:"debug" usage:"enable debug mode" dft:"false"`
	DBSource string `cli:"ds,db-source" usage:"database source" dft:"$APP_STAT_DB_SOURCE"`
}

func (t *daemonT) Validate(ctx *cli.Context) error {
	if t.Port == 0 {
		return fmt.Errorf("please don't use 0 as http port")
	}
	return nil
}

var daemon = &cli.Command{
	Name: "daemon",
	Desc: "startup app as daemon",
	Argv: func() interface{} { return new(daemonT) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*daemonT)
		if argv.Help {
			ctx.WriteUsage()
			return nil
		}
		debug.Switch(argv.Debug)

		addr := fmt.Sprintf("%s:%d", argv.Host, argv.Port)
		debug.Debugf("listening on: %s", addr)

		r := ctx.Command().Root()
		err := r.RegisterHTTP(ctx)
		if err != nil {
			return err
		}
		repo, err = Mysql(argv.DBSource)
		if err != nil {
			return err
		}
		return r.ListenAndServeHTTP(addr)
	},
}

//-----
// api
//-----
var api = &cli.Command{
	Name: "api",
	Desc: "display all api",
	Fn: func(ctx *cli.Context) error {
		ctx.String("Commands:\n")
		ctx.String("    ping\n")
		return nil
	},
}

//------
// ping
//------
var ping = &cli.Command{
	Name: "ping",
	Desc: "ping server",
	Fn: func(ctx *cli.Context) error {
		ctx.String("pong\n")
		return nil
	},
}

//---------
// website
//---------

type websiteT struct {
	cli.Helper
	AppId  string `cli:"appId" usage:"your app id"`
	AppKey string `cli:"appKey" usage:"your app key"`
	App    string `cli:"app" usage:"website name"`
	IP     string `cli:"ip" usage:"visitor's ip address"`
	CName  string `cli:"cname" usage:"visitor's country name"`
	Title  string `cli:"router" usage:"page's title"`
}

var website = &cli.Command{
	Name: "website",
	Desc: "statistic website access records",
	Argv: func() interface{} { return new(websiteT) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*websiteT)
		if argv.Help {
			ctx.String(ctx.Usage())
			return nil
		}
		debug.Debugf("%v", argv)
		if ctx.HTTPRequest != nil {
			ip := getip(ctx.HTTPRequest)
			if repo != nil {
				repo.SaveAccessRecord(argv.App, argv.Title, ip)
			}
		}
		ctx.String("(function(){})();")
		return nil
	},
}

func getip(req *http.Request) string {
	ips := getproxy(req)
	if len(ips) > 0 && ips[0] != "" {
		rip := strings.Split(ips[0], ":")
		if len(rip) > 0 {
			return rip[0]
		}
	}
	ip := strings.Split(req.RemoteAddr, ":")
	if len(ip) > 0 {
		if ip[0] != "[" {
			return ip[0]
		}
	}
	return "127.0.0.1"
}

func getproxy(req *http.Request) []string {
	if ips := req.Header.Get("X-Forwarded-For"); ips != "" {
		return strings.Split(ips, ",")
	}
	return []string{}
}

const createTableAccessRecord = "CREATE TABLE IF NOT EXISTS access_record (" +
	"`id` INT NOT NULL AUTO_INCREMENT," +
	"`app` varchar(128) NOT NULL," +
	"`title` varchar(128) NOT NULL," +
	"`ip` text NOT NULL," +
	"PRIMARY KEY ( id )" +
	")"

var repo *MysqlRepository

type MysqlRepository struct {
	locker sync.Mutex
	db     *sql.DB
}

func Mysql(dbsource string) (*MysqlRepository, error) {
	db, err := sql.Open("mysql", dbsource)
	if err != nil {
		return nil, err
	}

	repo := new(MysqlRepository)
	repo.db = db
	if err := multiExec(db,
		"CREATE DATABASE IF NOT EXISTS `app_stat`",
		"USE app_stat",
		createTableAccessRecord,
	); err != nil {
		return nil, err
	}
	return repo, nil
}

func (repo *MysqlRepository) SaveAccessRecord(app, title, ip string) error {
	repo.locker.Lock()
	defer repo.locker.Unlock()

	sqlStr := "insert into access_record(`app`,`title`,`ip`) values(?,?,?)"
	_, err := repo.db.Exec(sqlStr, app, title, ip)
	if err != nil {
		debug.Debugf("SaveAccessRecord error: %v", err)
	}
	return err
}

func multiExec(db *sql.DB, sqls ...string) error {
	for _, sql := range sqls {
		if _, err := db.Exec(sql); err != nil {
			return err
		}
	}
	return nil
}
