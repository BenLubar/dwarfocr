package main

import (
	"net/http"

	"github.com/BenLubar/dwarfocr/cmd/dwarfconvert_docker/impl"

	_ "image/png"

	_ "golang.org/x/image/bmp"
)

func main() {
	http.HandleFunc("/", impl.Handle)

	http.ListenAndServe(":80", nil)
}
