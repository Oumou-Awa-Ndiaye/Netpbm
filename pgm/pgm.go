package pgm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

	scanner.Scan()
	magicNumber := scanner.Text()

	scanner.Scan()
	header := scanner.Text()
	headerParts := strings.Fields(header)
	width, _ := strconv.Atoi(headerParts[0])
	height, _ := strconv.Atoi(headerParts[1])
	max, _ := strconv.Atoi(headerParts[2])

	data := make([][]uint8, height)
	for i := 0; i < height; i++ {
		scanner.Scan()
		rowData := scanner.Text()
		rowDataParts := strings.Fields(rowData)
		row := make([]uint8, width)
		for j := 0; j < width; j++ {
			value, _ := strconv.Atoi(rowDataParts[j])
			row[j] = uint8(value)
		}
		data[i] = row
	}

	return &PGM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
		max:         max,
	}, nil
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

	fmt.Fprintf(writer, "%s\n", pgm.magicNumber)

	fmt.Fprintf(writer, "%d %d %d\n", pgm.width, pgm.height, pgm.max)

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
	pgm.max = int(maxValue)
}

func (pgm *PGM) Rotate90CW() {
	newData := make([][]uint8, pgm.width)
	for i := 0; i < pgm.width; i++ {
		newData[i] = make([]uint8, pgm.height)
	}

	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width; x++ {
			newData[x][pgm.height-y-1] = pgm.data[y][x]
		}
	}

	pgm.data = newData
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

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	fmt.Fprintf(writer, "%s\n", pbm.magicNumber)

	fmt.Fprintf(writer, "%d %d\n", pbm.width, pbm.height)

	for _, row := range pbm.data {
		for _, value := range row {
			if value {
				fmt.Fprint(writer, "1 ")
			} else {
				fmt.Fprint(writer, "0 ")
			}
		}
		fmt.Fprintln(writer)
	}

	return writer.Flush()
}
