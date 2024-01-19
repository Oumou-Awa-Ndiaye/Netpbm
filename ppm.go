package Netpbm

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Pixel struct {
	R, G, B uint8
}

type PPM struct {
	data          [][]Pixel
	width, height int
	magicNumber   string
	max           uint8
}

type Point struct {
	X, Y int
}

func (ppm *PPM) Size() (int, int) {
	return ppm.width, ppm.height
}

func ReadPPM(filename string) (*PPM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	ppm := &PPM{}

	/* ce bloc nous permet de lire et analyser les informations d'en-tête*/
	if !scanner.Scan() {
		return nil, errors.New("EOF while reading magic number")
	}
	ppm.magicNumber = scanner.Text()

	if ppm.magicNumber != "P3" {
		return nil, errors.New("Unsupported PPM format. Only P3 is supported.")
	}

	if !scanner.Scan() {
		return nil, errors.New("EOF while reading width and height")
	}
	fmt.Sscanf(scanner.Text(), "%d %d", &ppm.width, &ppm.height)

	if !scanner.Scan() {
		return nil, errors.New("EOF while reading max value")
	}
	max, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return nil, fmt.Errorf("Error converting max value to integer: %v", err)
	}
	ppm.max = uint8(max)

	ppm.data = make([][]Pixel, ppm.height)
	for i := range ppm.data {
		ppm.data[i] = make([]Pixel, ppm.width)
		for j := 0; j < ppm.width; j++ {
			/*Exit if EOF is detected during pixel reading*/
			if !scanner.Scan() {
				if scanner.Err() == nil {
					return nil, fmt.Errorf("Unexpected EOF while reading pixel data (row: %d, column: %d)", i, j)
				}
				return nil, scanner.Err()
			}

			line := scanner.Text()
			values := strings.Fields(line)

			// Ensure that there are enough values to read
			if len(values) < 3 {
				return nil, fmt.Errorf("Insufficient values for a pixel (row: %d, column: %d)", i, j)
			}

			for k := 0; k < 3; k++ {
				r, err := strconv.Atoi(values[k])
				if err != nil {
					return nil, fmt.Errorf("Error converting pixel value to integer (row: %d, column: %d): %v", i, j, err)
				}
				ppm.data[i][j] = Pixel{uint8(r), uint8(r), uint8(r)} // Assuming grayscale, adjust if necessary
				j++
			}
		}
	}

	return ppm, nil
}

func (ppm *PPM) At(x, y int) Pixel {
	return ppm.data[y][x]
}

func (ppm *PPM) Set(x, y int, value Pixel) {
	ppm.data[y][x] = value
}

func (ppm *PPM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	fmt.Fprintf(writer, "%s\n", ppm.magicNumber)

	fmt.Fprintf(writer, "%d %d %d\n", ppm.width, ppm.height, ppm.max)

	for _, row := range ppm.data {
		for _, pixel := range row {
			fmt.Fprintf(writer, "%d %d %d ", pixel.R, pixel.G, pixel.B)
		}
		fmt.Fprintln(writer)
	}

	return writer.Flush()
}

func (ppm *PPM) Invert() {
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			ppm.data[y][x].R = uint8(ppm.max) - ppm.data[y][x].R
			ppm.data[y][x].G = uint8(ppm.max) - ppm.data[y][x].G
			ppm.data[y][x].B = uint8(ppm.max) - ppm.data[y][x].B
		}
	}
}

func (ppm *PPM) Flip() {
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width/2; x++ {
			ppm.data[y][x], ppm.data[y][ppm.width-x-1] = ppm.data[y][ppm.width-x-1], ppm.data[y][x]
		}
	}
}

func (ppm *PPM) Flop() {
	for y := 0; y < ppm.height/2; y++ {
		ppm.data[y], ppm.data[ppm.height-y-1] = ppm.data[ppm.height-y-1], ppm.data[y]
	}
}

func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber
}

func (ppm *PPM) SetMaxValue(maxValue uint8) {
	ppm.max = uint8(maxValue)

}

func (ppm *PPM) Rotate90CW() {
	newData := make([][]Pixel, ppm.width)
	for i := 0; i < ppm.width; i++ {
		newData[i] = make([]Pixel, ppm.height)
	}

	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			newData[x][ppm.height-y-1] = ppm.data[y][x]
		}
	}

	ppm.data = newData
	ppm.width, ppm.height = ppm.height, ppm.width
}

func (ppm *PPM) ToPGM() *PGM {
	pgmData := make([][]uint8, ppm.height)
	for y := 0; y < ppm.height; y++ {
		pgmData[y] = make([]uint8, ppm.width)
		for x := 0; x < ppm.width; x++ {
			pgmData[y][x] = uint8((int(ppm.data[y][x].R) + int(ppm.data[y][x].G) + int(ppm.data[y][x].B)) / 3)
		}
	}

	return &PGM{
		data:        pgmData,
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: "P2",
		max:         uint8(ppm.max),
	}
}

func (ppm *PPM) ToPBM() *PBM {

	pbmData := make([][]bool, ppm.height)

	for y := 0; y < ppm.height; y++ {

		pbmData[y] = make([]bool, ppm.width)

		for x := 0; x < ppm.width; x++ {
			// And Here I Calculate the grayscale value by averaging RGB values
			grayValue := (int(ppm.data[y][x].R) + int(ppm.data[y][x].G) + int(ppm.data[y][x].B)) / 3

			// Here I Set the corresponding PBM pixel value based on the grayscale value
			pbmData[y][x] = grayValue < int(ppm.max)/2
		}
	}

	// I Create a new PBM object with the converted data
	return &PBM{
		data:        pbmData,
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: "P1",
	}
}

func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {
	dx := p2.X - p1.X
	dy := p2.Y - p1.Y

	if dx == 0 && dy == 0 {
		ppm.Set(p1.X, p1.Y, color)
		return
	}

	steps := int(math.Max(math.Abs(float64(dx)), math.Abs(float64(dy)))) + 1
	xIncrement := float64(dx) / float64(steps)
	yIncrement := float64(dy) / float64(steps)

	x, y := float64(p1.X), float64(p1.Y)

	for i := 0; i <= steps; i++ {
		ppm.Set(int(x+0.5), int(y+0.5), color)
		x += xIncrement
		y += yIncrement
	}
}

func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
	p2 := Point{p1.X + width - 1, p1.Y}
	p3 := Point{p1.X, p1.Y + height - 1}
	p4 := Point{p1.X + width - 1, p1.Y + height - 1}

	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p4, color)
	ppm.DrawLine(p4, p3, color)
	ppm.DrawLine(p3, p1, color)
}
func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
	for y := p1.Y; y < p1.Y+height; y++ {
		for x := p1.X; x < p1.X+width; x++ {
			ppm.Set(x, y, color)
		}
	}
}

func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			if x*x+y*y <= radius*radius {
				ppm.Set(center.X+x, center.Y+y, color)
			}
		}
	}
}

func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			if x*x+y*y <= radius*radius {
				ppm.Set(center.X+x, center.Y+y, color)
			}
		}
	}
}

func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p1, color)
}

func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {
	var vertices [3]Point
	if p1.Y <= p2.Y && p1.Y <= p3.Y {
		vertices[0] = p1
		vertices[1] = p2
		vertices[2] = p3
	} else if p2.Y <= p1.Y && p2.Y <= p3.Y {
		vertices[0] = p2
		vertices[1] = p1
		vertices[2] = p3
	} else {
		vertices[0] = p3
		vertices[1] = p1
		vertices[2] = p2
	}
	slope1 := float64(vertices[2].X-vertices[0].X) / float64(vertices[2].Y-vertices[0].Y)
	slope2 := float64(vertices[2].X-vertices[1].X) / float64(vertices[2].Y-vertices[1].Y)

	x1 := float64(vertices[0].X)
	x2 := float64(vertices[1].X)

	for y := vertices[0].Y; y <= vertices[1].Y; y++ {
		ppm.DrawLine(Point{int(x1 + 0.5), y}, Point{int(x2 + 0.5), y}, color)
		x1 += slope1
		x2 += slope2
	}

	x2 = float64(vertices[1].X)

	for y := vertices[1].Y + 1; y <= vertices[2].Y; y++ {
		ppm.DrawLine(Point{int(x1 + 0.5), y}, Point{int(x2 + 0.5), y}, color)
		x1 += slope1
		x2 += slope2
	}
}

func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
	for i := 0; i < len(points)-1; i++ {
		ppm.DrawLine(points[i], points[i+1], color)
	}
	ppm.DrawLine(points[len(points)-1], points[0], color)
}
func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) error {
	if ppm == nil {
		return errors.New("PPM structure is nil")
	}

	minY := points[0].Y
	maxY := points[0].Y

	for _, point := range points {
		if point.Y < minY {
			minY = point.Y
		}
		if point.Y > maxY {
			maxY = point.Y
		}
	}

	if minY < 0 || maxY >= ppm.height {
		return errors.New("Invalid Y coordinates for DrawFilledPolygon")
	}

	xCoordinates := make([][]int, maxY-minY+1)

	for i := 0; i < len(points); i++ {
		p1 := points[i]
		p2 := points[(i+1)%len(points)]

		var start, end Point
		if p1.Y <= p2.Y {
			start, end = p1, p2
		} else {
			start, end = p2, p1
		}

		slope := float64(end.X-start.X) / float64(end.Y-start.Y)

		x := float64(start.X)

		for y := start.Y; y <= end.Y; y++ {
			index := y - minY
			xCoordinates[index] = append(xCoordinates[index], int(x+0.5))
			x += slope
		}
	}

	for i := 0; i < len(xCoordinates); i += 2 {
		for j := 0; j < len(xCoordinates[i])-1; j += 2 {
			ppm.DrawLine(Point{xCoordinates[i][j], i + minY},
				Point{xCoordinates[i][j+1], i + minY}, color)
		}

		if i+1 < len(xCoordinates) {
			for j := 0; j < len(xCoordinates[i+1])-1; j += 2 {
				ppm.DrawLine(Point{xCoordinates[i+1][j], i + minY + 1},
					Point{xCoordinates[i+1][j+1], i + minY + 1}, color)
			}
		}
	}

	return nil
}

/*func (ppm *PPM) DrawKochSnowflake(center Point, radius, depth int, color Pixel) {
	angles := []float64{
		0.0,
		(2.0 * math.Pi / 3.0),
		(-2.0 * math.Pi / 3.0),
	}

	points := make([]Point, 3)
	for i := 0; i < 3; i++ {
		x := int(float64(center.X) + float64(radius)*math.Cos(angles[i]))
		y := int(float64(center.Y) + float64(radius)*math.Sin(angles[i]))
		points[i] = Point{x, y}
	}

	var drawKochSegment func(p1, p2 Point, depth int)
	drawKochSegment = func(p1, p2 Point, depth int) {
		if depth == 0 {
			ppm.DrawLine(p1, p2, color)
		} else {
			dx := p2.X - p1.X
			dy := p2.Y - p1.Y

			p3 := Point{p1.X + dx/3, p1.Y + dy/3}
			p4 := Point{p1.X + dx*2/3, p1.Y + dy*2/3}

			h := int(float64(radius) * math.Sin(math.Pi/3))

			midpoint := Point{(p1.X + p2.X) / 2, (p1.Y + p2.Y) / 2}

			p5 := Point{
				midpoint.X + h*(dy)/radius,
				midpoint.Y - h*(dx)/radius,
			}

			drawKochSegment(p1, p3, depth-1)
			drawKochSegment(p3, p5, depth-1)
			drawKochSegment(p5, p4, depth-1)
			drawKochSegment(p4, p2, depth-1)
		}
	}

	drawKochSegment(points[0], points[1], depth)
	drawKochSegment(points[1], points[2], depth)
	drawKochSegment(points[2], points[0], depth)
}
*/
