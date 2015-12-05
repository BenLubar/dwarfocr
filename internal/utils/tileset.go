package utils

import (
	"os"

	"github.com/BenLubar/dwarfocr"
)

func ReadTilesetFromFile(name string) (*dwarfocr.Tileset, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return dwarfocr.ReadTileset(f)
}
