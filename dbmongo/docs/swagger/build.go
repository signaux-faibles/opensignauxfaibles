package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	file, _ := os.Open(os.Args[1])
	bytes, _ := ioutil.ReadAll(file)
	var byteGlyphs []string
	for _, b := range bytes {
		byteGlyphs = append(byteGlyphs, fmt.Sprintf("%d", b))
	}
	fmt.Println("package docs\n")
	fmt.Println("var jsonBytes = []byte{" + strings.Join(byteGlyphs, ", ") + "}")
	fmt.Println("var doc = string(jsonBytes)")

}
