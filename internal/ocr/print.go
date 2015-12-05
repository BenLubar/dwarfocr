package ocr

import (
	"errors"
	"fmt"
	"image"
	"io"

	"github.com/BenLubar/dwarfocr"
)

var cp437 = []rune(" ☺☻♥♦♣♠•◘○◙♂♀♪♫☼►◄↕‼¶§▬↨↑↓→←∟↔▲▼ !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~⌂ÇüéâäàåçêëèïîìÄÅÉæÆôöòûùÿÖÜ¢£¥₧ƒáíóúñÑªº¿⌐¬½¼¡«»░▒▓│┤╡╢╖╕╣║╗╝╜╛┐└┴┬├─┼╞╟╚╔╩╦╠═╬╧╨╤╥╙╘╒╓╫╪┘┌█▄▌▐▀αßΓπΣσµτΦΘΩδ∞φε∩≡±≥≤⌠⌡÷≈°∙·√ⁿ²■ ")

var ErrDimensionsMismatch = errors.New("dwarfocr: image dimensions do not match tileset dimensions")

func PrintOCR(w io.Writer, img *image.RGBA, tiles *dwarfocr.Tileset) error {
	var oldBg, oldFg, oldBright = -1, -1, -1

	setColor := func(bg, fg, bright int) {
		if oldFg != fg {
			fmt.Fprint(w, "\x1b[3", fg, "m")
			oldFg = fg
		}

		if oldBg != bg {
			fmt.Fprint(w, "\x1b[4", bg, "m")
			oldBg = bg
		}

		if oldBright != bright {
			if bright == 0 {
				fmt.Fprint(w, "\x1b[22m")
			} else {
				fmt.Fprint(w, "\x1b[1m")
			}
			oldBright = bright
		}
	}
	defer setColor(9, 9, 0)

	tileSize := tiles[0][0][0][0].Rect.Size()
	width, height := img.Rect.Dx(), img.Rect.Dy()
	if width%tileSize.X != 0 || height%tileSize.Y != 0 {
		return ErrDimensionsMismatch
	}

	width /= tileSize.X
	height /= tileSize.Y

	tileRect := image.Rectangle{image.ZP, tileSize}.Add(img.Rect.Min)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			bg, fg, bright, ch, err := tiles.Match(img.SubImage(tileRect.Add(image.Pt(x*tileSize.X, y*tileSize.Y))).(*image.RGBA))
			if err != nil {
				return err
			}
			setColor(bg, fg, bright)
			fmt.Fprint(w, string(cp437[ch]))
		}
		fmt.Fprintln(w)
	}
	return nil
}
