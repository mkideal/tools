# trafficd [![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/mkideal/tools/master/LICENSE)

## License

[The MIT License (MIT)](https://raw.githubusercontent.com/mkideal/tools/master/LICENSE)

## Install

```shell
go get github.com/tools/trafficd
```

## Usage

using `trafficd -h` show help information.

```
trafficda -h | --help
trafficd [--host] [--port=8080] [-M...]
```

**Opts**

```
--host
	local host, default is empty same as 0.0.0.0.

--port
	trafficd listening port, default is 8080.

-M
	define map rules. e.g.
	`-M www.example.com=www.example.com:9090 example.com=www.example.com:9090`
```

## Examples

```shell
sudo trafficd --port=80 -M 127.0.0.1=127.0.0.1:8080
```

```shell
sudo trafficd --port=80 -M www.a.com=www.me.com:8080 -M www.b.com=www.me.com.9090
```

So, you can build many web sites in a same host.
