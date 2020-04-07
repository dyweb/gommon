# mitchellh/cli

- `Command` is a interface with `Run`
- when register it uses string instead of nested command struct like cobra, e.g. `foo bar`
- has [CommandFactory](https://pkg.go.dev/github.com/mitchellh/cli?tab=doc#CommandFactory)
- need to handle flag manually using standard library, e.g. [consul kv get](https://github.com/hashicorp/consul/blob/master/command/kv/get/kv_get.go)
- use https://github.com/posener/complete for completion generation

If you use a CLI with nested subcommands, some semantics change due to ambiguities

> We use longest prefix matching to find a matching subcommand. This
    means if you register "foo bar" and the user executes "cli foo qux",
    the "foo" command will be executed with the arg "qux". It is up to
    you to handle these args. One option is to just return the special
    help return code `RunResultHelp` to display help and exit.

> Any parent commands that don't exist are automatically created as
    no-op commands that just show help for other subcommands. For example,
    if you only register "foo bar", then "foo" is automatically created.


```go
package main

func main() {
    c := cli.NewCLI("app", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"foo": fooCommandFactory,
		"bar": barCommandFactory,
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}
}
```