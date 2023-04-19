package converter

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"os"

	"github.com/Alexamakans/convert-png-to-tc/pkg/config"
)

const PIXEL_PACKET_IDENTIFIER = 0x10
const FLUSH_IDENTIFIER = 0x11

// writeOps is an indicator for file size/render time.
var writeOps = 0

// flushOps is an indicator for file size/render time.
var flushOps = 0

type Converter struct {
	Config config.Config
}

func New() Converter {
	return Converter{
		Config: config.DefaultConfig,
	}
}

func (c Converter) ConvertImage(image image.Image) []byte {
	convertedSlice := c.convert(image)
	convertedByteSlice := c.convertSliceToByteSlice(convertedSlice)
	return convertedByteSlice
}

func (c Converter) ConvertFile(inputFilePath string, outputFilePath string) error {
	imageFile, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer imageFile.Close()

	image, err := png.Decode(imageFile)
	if err != nil {
		return err
	}

	data := c.ConvertImage(image)
	if err := os.WriteFile(outputFilePath, data, 0644); err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("%d write ops, %d flush ops", writeOps, flushOps))

	return nil
}

func (c Converter) convert(img image.Image) []int {
	var output []int

	for cellX := 0; cellX < c.Config.GetWidthInCells(); cellX++ {
		for cellY := 0; cellY < c.Config.GetHeightInCells(); cellY++ {
			cellPackets := c.processCell(img, cellX, cellY)
			output = append(output, cellPackets...)
		}
	}

	if !c.Config.GetFlushAfterEveryWrite() && !c.Config.GetFlushAfterEveryCellWrite() {
		flushOps++
		output = append(output, FLUSH_IDENTIFIER)
	}

	return output
}

func (c Converter) convertSliceToByteSlice(s []int) []byte {
	var output []byte
	for _, val := range s {
		for i := 0; i < c.Config.GetDataWidthInBytesAsInt(); i++ {
			output = append(output, byte((val>>(i*8))&0xFF))
		}
	}
	return output
}

func at(img image.Image, x, y int) (r int, g int, b int) {
	_r, _g, _b, a := img.At(int(x), int(y)).RGBA()
	// _r, _g, and _b are 0.0-1.0 and then multiplied by alpha.
	// Undoing that here to make them in the range 0-255
	return int((_r * 255) / a), int((_g * 255) / a), int((_b * 255) / a)
}

func print8BytesAtOffset(s []byte, offset int) {
	fmt.Println(s[offset], s[offset+1], s[offset+2], s[offset+3], s[offset+4], s[offset+5], s[offset+6], s[offset+7])
}

func printAtOffset(s []int, offset int) {
	fmt.Println(s[offset], s[offset+1], s[offset+2], s[offset+3], s[offset+4], s[offset+5], s[offset+6], s[offset+7])
}

type color struct {
	r, g, b int
}

type pixel struct {
	c    color
	x, y int
}

func (c Converter) processCell(img image.Image, cellX, cellY int) []int {
	var pixels []pixel
	for x := 0; x < config.CellWidth; x++ {
		for y := 0; y < config.CellHeight; y++ {
			_x, _y := cellX*config.CellWidth+x, cellY*config.CellHeight+y
			r, g, b := at(img, _x, _y)
			pixels = append(pixels, pixel{
				c: color{r, g, b},
				x: _x,
				y: _y,
			})
		}
	}

	// consolidate
	pixelMap := map[color][]int{}
	for _, p := range pixels {
		if c, ok := c.getExistingSimilarColor(p.c, pixelMap); ok {
			pixelMap[c] = append(pixelMap[c], p.x, p.y)
		} else {
			pixelMap[p.c] = []int{p.x, p.y}
		}
	}

	var output []int
	for pixelColor, positions := range pixelMap {
		r, g, b := pixelColor.r, pixelColor.g, pixelColor.b
		positionBytes := [config.CellWidth]int{}
		for i := 0; i < len(positions); i += 2 {
			x := positions[i]
			y := positions[i+1]
			inCellX := x % config.CellWidth
			inCellY := y % config.CellHeight
			positionBytes[inCellX] |= 1 << inCellY
		}

		writeOps++
		output = append(output, PIXEL_PACKET_IDENTIFIER)
		row := cellY
		col := cellX
		output = append(output, 1<<col)
		output = append(output, 1<<row)
		output = append(output, r)
		output = append(output, g)
		output = append(output, b)
		var reversedPositionBytes []int
		for i := len(positionBytes) - 1; i >= 0; i-- {
			reversedPositionBytes = append(reversedPositionBytes, positionBytes[i])
		}
		for _, p := range reversedPositionBytes {
			output = append(output, p)
		}

		if c.Config.GetFlushAfterEveryWrite() {
			flushOps++
			output = append(output, FLUSH_IDENTIFIER)
		}
	}

	if c.Config.GetFlushAfterEveryCellWrite() {
		flushOps++
		output = append(output, FLUSH_IDENTIFIER)
	}

	return output
}

func (c Converter) getExistingSimilarColor(c1 color, m map[color][]int) (color, bool) {
	for c2 := range m {
		if getColorDistance(c1, c2) <= c.Config.GetColorSimilarityThreshold() {
			return c2, true
		}
	}

	return color{}, false
}

func getColorDistance(c1, c2 color) int {
	r := math.Abs(float64(c1.r) - float64(c2.r))
	g := math.Abs(float64(c1.g) - float64(c2.g))
	b := math.Abs(float64(c1.b) - float64(c2.b))
	return int(r + g + b)
}
