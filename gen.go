package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var o = os.Stdout

func main() {
	fmt.Printf("//generated by gen.go from edt.go\n")
	b, e := os.Open("edt.go")
	t(e)
	defer b.Close()
	io.Copy(o, b)

	for _, f := range []string{"cm.js", "cm.css"} {
		b, e := ioutil.ReadFile(f)
		t(e)
		n := strings.Replace(f, ".", "_", 1)
		fmt.Printf("const %s = %q\n", n, string(b))
	}
}
func t(e error) {
	if e != nil {
		panic(e)
	}
}
