package Netpbm

import (
	"bufio"
	"errors"
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
	pgm := PGM{}
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	/*Scanner for text-based information*/
	scanner := bufio.NewScanner(file)

	scanner.Scan()
	pgm.magicNumber = scanner.Text()
	if pgm.magicNumber != "P5" && pgm.magicNumber != "P2" {
		return nil, fmt.Errorf("invalid magic number: %s", pgm.magicNumber)
	}

	/*Donnons les dimensions ici*/
	scanner.Scan()
	sepa := strings.Fields(scanner.Text())
	if len(sepa) != 2 {
		return nil, errors.New("bad input format")
	}
	pgm.width, _ = strconv.Atoi(sepa[0])
	pgm.height, _ = strconv.Atoi(sepa[1])

	if pgm.width <= 0 || pgm.height <= 0 {
		return nil, fmt.Errorf("invalid size: %d x %d", pgm.width, pgm.height)
	}

	/*Get max value*/
	scanner.Scan()
	maxValue, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return nil, fmt.Errorf("error parsing max value: %v", err)
	}
	pgm.max = uint8(maxValue)

	/*Move to the beginning of binary data*/
	for scanner.Scan() {
		if scanner.Text() == "" || strings.HasPrefix(scanner.Text(), "#") {
			continue
		} else {
			break
		}
	}
	/*J'ai pas fait réussi  le p5 aussi c'est en commentaire */
	// P5 format (raw binary)
	if pgm.magicNumber == "P5" {
		/*  buffer := make([]byte, pgm.width*pgm.height)

			// Read the binary data directly into the buffer
		    _, err := file.Read(buffer)
		    if err != nil {
			    return nil, fmt.Errorf("error reading binary data: %v", err)
			 }

			// Populate pgm.data with the binary data
			pgm.data = make([][]uint8, pgm.height)
			for i := 0; i < pgm.height; i++ {
				pgm.data[i] = make([]uint8, pgm.width)
			for j := 0; j < pgm.width; j++ {
				 pgm.data[i][j] = buffer[i*pgm.width+j]
		   }
			}*/
	} else if pgm.magicNumber == "P2" {
		/*P2 format (ASCII)*/
		pgm.data = make([][]uint8, pgm.height)
		for i := 0; i < pgm.height; i++ {
			pgm.data[i] = make([]uint8, pgm.width)
			lineValues := strings.Fields(scanner.Text())
			if len(lineValues) != pgm.width {
				return nil, fmt.Errorf("bad row length: %d", len(lineValues))
			}
			for j := 0; j < pgm.width; j++ {
				value, err := strconv.Atoi(lineValues[j])
				if err != nil {
					return nil, fmt.Errorf("error reading pixel value: %v", err)
				}
				pgm.data[i][j] = uint8(value)
			}
			if !scanner.Scan() {
				break
			}
		}
	} else {
		return nil, fmt.Errorf("unsupported PGM format: %s", pgm.magicNumber)
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

	fmt.Fprintln(writer, pgm.magicNumber)

	/* Here we  Write width, height, and max*/
	fmt.Fprintf(writer, "%d %d\n", pgm.width, pgm.height)
	fmt.Fprintf(writer, "%d\n", pgm.max)
	/*It's pixel data*/
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
	scaleFactor := float64(maxValue) / float64(pgm.max)

	/*Here we Update max value*/
	pgm.max = maxValue

	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width; x++ {
			pgm.data[y][x] = uint8(float64(pgm.data[y][x]) * scaleFactor)
		}
	}
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
			pbmData[y][x] = pgm.data[y][x] < uint8(pgm.max)/2
		}
	}

	return &PBM{
		data:        pbmData,
		width:       pgm.width,
		height:      pgm.height,
		magicNumber: "P1",
	}
}
