// Test with `$ go build -o "go" . && ./go --help`

package main

import (
	"errors"
	"fmt"
	"os"

	flag "github.com/cosiner/flag"
)

type command interface {
	IsEnabled() bool
	Validate() error
	Run()
}

type buildParams struct {
	Enable   bool
	Already  bool     `names:"-a" important:"1" desc:"force rebuilding of packages that are already up-to-date."`
	Race     bool     `important:"1" desc:"enable data race detection.\nSupported only on linux/amd64, freebsd/amd64, darwin/amd64 and windows/amd64."`
	Output   string   `names:"-o" arglist:"output" important:"1" desc:"only allowed when compiling a single package"`
	LdFlags  string   `names:"-ldflags" arglist:"'flag list'" desc:"arguments to pass on each go tool link invocation."`
	Packages []string `args:"true"`
}

var buildMetadata = flag.Flag{
	Arglist: "[-o output] [-i] [build flags] [packages]",
	Desc: `
	Build compiles the packages named by the import paths,
	along with their dependencies, but it does not install the results.
	...
	The build flags are shared by the build, clean, get, install, list, run,
	and test commands:
		`,
}

func (p buildParams) IsEnabled() bool {
	return p.Enable
}

func (p buildParams) Validate() error {
	if len(p.Packages) == 0 {
		return errors.New("Error: you should at least specify one package")
	}
	return nil
}

func (p buildParams) Run() {
	fmt.Println("Going to build with the following parameters:")
	fmt.Println(p)
}

type cleanParams struct {
	Enable bool
}

func (p cleanParams) Validate() error {
	return nil
}

func (p cleanParams) Run() {
	fmt.Println("Going to clean with the following parameters:")
	fmt.Println(p)
}

type goCmd struct {
	Build buildParams `usage:"compile packages and dependencies"`
	Clean cleanParams `usage:"remove object files"`
}

var goCmdMetadata = flag.Flag{
	Usage:   "Go is a tool for managing Go source code.",
	Arglist: "command [argument]",
}

func (*goCmd) Metadata() map[string]flag.Flag {
	return map[string]flag.Flag{
		"":      goCmdMetadata,
		"build": buildMetadata,
	}
}

func main() {
	var actualArgs goCmd

	flagSet := flag.NewFlagSet(flag.Flag{})
	flagSet.ParseStruct(&actualArgs, os.Args...)

	var commands = map[string]command{
		"build": actualArgs.Build,
	}

	for cmdName, cmdArgs := range commands {
		if cmdArgs.IsEnabled() {
			err := cmdArgs.Validate()
			if err == nil {
				cmdArgs.Run()
				os.Exit(0)
			} else {
				fmt.Fprintln(os.Stderr, err)
				fmt.Println("")
				build, _ := flagSet.FindSubset(cmdName)
				build.Help(false) // display usage information for this command only
				os.Exit(1)
			}
		}
	}

	// no command was recognized in args
	flagSet.Help(false) // display usage information, with list of supported commands
	os.Exit(1)
}
