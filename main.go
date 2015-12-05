// +build !translate

package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	tileset := flag.String("t", "curses_800x600.png", "the tileset to use")
	flag.Parse()

	tiles, err := ReadTilesetFromFile(*tileset)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Fatal error loading tileset:", err)
		os.Exit(2)
	}

	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "Warning: no filenames given")
	}
	exit := 0
	for _, name := range flag.Args() {
		img, err := RGBAFromFile(name)
		if err == nil {
			err = PrintOCR(os.Stdout, img, tiles)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error for", name+":", err)
			exit = 3
		}
	}
	os.Exit(exit)
}
