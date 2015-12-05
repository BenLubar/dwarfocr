package dwarfocr

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"io"
)

type Tileset [8][8][2][256]*image.RGBA

var ErrInvalidDimensions = errors.New("dwarfocr: invalid tileset dimensions")

var colors = [8][2]color.RGBA{
	{
		{0x00, 0x00, 0x00, 0xff},
		{0x80, 0x80, 0x80, 0xff},
	},
	{
		{0x80, 0x00, 0x00, 0xff},
		{0xff, 0x00, 0x00, 0xff},
	},
	{
		{0x00, 0x80, 0x00, 0xff},
		{0x00, 0xff, 0x00, 0xff},
	},
	{
		{0x80, 0x80, 0x00, 0xff},
		{0xff, 0xff, 0x00, 0xff},
	},
	{
		{0x00, 0x00, 0x80, 0xff},
		{0x00, 0x00, 0xff, 0xff},
	},
	{
		{0x80, 0x00, 0x80, 0xff},
		{0xff, 0x00, 0xff, 0xff},
	},
	{
		{0x00, 0x80, 0x80, 0xff},
		{0x00, 0xff, 0xff, 0xff},
	},
	{
		{0xc0, 0xc0, 0xc0, 0xff},
		{0xff, 0xff, 0xff, 0xff},
	},
}

func ReadTileset(r io.Reader) (*Tileset, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	if bounds.Empty() || bounds.Dx()%16 != 0 || bounds.Dy()%16 != 0 {
		return nil, ErrInvalidDimensions
	}

	var t Tileset

	size := image.Pt(bounds.Dx()/16, bounds.Dy()/16)
	tileRect := image.Rectangle{image.ZP, size}.Add(bounds.Min)

	mask := image.NewAlpha16(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			if r, g, b, _ := img.At(x, y).RGBA(); r == 0xffff && g == 0 && b == 0xffff {
				mask.SetAlpha16(x, y, color.Transparent)
			} else {
				mask.SetAlpha16(x, y, color.Opaque)
			}
		}
	}

	base := image.NewRGBA(bounds)
	draw.DrawMask(base, bounds, img, image.ZP, mask, image.ZP, draw.Src)

	for bg := range t {
		bgc := colors[bg][0]
		for fg := range t[bg] {
			for b := range t[bg][fg] {
				fgc := colors[fg][b]
				p := image.NewRGBA(bounds)
				draw.Draw(p, bounds, image.NewUniform(bgc), image.ZP, draw.Src)
				draw.Draw(p, bounds, &multiply{base: base, color: fgc}, image.ZP, draw.Over)
				for y := 0; y < 16; y++ {
					for x := 0; x < 16; x++ {
						t[bg][fg][b][x+y*16] = p.SubImage(tileRect.Add(image.Pt(size.X*x, size.Y*y))).(*image.RGBA)
					}
				}
			}
		}
	}

	return &t, nil
}

type multiply struct {
	base  *image.RGBA
	color color.RGBA
}

func (m *multiply) At(x, y int) color.Color {
	return m.RGBAAt(x, y)
}

func (m *multiply) RGBAAt(x, y int) color.RGBA {
	rgba := m.base.RGBAAt(x, y)
	return color.RGBA{
		R: uint8(uint16(rgba.R) * uint16(m.color.R) / 0xff),
		G: uint8(uint16(rgba.G) * uint16(m.color.G) / 0xff),
		B: uint8(uint16(rgba.B) * uint16(m.color.B) / 0xff),
		A: rgba.A,
	}
}

func (m *multiply) Bounds() image.Rectangle {
	return m.base.Rect
}

func (m *multiply) ColorModel() color.Model {
	return color.RGBAModel
}

var ErrNoMatch = errors.New("dwarfocr: could not match tile")

func (tiles *Tileset) Match(tile *image.RGBA) (bg, fg, bright, ch int, err error) {
	for bg = range tiles {
		for fg = range tiles[bg] {
			for bright = range tiles[bg][fg] {
				for ch = range tiles[bg][fg][bright] {
					if tiles.equal(tiles[bg][fg][bright][ch], tile) {
						return
					}
				}
			}
		}
	}

	err = ErrNoMatch
	return
}

func (tiles *Tileset) equal(a, b *image.RGBA) bool {
	dx := a.Rect.Dx() * 4
	dy := a.Rect.Dy()
	ap := a.Pix
	bp := b.Pix
	as := a.Stride
	bs := b.Stride

	for y := 0; y < dy; y++ {
		if y != 0 {
			// at the beginning of the loop instead of the end
			// to prevent index out of bounds errors on the last
			// iteration.
			ap, bp = ap[as:], bp[bs:]
		}
		if !bytes.Equal(ap[:dx], bp[:dx]) {
			return false
		}
	}
	return true
}
