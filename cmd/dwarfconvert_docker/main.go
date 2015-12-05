package main

import (
	"bytes"
	"io"
	"net/http"
	"time"

	_ "image/png"

	_ "golang.org/x/image/bmp"

	"github.com/BenLubar/dwarfocr/internal/convert"
	"github.com/BenLubar/dwarfocr/internal/utils"
)

func must(tiles *Tileset, err error) *Tileset {
	if err != nil {
		panic(err)
	}
	return tiles
}

var curses640x300 = must(utils.ReadTilesetFromFile("curses_640x300.png"))
var curses800x600 = must(utils.ReadTilesetFromFile("curses_800x600.png"))

func main() {
	http.HandleFunc("/", handle)

	http.ListenAndServe(":80", nil)
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method == "POST" {
		f, header, err := r.FormFile("input")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer f.Close()

		img, err := utils.ReadRGBA(f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var buf bytes.Buffer
		err = convert.ConvertOCR(&buf, img, curses640x300, curses800x600)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.ServeContent(w, r, header.Filename, time.Now(), bytes.NewReader(buf.Bytes()))
		return
	}

	io.WriteString(w, `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>Dwarf Fortress Screenshot Converter</title>
</head>
<body>
<form action="" method="post">
<input name="input" id="input" type="file" accept=".png,.bmp">
<br>
<input type="submit" value="!!SCIENCE!!">
</form>
<script>
  (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
  (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
  m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
  })(window,document,'script','//www.google-analytics.com/analytics.js','ga');

  ga('create', 'UA-41367436-1', 'auto');
  ga('send', 'pageview');

</script>
</body>
</html>`)
}
