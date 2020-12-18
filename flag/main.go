// Test with `$ go build -o "go" . && ./go --help`

package main

import (
	"errors"
	"fmt"
	"os"

	flag "github.com/cosiner/flag"
)

type BuildParams struct {
	Enable   bool
	Already  bool     `names:"-a" important:"1" desc:"force rebuilding of packages that are already up-to-date."`
	Race     bool     `important:"1" desc:"enable data race detection.\nSupported only on linux/amd64, freebsd/amd64, darwin/amd64 and windows/amd64."`
	Output   string   `names:"-o" arglist:"output" important:"1" desc:"only allowed when compiling a single package"`
	LdFlags  string   `names:"-ldflags" arglist:"'flag list'" desc:"arguments to pass on each go tool link invocation."`
	Packages []string `args:"true"`
}

var BuildMetadata = flag.Flag{
	Arglist: "[-o output] [-i] [build flags] [packages]",
	Desc: `
	Build compiles the packages named by the import paths,
	along with their dependencies, but it does not install the results.
	...
	The build flags are shared by the build, clean, get, install, list, run,
	and test commands:
		`,
}

func (p BuildParams) Validate() error {
	if len(p.Packages) == 0 {
		return errors.New("Error: you should at least specify one package")
	}
	return nil
}

func (p BuildParams) Run() {
	fmt.Println("Going to build with the following parameters:")
	fmt.Println(p)
}

type CleanParams struct {
	Enable bool
}

func (p CleanParams) Validate() error {
	return nil
}

func (p CleanParams) Run() {
	fmt.Println("Going to clean with the following parameters:")
	fmt.Println(p)
}

type GoCmd struct {
	Build BuildParams `usage:"compile packages and dependencies"`
	Clean CleanParams `usage:"remove object files"`
}

var GoCmdMetadata = flag.Flag{
	Usage:   "Go is a tool for managing Go source code.",
	Arglist: "command [argument]",
}

func (*GoCmd) Metadata() map[string]flag.Flag {
	return map[string]flag.Flag{
		"":      GoCmdMetadata,
		"build": BuildMetadata,
	}
}

func main() {
	var g GoCmd

	set := flag.NewFlagSet(flag.Flag{})
	set.ParseStruct(&g, os.Args...)

	if g.Build.Enable {
		err := g.Build.Validate()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			fmt.Println("")
			build, _ := set.FindSubset("build")
			build.Help(false) // display usage information for the "go build" command only
		} else {
			g.Build.Run()
		}
	} else if g.Clean.Enable {
		err := g.Clean.Validate()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			fmt.Println("")
			clean, _ := set.FindSubset("clean")
			clean.Help(false) // display usage information for the "go clean" command only
		} else {
			g.Clean.Run()
		}
	} else {
		set.Help(false) // display usage information, with list of supported commands
	}
}
