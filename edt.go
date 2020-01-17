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

var mainpage = `<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>E</title>
<link rel="stylesheet" href="/cm.css" type="text/css" />
<link rel="icon" type="image/png" sizes="48x48" href="/favicon.png">
<script type="application/javascript" src="/cm.js"></script>
<style>
   *{ font-family:monospace; margin:0; padding:0; }
#hdr{ position:fixed;display:flex;top:0;left:0;right:0; }
#tag{ background:#FFFFEA; border:1px solid black;}
#edt{height:auto;top:0;buttom:0;position:fixed;}
</style>
</head>
<body>
<div id="hdr">
 <button id="but" onclick="wr()" onkey="exec">write</button>
 <input id="tag" style="flex-grow:1"></input>
</div>
<textarea id="edt"></textarea>
<script>
var tag = document.getElementById("tag")
var edt = document.getElementById("edt")
var ed = CodeMirror.fromTextArea(edt, {"lineNumbers":true})

function get(p, f) {
 var r = new XMLHttpRequest()
 r.onreadystatechange = function() { if (this.readyState == 4 && this.status == 200) { if (f) f(this.response, this); } }
 r.open("GET", p)
 r.send()
}
function post(p, f, b) {
 var r = new XMLHttpRequest()
 r.onreadystatechange = function() { if (this.readyState == 4 && this.status == 200) { if (f) f(this.response, this); } }
 r.open("POST", p)
 r.send(b)
}
function hash(s){window.location.hash=encodeURIComponent(s.trim())}
function rd(file) { 
 get('/r?'+file, function(s){ed.setValue(s); document.title=file})
}
function wr() {
 var file = window.location.pathname.substr(1)
 post('/w?'+file, function(s){ if(s.length>0){tag.value=s} }, ed.getValue())
}
rd(window.location.pathname.substr(1))

// search selected: middle-button(all), right(next) (bug: mouseup(chrome), ff ok)
function pd(e){e.preventDefault();e.stopPropagation()}
document.addEventListener('contextmenu',function(e){e.preventDefault()})
ed.on('mousedown', function(cm, e) {
 if     (e.button==2 && (ed.getSelection().length>0)){search(ed.getSelection(),false);pd(e)}
 else if(e.button==1 && (ed.getSelection().length>0)){search(ed.getSelection(),true );pd(e)}
})
function indexAll(a, s) { var r = [], i = -1; while ((i = a.indexOf(s, i+1)) != -1){ r.push(i); }; return r; }
function search(t, all){
 var v = ed.getValue()
 var p = ed.getCursor()
 if(all) {
  var n = indexAll(v,t)
  for(var i=0; i<n.length; i++) ed.addSelection(ed.posFromIndex(n[i]), ed.posFromIndex(n[i]+t.length), {"scroll":true})
 } else {
  var c = ed.indexFromPos(ed.getCursor())
  c = (p.sticky == "after") ? c+t.length : c
  var n = v.indexOf(t, c)
  if (n == -1) { n = v.indexOf(t, 0) }
  ed.setSelection(ed.posFromIndex(n), ed.posFromIndex(n+t.length), {"scroll":true})
 }
}

// tag-bar: return(search/goto line/re→replace), mark+button-click: middle(search all), right(search next)
function goto(line){ ed.setSelection({"line":line,"ch":0},{"line":line}) }
function tagKey(e){
 if(e.keyCode!=13) return
 var v = tag.value
 if(v.length==0) return
 var line = Number(v)
 if(isNaN(line)==false&&line>=0&&Math.floor(line)==line){goto(line-1);return}
 var i = v.indexOf("→")
 if(i==-1){search(v,false);return}
 var a = v.slice(0,i)
 var b = v.slice(i+1)
 var s = ed.getSelections()
 var re = RegExp(a, "gm")
 for(var i=0; i<s.length; i++) { s[i]=s[i].replace(re,b) }
 ed.replaceSelections(s)
}
function tagSelection(){ return tag.value.slice(tag.selectionStart,tag.selectionEnd) }
function tagMouse(evt){
 var s = tagSelection()
 if     (evt.button==2 && (s.length>0)){search(s,false);evt.preventDefault()}
 else if(evt.button==1 && (s.length>0)){search(s,true );evt.preventDefault()}
}
tag.onkeydown = tagKey
tag.onmousedown = tagMouse

</script>
</body></html>
`
