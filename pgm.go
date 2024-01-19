package Netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PGM struct {
	magicNumber string
	width       int
	height      int
	max         uint8
	data        [][]uint8
}

func ReadPGM(filename string) (*PGM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	pgm := &PGM{}

	// Read magic number
	scanner.Scan()
	pgm.magicNumber = scanner.Text()

	// Read width, height, and max
	for i := 0; i < 3; i++ {
		scanner.Scan()
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			// Skip comments
			i--
			continue
		}

		switch i {
		case 0:
			fmt.Sscanf(line, "%d %d", &pgm.width, &pgm.height)
		case 2:
			maxValue, _ := strconv.Atoi(line)
			pgm.max = uint8(maxValue)
		}
	}

	// Read pixel data
	pgm.data = make([][]uint8, pgm.height)
	for i := range pgm.data {
		pgm.data[i] = make([]uint8, pgm.width)
		scanner.Scan()
		line := scanner.Text()
		values := strings.Fields(line)
		for j, v := range values {
			pixelValue, _ := strconv.Atoi(v)
			pgm.data[i][j] = uint8(pixelValue)
		}
	}

	return pgm, nil
}

func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[y][x]
}

func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[y][x] = value
}

func (pgm *PGM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write magic number
	fmt.Fprintln(writer, pgm.magicNumber)

	// Write width, height, and max
	fmt.Fprintf(writer, "%d %d\n", pgm.width, pgm.height)
	fmt.Fprintf(writer, "%d\n", pgm.max)

	// Write pixel data
	for _, row := range pgm.data {
		for _, value := range row {
			fmt.Fprintf(writer, "%d ", value)
		}
		fmt.Fprintln(writer)
	}

	return writer.Flush()
}

func (pgm *PGM) Invert() {
	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width; x++ {
			pgm.data[y][x] = uint8(pgm.max) - pgm.data[y][x]
		}
	}
}

func (pgm *PGM) Flip() {
	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width/2; x++ {
			pgm.data[y][x], pgm.data[y][pgm.width-x-1] = pgm.data[y][pgm.width-x-1], pgm.data[y][x]
		}
	}
}

func (pgm *PGM) Flop() {
	for y := 0; y < pgm.height/2; y++ {
		pgm.data[y], pgm.data[pgm.height-y-1] = pgm.data[pgm.height-y-1], pgm.data[y]
	}
}

func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

func (pgm *PGM) SetMaxValue(maxValue uint8) {
	pgm.max = uint8(maxValue)
}

func (pgm *PGM) Rotate90CW() {
	rotated := make([][]uint8, pgm.width)
	for i := 0; i < pgm.width; i++ {
		rotated[i] = make([]uint8, pgm.height)
		for j := 0; j < pgm.height; j++ {
			rotated[i][j] = pgm.data[pgm.height-j-1][i]
		}
	}

	pgm.data = rotated

	pgm.width, pgm.height = pgm.height, pgm.width
}

func (pgm *PGM) ToPBM() *PBM {
	pbmData := make([][]bool, pgm.height)
	for y := 0; y < pgm.height; y++ {
		pbmData[y] = make([]bool, pgm.width)
		for x := 0; x < pgm.width; x++ {
			pbmData[y][x] = pgm.data[y][x] > uint8(pgm.max/2)
		}
	}

	return &PBM{
		data:        pbmData,
		width:       pgm.width,
		height:      pgm.height,
		magicNumber: "P1",
	}
}
