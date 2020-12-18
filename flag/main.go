// Test with `$ go build -o "go" . && ./go --help`

package main

import (
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
type CleanParams struct {
	Enable bool
}
type GoCmd struct {
	Build BuildParams `usage:"compile packages and dependencies"`
	Clean CleanParams `usage:"remove object files"`
}

func (*GoCmd) Metadata() map[string]flag.Flag {
	return map[string]flag.Flag{
		"": {
			Usage:   "Go is a tool for managing Go source code.",
			Arglist: "command [argument]",
		},
		"build": {
			Arglist: "[-o output] [-i] [build flags] [packages]",
			Desc: `
		Build compiles the packages named by the import paths,
		along with their dependencies, but it does not install the results.
		...
		The build flags are shared by the build, clean, get, install, list, run,
		and test commands:
			`,
		},
	}
}

func main() {
	var g GoCmd

	set := flag.NewFlagSet(flag.Flag{})
	set.ParseStruct(&g, os.Args...)

	if g.Build.Enable {
		if len(g.Build.Packages) == 0 {
			fmt.Fprintln(os.Stderr, "Error: you should at least specify one package")
			fmt.Println("")
			build, _ := set.FindSubset("build")
			build.Help(false) // display usage information for the "go build" command only
		} else {
			fmt.Println("Going to build with the following parameters:")
			fmt.Println(g.Build)
		}
	} else if g.Clean.Enable {
		fmt.Println("Going to clean with the following parameters:")
		fmt.Println(g.Clean)
	} else {
		set.Help(false) // display usage information, with list of supported commands
	}
}
