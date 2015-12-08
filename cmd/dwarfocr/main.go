//go:generate go get github.com/tv42/becky
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
	"github.com/BenLubar/dwarfocr/internal/ocr"
	"github.com/BenLubar/dwarfocr/internal/utils"

	_ "golang.org/x/image/bmp"
)

func main() {
	commander.RegisterFlags(flag.CommandLine)
	tileset := flag.String("t", "curses_800x600.png", "the tileset to use")

	flag.Parse()

	err := commander.Init()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Warning: some profiles failed:", err)
		// don't exit
	}
	defer commander.Close()

	tiles, err := utils.ReadTilesetFromFile(*tileset)
	if *tileset == "curses_800x600.png" && os.IsNotExist(err) {
		tiles, err = dwarfocr.ReadTileset(strings.NewReader(curses_800x600))
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Fatal error loading tileset:", err)
		os.Exit(2)
	}

	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "Warning: no filenames given")
	}
	exit := 0
	for _, name := range flag.Args() {
		img, err := utils.RGBAFromFile(name)
		if err == nil {
			err = ocr.PrintOCR(os.Stdout, img, tiles)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error for", name+":", err)
			exit = 3
		}
	}
	os.Exit(exit)
}

func png(a asset) string { return a.Content }
