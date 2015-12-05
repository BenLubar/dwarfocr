package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/BenLubar/dwarfocr/internal/utils"

	_ "golang.org/x/image/bmp"
)

func main() {
	from := flag.String("f", "curses_640x300.png", "the tileset to convert from")
	to := flag.String("t", "curses_800x600.png", "the tileset to convert to")
	flag.Parse()

	tilesFrom, err := utils.ReadTilesetFromFile(*from)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Fatal error loading tileset:", err)
		os.Exit(2)
	}

	tilesTo, err := utils.ReadTilesetFromFile(*to)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Fatal error loading tileset:", err)
		os.Exit(2)
	}

	if flag.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "need 0 filenames")
		os.Exit(2)
	}

	img, err := utils.ReadRGBA(os.Stdin)
	if err == nil {
		err = ConvertOCR(os.Stdout, img, tilesFrom, tilesTo)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}
}
