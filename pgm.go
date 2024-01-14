package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           int
}

func ReadPGM(filename string) (*PGM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var pgm PGM
	scanner.Scan()
	pgm.magicNumber = scanner.Text()

	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "%d %d", &pgm.width, &pgm.height)

	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "%d", &pgm.max)

	pgm.data = make([][]uint8, pgm.height)
	for i := 0; i < pgm.height && scanner.Scan(); i++ {
		line := strings.Fields(scanner.Text())
		pgm.data[i] = make([]uint8, pgm.width)
		for j := 0; j < pgm.width; j++ {
			fmt.Sscanf(line[j], "%d", &pgm.data[i][j])
		}
	}

	return &pgm, nil
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

	fmt.Fprintf(writer, "%s\n%d %d\n%d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max)

	for _, row := range pgm.data {
		for _, pixel := range row {
			fmt.Fprintf(writer, "%d ", pixel)
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
	pgm.max = int(maxValue)
}

func (pgm *PGM) Rotate90CW() {
	rotated := make([][]uint8, pgm.width)
	for i := 0; i < pgm.width; i++ {
		rotated[i] = make([]uint8, pgm.height)
		for j := 0; j < pgm.height; j++ {
			rotated[i][j] = pgm.data[pgm.height-j-1][i]
		}
	}
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
