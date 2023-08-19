// Package fonts embeds a font and provides for loading it
// magics from https://github.com/hajimehoshi/ebiten/blob/main/examples/resources/fonts/embed.go
package fonts

import (
	_ "embed"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	//go:embed mplus-1p-regular.ttf
	MPlus1pRegular_ttf []byte
)

func Font(size, dpi float64) (fnt font.Face, err error) {

	tt, err := opentype.Parse(MPlus1pRegular_ttf)
	if err != nil {
		return
	}

	opt := &opentype.FaceOptions{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	}

	fnt, err = opentype.NewFace(tt, opt)
	return
}
