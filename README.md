# weso

CLI for Websocket.

## install

```bash
go get github.com/morikuni/weso/cmd/weso
```

## Usage

```bash
$ cat template.txt
{{define "Name"}}{"name": "{{._1}}"}{{end}}

{{define "Decorate"}}{{._1}}{{._2}}{{._1}}{{end}}

{{define "Twice"}}{{._1}}{{._1}}{{end}}

$ weso -template template.txt ws://echo.websocket.org
> hello
<< hello
> .Name Jack
<< {"name": "Jack"}
> .Decorate *** "I'm Jack"
<< ***I'm Jack***
> .Twice Hello
<< HelloHello
```

Start line with dot `.` to use template.

## Template

Template file is written in [Go's standard template library](https://golang.org/pkg/text/template/).

Arguments for template are expanded to `._1`, `._2` ... `._N`.

## License

[MIT](http://github.com/morikuni/weso/LICENSE)

## Credit

weso uses

- [Liner](https://github.com/peterh/liner) for line editor.
