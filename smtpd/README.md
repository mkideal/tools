# smtpd [![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/mkideal/tools/master/LICENSE)

## License

[The MIT License (MIT)](https://raw.githubusercontent.com/mkideal/tools/master/LICENSE)

## Install

```shell
go get github.com/mkideal/tools/smtpd
```

## Usage

using `smtpd -h` show help information.

```
exp [-h | --help]
exp [--debug] [-H | --host=HOST] [-p | --port=PORT]
```

**Opts**

```
--debug
	enable debug mode

-H, --host=[0.0.0.0]
	local host

-p, --port=[8080]
	listening port for smtpd
```

## Examples

```shell
smtpd
smtpd -p 25
smtpd --debug -H 127.0.0.1 -p 25
```
