// +build translate

package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"os"
)

func main() {
	from := flag.String("f", "curses_640x300.png", "the tileset to convert from")
	to := flag.String("t", "curses_800x600.png", "the tileset to convert to")
	flag.Parse()

	tilesFrom, err := ReadTilesetFromFile(*from)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Fatal error loading tileset:", err)
		os.Exit(2)
	}

	tilesTo, err := ReadTilesetFromFile(*to)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Fatal error loading tileset:", err)
		os.Exit(2)
	}

	if flag.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "need 0 filenames")
		os.Exit(2)
	}

	img, err := RGBAFromFile("/dev/stdin")
	if err == nil {
		err = ConvertOCR(os.Stdout, img, tilesFrom, tilesTo)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}
}

func ConvertOCR(w io.Writer, img *image.RGBA, from, to *Tileset) error {
	sizeFrom := from[0][0][0][0].Rect.Size()
	sizeTo := to[0][0][0][0].Rect.Size()
	width, height := img.Rect.Dx(), img.Rect.Dy()
	if width%sizeFrom.X != 0 || height%sizeFrom.Y != 0 {
		return ErrDimensionsMismatch
	}

	width /= sizeFrom.X
	height /= sizeFrom.Y

	out := image.NewRGBA(image.Rect(0, 0, width*sizeTo.X, height*sizeTo.Y))

	tileRect := image.Rectangle{image.ZP, sizeFrom}.Add(img.Rect.Min)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			bg, fg, bright, ch, err := from.Match(img.SubImage(tileRect.Add(image.Pt(x*sizeFrom.X, y*sizeFrom.Y))).(*image.RGBA))
			if err != nil {
				return err
			}
			draw.Draw(out, image.Rectangle{image.ZP, sizeTo}.Add(image.Pt(sizeTo.X*x, sizeTo.Y*y)), to[bg][fg][bright][ch], to[bg][fg][bright][ch].Rect.Min, draw.Src)
		}
	}
	return png.Encode(w, out)
}
