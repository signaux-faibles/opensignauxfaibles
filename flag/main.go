package main

import (
	"fmt"
	"os"

	flag "github.com/cosiner/flag"
)

type Tar struct {
	GZ          bool     `names:"-z, --gz" usage:"gzip format"`
	BZ          bool     `names:"-j, --bz" usage:"bzip2 format"`
	XZ          bool     `names:"-J, --xz" usage:"xz format"`
	Create      bool     `names:"-c" usage:"create tar file"`
	Extract     bool     `names:"-x" usage:"extract tar file"`
	File        string   `names:"-f" usage:"output file for create or input file for extract"`
	Directory   string   `names:"-C" usage:"extract directory"`
	SourceFiles []string `args:"true"`
}

func (t *Tar) Metadata() map[string]flag.Flag {
	const (
		usage   = "tar is a tool for manipulate tape archives."
		version = `
			version: v1.0.0
			commit: 10adf10dc10
			date:   2017-01-01 10:00:01
		`
		desc = `
		tar creates and manipulates streaming archive files.  This implementation can extract
		from tar, pax, cpio, zip, jar, ar, and ISO 9660 cdrom images and can create tar, pax,
		cpio, ar, and shar archives.
		`
	)
	return map[string]flag.Flag{
		"": {
			Usage:   usage,
			Version: version,
			Desc:    desc,
		},
		"--gz": {
			Desc: "use gzip format",
		},
	}
}

func main() {
	var tar Tar

	flag.NewFlagSet(flag.Flag{}).ParseStruct(&tar, os.Args...)
	// flag.NewFlagSet(flag.Flag{}).ParseStruct(&tar, "tar", "-zcf", "a.tgz", "a.go", "b.go")
	fmt.Println(tar.GZ)
	fmt.Println(tar.Create)
	fmt.Println(tar.File)
	fmt.Println(tar.SourceFiles)

	// Output:
	// true
	// true
	// a.tgz
	// [a.go b.go]
}
