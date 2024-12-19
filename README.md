# kagofunge

An interpreter & debugger written in Go for the Befunge-93 esoteric programming language.

This was mainly done as a pet project to get more familiar with Go, so please use with caution.

## Building

This project is written in Go 1.23. To build it, simply run
```sh
go build
```

The only non-stdlib dependency of this program is [`github.com/fatih/color v1.18.0`](https://github.com/fatih/color) which is used to color the debugger output.

## Usage
Once you have a binary, the usage is as follows:
```
Usage of kagofunge:
kagofunge [OPTIONS] befungeFile

Examples:
         kagofunge filename.bf
         kagofunge -debug -breakpoint '(0,0)' filename.bf
         kagofunge -debug -breakpoint 0,0 -breakpoint 1,2 filename.bf
         kagofunge -i input.txt -o output.txt filename.bf
         kagofunge -inline '"olleh",,,,,@'

  -breakpoint value
        Breakpoints in the program, if in debug mode. Multiple supported
  -debug
        toggle on/off debug mode
  -i string
        input file (default: stdin)
  -inline string
    	an inline Befunge-93 program
  -o string
        output file (default: stdout)


```

### Debugging

There is a Terminal based debugger which can be used to step through the program, see the state of the code and the stack, and generally see how a Befunge-93 program is executing.

![debugging demo](img/_debug_demo.gif)

## About Befunge

Befunge-93 is an esoteric programming language created by [cpressey](https://catseye.tc/). It is a stack-based language operating on a two-dimensional plane (technically a torus) where the program counter's direction is determined by certain control characters (`>`, `^`, `<`, `v`). Other control characters include conditional directions (`|`, `_`), and skips (`#`). It also allows for modification of the code at runtime (`p`), which can lead to interesting results. To learn more, the [Wikipedia page](https://en.wikipedia.org/wiki/Befunge) has a good summary of how the language operates.

Note that although the spec for Befunge-93 states that the program should only be 80x25, this interpreter does not impose such a restriction (which seems to be common among other implementations as well).

An example "hello world" program in Befunge-93 can look like:
```befunge-93
>               v
v"hello, world!"<
>,:# _@
```
The program counter starts in the upper left, then is directed down and moves right-to-left through the second row. `"` turns on "string mode" so the characters in `hello, world!` get added to the stack backwards. `"` turns string mode off again. Finally, `>,:# _@` is an outputting loop, printing all values on the stack as ASCII until the stack is empty, at which point the program terminates.

### Funge-98

There is also a more ambitious spec for Funge-98, which defines many more control characters, uses a stack of stacks, and allows for different dimensional "funges" (unefunge, trefunge, etc). This is much more difficult to implement and also doesn't have as well defined behaviour as Befunge-93, so I've limited this project to the 93 spec.

## See Also

* the [Befunge-93 reference implementation](https://codeberg.org/catseye/Befunge-93)
* the [Befunge-93 spec](https://codeberg.org/catseye/Befunge-93/src/branch/master/doc/Befunge-93.markdown)
* my repo of [Befunge-93 programs](https://github.com/kagof/BefungeRepo)
* my (largely abandoned) Befunge-93 [extension for VSCode](https://marketplace.visualstudio.com/items?itemName=kagof.befunge)
* [gofunge98](https://github.com/adyxax/gofunge98), an unrelated Funge-98 intepretter in Go. This project has been renamed from Gofunge to kagofunge to avoid any confusion