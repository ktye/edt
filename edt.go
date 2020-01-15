package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	file, args, port := "", os.Args, 2020
	if len(args) > 2 {
		panic("args")
	} else if len(args) == 2 {
		dir, e := os.Getwd()
		fatal(e)
		file = filepath.Join(dir, args[1])
	}
	fmt.Printf("http://localhost:%d/%s\n", port, filepath.ToSlash(file))

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
	if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil {
		fmt.Println(err)
	}
}
func rd(w http.ResponseWriter, r *http.Request) {
	q := r.URL.RawQuery
	fmt.Println("r", q)
	f, e := os.Open(q)
	if e != nil {
		fmt.Fprintf(w, "%s\n", e)
	}
	defer f.Close()
	io.Copy(w, f)
}
func wr(w http.ResponseWriter, r *http.Request) {
	q := r.URL.RawQuery
	if b, err := ioutil.ReadAll(r.Body); err != nil {
		fmt.Fprintf(w, "%s\n", err)
	} else {
		fmt.Printf("w %s (%d)\n", q, len(b))
		// TODO: write r.Body to q
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

var mainpage = `<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>E</title>
<link rel="stylesheet" href="/cm.css" type="text/css" />
<link rel="icon" type="image/png" sizes="48x48" href="/favicon.png">
<script type="application/javascript" src="/cm.js"></script>
<style>*{ font-family:monospace; margin:0; padding:0;  }
html, body, #edt {position:absolute; left:0;right:0;top:0;bottom:0; }
#but{ position:absolute;right:0;bottom:0; }
</style>
</head>
<body>
<textarea id="edt"></textarea>
<button id="but" onclick="wr()">write</button>
<script>
ed = CodeMirror.fromTextArea(edt, {"lineNumbers":true})

function get(p, f) {
 console.log("get", p)
 var r = new XMLHttpRequest()
 r.onreadystatechange = function() { if (this.readyState == 4 && this.status == 200) { if (f) f(this.response, this); } }
 r.open("GET", p)
 r.send()
}
function post(p, f, b) {
 console.log("post", p)
 var r = new XMLHttpRequest()
 r.onreadystatechange = function() { if (this.readyState == 4 && this.status == 200) { if (f) f(this.response, this); } }
 r.open("POST", p)
 r.send(b)
}
function hash(s){window.location.hash=encodeURIComponent(s.trim())}
function rd(file) { 
 get('/r?'+file, function(s){console.log("content", s); ed.setValue(s); document.title=file})
}
function wr() {
 var file = window.location.pathname.substr(1)
 console.log("write!", file)
 post('/w?'+file, function(s){ console.log("post returns", s) }, ed.getValue())
}
rd(window.location.pathname.substr(1))

// button-3 search
document.addEventListener('contextmenu',function(e){e.preventDefault()})
ed.on('mousedown', function(cm, evt) {if(evt.button==2 && (ed.getSelection().length>0)){search();evt.preventDefault()}})
function search() {
 var t = ed.getSelection()
 var v = ed.getValue()
 var p = ed.getCursor()
 var c = ed.indexFromPos(ed.getCursor())
 c = (p.sticky == "after") ? c+t.length : c
 var n = v.indexOf(t, c)
 if (n == -1) { n = v.indexOf(t, 0) }
 ed.setSelection(ed.posFromIndex(n), ed.posFromIndex(n+t.length), {"scroll":true})
}

</script>
</body></html>
`
