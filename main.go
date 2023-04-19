/**
 * This is free and unencumbered software released into the public domain.
 *
 * Anyone is free to copy, modify, publish, use, compile, sell, or distribute this software, either in source code form or as a compiled binary, for any purpose, commercial or non-commercial, and by any means.
 *
 * In jurisdictions that recognize copyright laws, the author or authors of this software dedicate any and all copyright interest in the software to the public domain. We make this dedication for the benefit of the public at large and to the detriment of our heirs and successors. We intend this dedication to be an overt act of relinquishment in perpetuity of all present and future rights to this software under copyright law.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 *	For more information, please refer to <http://unlicense.org/>
 **/
package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"os"

	"github.com/Alexamakans/convert-png-to-tc/pkg/config"
)

var cfg config.Config

// writeOps is an indicator for file size/render time.
var writeOps = 0

// flushOps is an indicator for file size/render time.
var flushOps = 0

const PIXEL_PACKET_IDENTIFIER = 0x10
const FLUSH_IDENTIFIER = 0x11

func main() {
	imgFile, err := os.Open(inputFilePath)
	panicIfNotNil(err)

	img, err := png.Decode(imgFile)
	panicIfNotNil(err)

	panicIfNotNil(imgFile.Close())

	convertedSlice := convert(img)
	convertedByteSlice := convertSliceToByteSlice(convertedSlice)

	printAtOffset(convertedSlice, 0)
	printAtOffset(convertedSlice, 8)
	printAtOffset(convertedSlice, 16)
	printAtOffset(convertedSlice, 24)
	printAtOffset(convertedSlice, 32)
	printAtOffset(convertedSlice, 40)

	fmt.Printf("Image will require:\n - %d write ops\n - %d flush ops", writeOps, flushOps)
	panicIfNotNil(os.WriteFile(outputFilePath, convertedByteSlice, 0644))
}

func panicIfNotNil(err error) {
	if err != nil {
		panic(err)
	}
}

func convert(img image.Image) []int {
	var output []int

	for cellX := 0; cellX < widthInCells; cellX++ {
		for cellY := 0; cellY < heightInCells; cellY++ {
			cellPackets := processCell(img, cellX, cellY)
			output = append(output, cellPackets...)
		}
	}

	if !flushAfterEveryWrite && !flushAfterEveryCellWrite {
		flushOps++
		output = append(output, FLUSH_IDENTIFIER)
	}

	return output
}

func at(img image.Image, x, y int) (r int, g int, b int) {
	_r, _g, _b, a := img.At(int(x), int(y)).RGBA()
	// _r, _g, and _b are 0.0-1.0 and then multiplied by alpha.
	// Undoing that here to make them in the range 0-255
	return int((_r * 255) / a), int((_g * 255) / a), int((_b * 255) / a)
}

func convertSliceToByteSlice(s []int) []byte {
	var output []byte
	for _, val := range s {
		for i := 0; i < dataWidthInBytes; i++ {
			output = append(output, byte((val>>(i*8))&0xFF))
		}
	}
	return output
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

func processCell(img image.Image, cellX, cellY int) []int {
	var pixels []pixel
	for x := 0; x < cellWidth; x++ {
		for y := 0; y < cellHeight; y++ {
			_x, _y := cellX*cellWidth+x, cellY*cellHeight+y
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
		if c, ok := getExistingSimilarColor(p.c, pixelMap); ok {
			pixelMap[c] = append(pixelMap[c], p.x, p.y)
		} else {
			pixelMap[p.c] = []int{p.x, p.y}
		}
	}

	var output []int
	for c, positions := range pixelMap {
		r, g, b := c.r, c.g, c.b
		positionBytes := [cellWidth]int{}
		for i := 0; i < len(positions); i += 2 {
			x := positions[i]
			y := positions[i+1]
			inCellX := x % cellWidth
			inCellY := y % cellHeight
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

		if flushAfterEveryWrite {
			flushOps++
			output = append(output, FLUSH_IDENTIFIER)
		}
	}

	if flushAfterEveryCellWrite {
		flushOps++
		output = append(output, FLUSH_IDENTIFIER)
	}

	return output
}

func getExistingSimilarColor(c color, m map[color][]int) (color, bool) {
	for c2 := range m {
		if getColorDistance(c, c2) <= colorSimilarityThreshold {
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
