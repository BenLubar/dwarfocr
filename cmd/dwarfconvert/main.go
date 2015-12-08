//go:generate go get github.com/tv42/becky
//go:generate becky curses_640x300.png
//go:generate becky curses_800x600.png

package main

import (
	"flag"
	"fmt"
	_ "image/png"
	"os"
	"strings"

	"github.com/BenLubar/commander"
	"github.com/BenLubar/dwarfocr"
	"github.com/BenLubar/dwarfocr/internal/convert"
	"github.com/BenLubar/dwarfocr/internal/utils"

	_ "golang.org/x/image/bmp"
)

func main() {
	commander.RegisterFlags(flag.CommandLine)
	from := flag.String("f", "curses_640x300.png", "the tileset to convert from")
	to := flag.String("t", "curses_800x600.png", "the tileset to convert to")
	flag.Parse()

	err := commander.Init()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Warning: some profiles failed:", err)
		// don't exit
	}
	defer commander.Close()

	tilesFrom, err := utils.ReadTilesetFromFile(*from)
	if *from == "curses_640x300.png" && os.IsNotExist(err) {
		tilesFrom, err = dwarfocr.ReadTileset(strings.NewReader(curses_640x300))
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Fatal error loading tileset:", err)
		os.Exit(2)
	}

	tilesTo, err := utils.ReadTilesetFromFile(*to)
	if *to == "curses_800x600.png" && os.IsNotExist(err) {
		tilesTo, err = dwarfocr.ReadTileset(strings.NewReader(curses_800x600))
	}
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
		err = convert.ConvertOCR(os.Stdout, img, tilesFrom, tilesTo)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}
}

func png(a asset) string { return a.Content }
