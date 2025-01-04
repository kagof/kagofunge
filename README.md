# kagofunge
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/kagof/kagofunge/test.yml)


An interpreter & debugger written in Go for the Befunge-93 esoteric programming language.

This was mainly done as a pet project to get more familiar with Go, so please use with caution.

## Building

This project is written in Go 1.23. To compile it into the current directory, simply run:
```sh
go build
```

or to compile it and add it to your `$GOPATH/bin`:

```sh
go install
```

You may optionally want to link it to a shorthand `kgf` if you (for some reason) are using it frequently. This can be done with the following.

On a *nix shell:
```sh
ln -s $(go env GOPATH)/bin/kagofunge $(go env GOPATH)/bin/kgf
```

In PowerShell:
```pwsh
$GOPATH = (go env GOPATH)
New-Item -Path "$GOPATH\bin\kgf.exe" -ItemType SymbolicLink -Target "$GOPATH\bin\kagofunge.exe"
```

### Dependencies

The only non-stdlib direct dependencies of this program are:

* [`github.com/fatih/color v1.18.0`](https://github.com/fatih/color) used to color the debugger output
* [`github.com/spf13/cobra v1.8.1`](https://github.com/spf13/cobra) used for the CLI
* [`github.com/stretchr/testify v1.10.0`](https://github.com/stretchr/testify) used for assertions in tests
* [`github.com/goccy/go-yaml v1.15.13`](https://github.com/goccy/go-yaml) used for parsing YAML config files

## Usage

`kagofunge` is an interpreter and debugger for Befunge-93 written in Go.
For detailed usage, use `kagofunge run --help` or `kagofunge debug --help`.

```sh
kagofunge <run|debug> <program> [flags]
```

### Examples

```sh
kagofunge run hello-world.bf
kagofunge run '<> #,:# _@#:"Hello, World!"' -I
kagofunge run hello-world.bf -o output.txt -i input.txt
```

```sh
kagofunge debug hello-world.bf --breakpoint "(0,0)"
kagofunge debug '<> #,:# _@#:"Hello, World!"' -I -b 0,0 -b 8,0
kagofunge debug hello-world.bf -o output.txt -i input.txt -b '[1,1]'
```

### Available Sub-Commands

| Name    | Description                |
|---------|----------------------------|
| `debug` | Debug a Befunge-93 program |
| `run`   | Run a Befunge-93 program   | 

### Flags

#### global
| Shortcut | Name            | Type      | Repeatable | Description                                                                                                                                      |
|----------|-----------------|-----------|------------|--------------------------------------------------------------------------------------------------------------------------------------------------|
| `-h`     | `--help`        | boolean   | false      | help for the given command                                                                                                                       |
| `-I`     | `--inline`      | boolean   | false      | If set, then the `<program>` is interpreted as an inline Befunge-93 program, otherwise it is interpreted as a path to a Befunge-93 program file. |
| `-i`     | `--input`       | string    | false      | Output file path. Default: `stdin`                                                                                                               |
| `-o`     | `--output`      | string    | false      | Output file path. Default: `stdout`                                                                                                              |
| `-c`     | `--config`      | key=value | true       | Override specific config values as key=value pairs.                                                                                              |
| `-C`     | `--config-file` | string    | false      | Config file path. Default: `$KAGOFUNGE_CONFIG_FILE` if set or `$HOME/.kgf/config.yml` if not                                                     |

#### root command only
| Shortcut | Name        | Type    | Repeatable | Description             |
|----------|-------------|---------|------------|-------------------------|
| `-v`     | `--version` | boolean | false      | version for `kagofunge` |

#### debug sub-command only
| Shortcut | Name           | type        | Repeatable | Description                                                                                                            |
|----------|----------------|-------------|------------|------------------------------------------------------------------------------------------------------------------------|
| `-b`     | `--breakpoint` | stringArray | true       | Breakpoints to set in the program while executing. can be in the formats `(x,y)`, `(x y)`, `[x,y]`, `[x y]`, or `x,y`. |
| `-s`     | `--speed`      | duration    | false      | If set, the program will progress automatically at the specified speed. Should be a duration. Eg 100ms, 1s             |s

### Debugging

There is a Terminal based debugger which can be used to step through the program, see the state of the code and the stack, and generally see how a Befunge-93 program is executing.

![debugging demo](img/_debug_demo.gif)

## Testing

The Go test suite can be executed by running

```sh
go test ./...
```

## About Befunge

Befunge-93 is an esoteric programming language created by [cpressey](https://catseye.tc/). It is a stack-based language operating on a two-dimensional plane (technically a torus) where the program counter's direction is determined by certain control characters (`>`, `^`, `<`, `v`). Other control characters include conditional directions (`|`, `_`), and skips (`#`). It also allows for modification of the code at runtime (`p`), which can lead to interesting results. To learn more, the [Wikipedia page](https://en.wikipedia.org/wiki/Befunge) has a good summary of how the language operates.

Note that although the spec for Befunge-93 states that the program should only be 80x25, this interpreter does not impose such a restriction (which seems to be common among other implementations as well).

An example "Hello World" program in Befunge-93 can look like:

```Befunge
>               v
v"Hello, World!"<
> #,:# _@
```
The program counter starts in the upper left, then is directed down and moves right-to-left through the second row. `"` turns on "string mode" so the characters in `Hello, World!` get added to the stack backwards. `"` turns string mode off again. Finally, `> #,:# _@` is an outputting loop, printing all values on the stack as ASCII until the stack is empty, at which point the program terminates.

### Funge-98

There is also a more ambitious spec for Funge-98, which defines many more control characters, uses a stack of stacks, and allows for different dimensional "funges" (unefunge, trefunge, etc). This is much more difficult to implement and also doesn't have as well defined behaviour as Befunge-93, so I've limited this project to the 93 spec.

## Configuration

The kagofunge interpreter and debugger can both be configured to change the behaviour. Some of these changes accommodate areas of the Befunge-93 spec that are left to interpretation or historically have had many different possibilities.

There are several ways that configuration is picked up. In order from lowest precedence to highest, they are:

* built in defaults
* configuration file
* environment variables
* `-c`/`--config` flag overrides

### Configuration values

| parent      | name                           | possible values                                                                                        | description                                                                                                                                                                                                                                                                                                                              |
|-------------|--------------------------------|--------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| interpreter | divide-by-zero-behaviour       | <ul><li>`PROMPT_FOR_INPUT` (default)</li><li>`RETURN_ZERO`</li><li>`REFLECT`</li><li>`PANIC`</li></ul> | Behaviour when dividing by zero. The default for Befunge-93 is to prompt the user for a value, however other possibilities are to either push 0 onto the stack, reflect the instruction pointer, or panic (exit the program with an error).                                                                                              |
| interpreter | modulus-by-zero-behaviour      | <ul><li>`PROMPT_FOR_INPUT` (default)</li><li>`RETURN_ZERO`</li><li>`REFLECT`</li><li>`PANIC`</li></ul> | Behaviour when performing modulus by zero. The default for Befunge-93 is to prompt the user for a value, however other possibilities are to either push 0 onto the stack, reflect the instruction pointer, or panic (exit the program with an error).                                                                                    |
| interpreter | put-out-of-bounds-behaviour    | <ul><li>`NO_OP` (default)</li><li>`ZERO`</li><li>`WRAP`</li><li>`PANIC`</li></ul>                      | Behaviour when performing the `p` command with coordinates that lie outside of the torus. The default behaviour for Befunge-93 is to do nothing, however you can also choose to wrap the value across the torus, or panic (exit the program with an error). Note that `ZERO` is meaningless for `p` and will behave the same as `NO_OP`. |
| interpreter | get-out-of-bounds-behaviour    | <ul><li>`NO_OP`</li><li>`ZERO` (default)</li><li>`WRAP`</li><li>`PANIC`</li></ul>                      | Behaviour when performing the `g` command with coordinates that lie outside of the torus. The default behaviour for Befunge-93 is to add 0 to the stack, however you can also choose to do nothing, wrap the value across the torus, or panic (exit the program with an error).                                                          |
| interpreter | enforce-torus-size-restriction | <ul><li>`true`</li><li>`false` (default)</li></ul>                                                     | Whether or not to enforce the torus size restriction. Traditionally, Befunge-93 programs can only be 80x25 characters, though many interpreters ignore this restriction (including this one by default). If set to true, then program inputs will be truncated or padded to fit the size restriction.                                    |
| interpreter | torus-size-restriction-width   | integer > 0 (default 80)                                                                               | If enforce-torus-size-restriction is true, the width to restrict the torus to.                                                                                                                                                                                                                                                           |
| interpreter | torus-size-restriction-height  | integer > 0 (default 25)                                                                               | If enforce-torus-size-restriction is true, the height to restrict the torus to.                                                                                                                                                                                                                                                          |
| debugger    | show-torus                     | <ul><li>`true` (default)</li><li>`false`</li></ul>                                                     | Whether or not to show the code torus in the debugger output.                                                                                                                                                                                                                                                                            |
| debugger    | show-torus-coordinates         | <ul><li>`true` (default)</li><li>`false`</li></ul>                                                     | If showing the code torus, whether or not to show the coordinates.                                                                                                                                                                                                                                                                       |
| debugger    | show-stack                     | <ul><li>`true` (default)</li><li>`false`</li></ul>                                                     | Whether or not to show the stack in the debugger output.                                                                                                                                                                                                                                                                                 |
| debugger    | enable-colors                  | <ul><li>`true` (default)</li><li>`false`</li></ul>                                                     | Whether or not to use ANSI colors in the debugger output.                                                                                                                                                                                                                                                                                |

### Configuration file

The format of the configuration file can be YAML or JSON. It must adhere to [this JSON schema](./kagofunge-config.schema.json). An example filled out with default values can be found [here](default-kgf-config.yaml) and adapted to suit your needs. This is either the file pointed to by the `-C`/`--config-file` flag, the `$KGF_CONFIG_FILE` environment variable, or the `$HOME/.kgf/config.yaml` file (in that order). Note that at most one of these three will be used.

### Environment variables

The format for environment variable names is `KGF_$parent_$name`, changing sausage-case to SCREAMING_SNAKE_CASE. For example, `KGF_INTERPRETER_DIVIDE_BY_ZERO_BEHAVIOUR`.

### Flag overrides

The format for flag overrides is just `$parent.$name`, for example `interpreter.divide-by-zero-behaviour`.

## See Also

* the [Befunge-93 reference implementation](https://codeberg.org/catseye/Befunge-93)
* the [Befunge-93 spec](https://codeberg.org/catseye/Befunge-93/src/branch/master/doc/Befunge-93.markdown)
* my repo of [Befunge-93 programs](https://github.com/kagof/BefungeRepo)
* my (largely abandoned) Befunge-93 [extension for VSCode](https://marketplace.visualstudio.com/items?itemName=kagof.befunge)
* [gofunge98](https://github.com/adyxax/gofunge98), an unrelated Funge-98 intepretter in Go. This project has been renamed from Gofunge to kagofunge to avoid any confusion