package game

import (
	"image/color"
	"math"
	"os"

	"github.com/clarktrimble/stam/fonts"
	"github.com/clarktrimble/stam/ifc"
	"github.com/clarktrimble/stam/stam"
	"github.com/clarktrimble/stam/vektor"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/pkg/errors"
	"golang.org/x/image/font"
)

// Game implements ebiten callbacks for fluid simulation demo
type Game struct {
	keys      []ebiten.Key
	button    bool
	help      bool
	gs        float64
	font      font.Face
	gridImage *ebiten.Image
	gridSize  int
	gridMid   int
	scale     int
	fluid     stam.Fluid
}

// New creates a game, given:
//
//	gridSize - width and height of grid used to model fluid
//	scale    - width and height used to display of each cell
func New(gridSize, scale int) (game *Game, err error) {

	fnt, err := fonts.Font(fontSize, fontDpi)
	if err != nil {
		err = errors.Wrapf(err, "failed to load font")
		return
	}

	game = &Game{
		font:      fnt,
		help:      true,
		gridImage: ebiten.NewImage(gridSize, gridSize),
		gridSize:  gridSize,
		gridMid:   gridSize / 2,
		scale:     scale,
		fluid:     stam.NewFluid(gridSize, visc, diff, dt, factory),
	}

	return
}

// Todo: look at Update and Draw performance and learn!!

// Update updates internal state with each tick
func (g *Game) Update() (err error) {

	g.checkKeys()
	g.checkButton()

	if g.button {
		g.fluid.AddDensity(g.gridMid, g.gridMid, densitySize, densityAdd)

		u, v := g.velocityFromCursor()
		g.fluid.AddVelocity(g.gridMid, g.gridMid, velocitySize, u, v)
	}

	g.fluid.Step()

	// remove minimum density and recalculate grayscale
	// all in the name of better contrast on-screen

	g.fluid.Level(g.fluid.Min())
	g.gs = float64(math.MaxUint8) / g.fluid.Max()

	return
}

// Draw draws on the screen with each tick
func (g *Game) Draw(screen *ebiten.Image) {

	if g.help {
		// Todo: relate font to cell size
		msg := "click to add dye and velocity\nesc to exit"
		text.Draw(screen, msg, g.font, 80, 100, color.Gray{Y: 196})

		return
	}

	var clr color.Color
	for i := 1; i <= g.gridSize; i++ {
		for j := 1; j <= g.gridSize; j++ {

			val := g.fluid.Density(i, j) * g.gs
			clr = color.Gray{Y: uint8(val)}

			g.gridImage.Set(i-1, j-1, clr)
		}
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(g.scale), float64(g.scale))
	screen.DrawImage(g.gridImage, op)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {

	size := g.gridSize * g.scale
	return size, size
}

// unexported

const (
	visc float64 = 0.01
	diff float64 = 0.0001
	dt   float64 = 0.001

	densityAdd       float64 = 9
	densitySize      int     = 5
	velocityMultiply float64 = 1000
	velocitySize     int     = 1

	fontSize float64 = 24
	fontDpi  float64 = 72
)

func factory(size int) (gridder ifc.Gridder) {
	return vektor.New(size)
}

func (g *Game) checkKeys() {

	g.keys = inpututil.AppendJustPressedKeys(g.keys[:0])
	for _, key := range g.keys {
		switch key {
		case ebiten.KeyEscape:
			os.Exit(0)
		}
	}
}

func (g *Game) checkButton() {

	mbl := ebiten.MouseButtonLeft
	if inpututil.IsMouseButtonJustPressed(mbl) {
		g.button = true
		g.help = false
	}
	if inpututil.IsMouseButtonJustReleased(mbl) {
		g.button = false
	}
}

func (g *Game) velocityFromCursor() (u, v float64) {

	// since grid is in upper left corner,
	// dividing pixel coord by scale gives the grid coord
	// almost, add one to account for one based grid coords
	// Todo: skip calculating grid coords and subtract middle pixel coord

	cursorX, cursorY := ebiten.CursorPosition()
	gridX, gridY := cursorX/g.scale+1, cursorY/g.scale+1

	// subtract midpoint for a vector pointing to cursor

	u = float64(gridX-g.gridMid) * velocityMultiply
	v = float64(gridY-g.gridMid) * velocityMultiply

	return
}
