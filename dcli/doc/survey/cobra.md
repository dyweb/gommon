# Cobra

https://github.com/spf13/cobra

- support persistent flag, i.e. flag that can be applied to sub commands `gommon gen --foo bar --verbose`
- https://github.com/at15/code-i-read/tree/master/go/spf13/cobra has much more detail

## Subcommand

If you defined a subcommand like `git clone` and you do `git glone`,
it will error out instead of passing `clone` as argument to the `git` command.
However if you just use `git`, it will call the `git` command without argument.

```text
// pseudo code for execute logic in cobra
func ExecuteC() {
    commands = stripFlags(args)
    nextCommand = commands[0]
    var cmd
    for _, c := range c.commands {
        if c.Name() == nextCommand {
            cmd = c
            break
        }
    }
    cmd.execute()
}
```

```go
// https://github.com/spf13/cobra/blob/master/command.go#L605

// Find the target command given the args and command tree
// Meant to be run on the highest node. Only searches down.
func (c *Command) Find(args []string) (*Command, []string, error) {
	var innerfind func(*Command, []string) (*Command, []string)

	innerfind = func(c *Command, innerArgs []string) (*Command, []string) {
		argsWOflags := stripFlags(innerArgs, c)
		if len(argsWOflags) == 0 {
			return c, innerArgs
		}
		nextSubCmd := argsWOflags[0]

		cmd := c.findNext(nextSubCmd)
		if cmd != nil {
			return innerfind(cmd, argsMinusFirstX(innerArgs, nextSubCmd))
		}
		return c, innerArgs
	}

	commandFound, a := innerfind(c, args)
	if commandFound.Args == nil {
		return commandFound, a, legacyArgs(commandFound, stripFlags(a, commandFound))
	}
	return commandFound, a, nil
}

// https://github.com/spf13/cobra/blob/6607e6b8603f56adb027298ee6695e06ffb3a819/command.go#L546
func stripFlags(args []string, c *Command) []string {
	if len(args) == 0 {
		return args
	}
	c.mergePersistentFlags()

	commands := []string{}
	flags := c.Flags()

Loop:
	for len(args) > 0 {
		s := args[0]
		args = args[1:]
		switch {
		case s == "--":
			// "--" terminates the flags
			break Loop
		case strings.HasPrefix(s, "--") && !strings.Contains(s, "=") && !hasNoOptDefVal(s[2:], flags):
			// If '--flag arg' then
			// delete arg from args.
			fallthrough // (do the same as below)
		case strings.HasPrefix(s, "-") && !strings.Contains(s, "=") && len(s) == 2 && !shortHasNoOptDefVal(s[1:], flags):
			// If '-f arg' then
			// delete 'arg' from args or break the loop if len(args) <= 1.
			if len(args) <= 1 {
				break Loop
			} else {
				args = args[1:]
				continue
			}
		case s != "" && !strings.HasPrefix(s, "-"):
			commands = append(commands, s)
		}
	}

	return commands
}
```