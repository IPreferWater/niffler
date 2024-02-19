package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	nifflerImg *ebiten.Image
)

const (
	screenWidth  = 600
	screenHeight = 360

	spriteNifflerWidth  = 120
	spriteNifflerHeight = 79
)

type Game struct {
	worldSpeed float64
	frameCount int
	niffler    Niffler
	rfidTicket RfidTicket
	keys       []ebiten.Key
}

type Niffler struct {
	x      int
	y      int
	radius int
	speed  int
}

type RfidTicket struct {
	x      int
	y      int
	width int
	height  int
}

var (
	emptyCircle *ebiten.Image
)

func (g *Game) Update() error {
	g.frameCount++
	//1 second tick
	if g.frameCount%60 == 0 {
		rssi := getRssiPower(float64(g.niffler.x), float64(g.niffler.y), float64(g.rfidTicket.x), float64(g.rfidTicket.y), float64(g.niffler.radius))
		publish(rssi)
		g.frameCount=0
	}



	keysPressed := inpututil.AppendPressedKeys(g.keys)
	if len(keysPressed) != 0 {
		if contains(keysPressed, ebiten.KeyUp) {
			g.niffler.y -= g.niffler.speed
		}

		if contains(keysPressed, ebiten.KeyDown) {
			g.niffler.y += g.niffler.speed
		}

		if contains(keysPressed, ebiten.KeyLeft) {
			g.niffler.x -= g.niffler.speed
		}

		if contains(keysPressed, ebiten.KeyRight) {
			g.niffler.x += g.niffler.speed
		}
	}

	//rssi := getRssiPower(float64(g.niffler.x), float64(g.niffler.y), float64(g.rfidTicket.x), float64(g.rfidTicket.y), float64(g.niffler.radius))
	return nil
}

//TODO it might be a good idea to remove the key that was detected ? we won't need it anymore
func contains(keys []ebiten.Key, key ebiten.Key) bool {
	for _, v := range keys {
		if v == key{
			return true
		}
	}
	return false
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(spriteNifflerWidth)/2, -float64(spriteNifflerHeight)/2)
	op.GeoM.Translate(float64(g.niffler.x), float64(g.niffler.y))

	screen.DrawImage(nifflerImg, op)
	// circle that represent the radius of the niffler
	vector.StrokeCircle(screen, float32(g.niffler.x), float32(g.niffler.y), 150, 6, color.White, false)


	ticket := g.rfidTicket
	vector.DrawFilledRect(screen,float32(ticket.x),float32(ticket.y),float32(ticket.width),float32(ticket.height),color.RGBA{
		R: 251,
		G: 192,
		B: 147,
		A: 1,
	}, false)

	rssi := getRssiPower(float64(g.niffler.x), float64(g.niffler.y), float64(g.rfidTicket.x), float64(g.rfidTicket.y), float64(g.niffler.radius))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("rssi %d",rssi))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func startGraphical() {
	img, err := getEbitenImageFromRes("./res/niffler.png")

	if err != nil {
		log.Fatal(err)
	}
	nifflerImg = img

	emptyCircle = ebiten.NewImage(200, 200)

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Animation (Ebitengine Demo)")
	if err := ebiten.RunGame(&Game{
		worldSpeed: 10,
		frameCount: 0,
		niffler: Niffler{
			x:      100,
			y:      100,
			radius: 50,
			speed:  4,
		},
		rfidTicket: RfidTicket{
			x:      400,
			y:      250,
			width:  50,
			height: 10,
		},
		keys: []ebiten.Key{},
	}); err != nil {
		log.Fatal(err)
	}
}

func getEbitenImageFromRes(path string) (*ebiten.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return ebiten.NewImageFromImage(img), nil
}

func getRssiPower(x1 float64, y1 float64, x2 float64, y2 float64, radius float64) int{
	distance := math.Sqrt(math.Pow(x2-x1,2)+math.Pow((y2-y1),2))
	if distance> radius*2.9 {
		return 0
	}
	return int(radius*2.9 - distance)
}
