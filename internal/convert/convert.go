package convert

import (
	"errors"
	"image"
	"image/draw"
	"image/png"
	"io"

	"github.com/BenLubar/dwarfocr"
)

var ErrDimensionsMismatch = errors.New("dwarfocr: image dimensions do not match tileset dimensions")

func ConvertOCR(w io.Writer, img *image.RGBA, from, to *dwarfocr.Tileset) error {
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
