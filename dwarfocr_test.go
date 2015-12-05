package dwarfocr

import (
	"bytes"
	_ "image/png"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestOCR(t *testing.T) {
	tiles, err := ReadTilesetFromFile("curses_640x300.png")
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer

	for _, id := range []string{"adv10", "adv11", "adv22", "adv33", "adv44", "adv6", "adv7", "adv8", "adv9", "arena1", "arena2", "dwf1", "dwf2", "dwf3", "dwf4", "dwf5", "dwf6", "dwf7", "dwf8", "dwf9", "legends1", "legends2"} {
		img, err := RGBAFromFile(filepath.Join("testdata", id+".png"))
		if err != nil {
			t.Error(err)
			continue
		}

		expected, err := ioutil.ReadFile(filepath.Join("testdata", id+".txt"))
		if err != nil {
			t.Error(err)
			continue
		}

		buf.Reset()
		err = PrintOCR(&buf, img, tiles)
		if err != nil {
			t.Error(err)
			continue
		}

		if !bytes.Equal(expected, buf.Bytes()) {
			t.Errorf("id: %q\nexpected:\n%sactual:\n%s", id, expected, buf.Bytes())
		}
	}
}
