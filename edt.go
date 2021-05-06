package main

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	file, args := "", os.Args
	if len(args) > 2 {
		panic("args")
	} else if len(args) == 2 {
		dir, e := os.Getwd()
		fatal(e)
		file = filepath.Join(dir, args[1])
	}

	http.HandleFunc("/cm.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/javascript")
		io.Copy(w, strings.NewReader(cm_js))
	})
	http.HandleFunc("/cm.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/css")
		io.Copy(w, strings.NewReader(cm_css))
	})
	http.HandleFunc("/favicon.png", fav)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(mainpage)) })
	http.HandleFunc("/r", rd)
	http.HandleFunc("/w", wr)

	addr := "127.0.0.1:0"
	srv := &http.Server{Addr: addr}
	ln, e := net.Listen("tcp", addr)
	fmt.Printf("http://%s/%s\n", ln.Addr(), filepath.ToSlash(file))
	fatal(e)
	e = srv.Serve(ln)
	fatal(e)
}
func rd(w http.ResponseWriter, r *http.Request) {
	q := r.URL.RawQuery
	fmt.Println("r", q)
	f, e := os.Open(q)
	if e != nil {
		fmt.Fprintf(w, "%s", e)
		return
	}
	defer f.Close()
	if fi, e := f.Stat(); e != nil {
		fmt.Fprintf(w, "%s", e)
	} else if fi.IsDir() == false {
		io.Copy(w, f)
	} else {
		if names, e := f.Readdirnames(-1); e != nil {
			fmt.Fprintf(w, "%s\n", e)
		} else {
			for _, s := range names {
				fmt.Fprintf(w, "%s\n", s)
			}
		}
	}
}
func wr(w http.ResponseWriter, r *http.Request) {
	q := r.URL.RawQuery
	b, e := ioutil.ReadAll(r.Body)
	if e != nil {
		fmt.Fprintf(w, "%s", e)
	} else {
		fmt.Printf("w %s (%d)\n", q, len(b))
	}
	fi, e := os.Stat(q)
	if e != nil {
		fmt.Fprintf(w, "w %s: %s", q, e) // only write over existing files
		return
	} else if fi.IsDir() {
		fmt.Fprintf(w, "w %s (is directory)", q)
		return
	}
	if e := ioutil.WriteFile(q, b, fi.Mode()); e != nil {
		fmt.Fprintf(w, "w %s: %s", q, e)
	}
}
func fav(w http.ResponseWriter, r *http.Request) {
	m := image.NewRGBA(image.Rectangle{Max: image.Point{48, 48}})
	for i := 0; i < 48; i++ {
		for k := i; k < 48; k++ {
			m.Set(i, k, color.RGBA{168, 50, 111, 255})
		}
	}
	w.Header().Set("content-type", "image/png")
	png.Encode(w, m)
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}

//go:embed main.html
var mainpage string

//go:embed cm.js
var cm_js string

//go:embed cm.css
var cm_css string
