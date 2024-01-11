package pbm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	magicNumber := scanner.Text()

	scanner.Scan()
	dimensions := strings.Split(scanner.Text(), " ")
	if len(dimensions) != 2 {
		return nil, fmt.Errorf("dimensions line should contain 2 elements, got: %v", dimensions)
	}

	width, err := strconv.Atoi(dimensions[0])
	if err != nil {
		return nil, fmt.Errorf("failed to convert width to int: %v", err)
	}
	height, err := strconv.Atoi(dimensions[1])
	if err != nil {
		return nil, fmt.Errorf("failed to convert height to int: %v", err)
	}

	data := make([][]bool, height)
	for i := range data {
		data[i] = make([]bool, width)
		scanner.Scan()
		line := scanner.Text()
		for j, c := range line {
			if j >= width {
				return nil, fmt.Errorf("index out of range in line %d, width: %d, line length: %d", i+2, width, len(line))
			}
			data[i][j] = (c == '1')
		}
	}

	return &PBM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
	}, nil
}

func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height
}

func (pbm *PBM) At(x, y int) bool {
	return pbm.data[y][x]
}

func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[y][x] = value
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
		for _, pixel := range row {
			if pixel {
				fmt.Fprint(writer, "1 ")
			} else {
				fmt.Fprint(writer, "0 ")
			}
		}
		fmt.Fprintln(writer)
	}

	return writer.Flush()
}

func (pbm *PBM) Invert() {
	for y := 0; y < pbm.height; y++ {
		for x := 0; x < pbm.width; x++ {
			pbm.data[y][x] = !pbm.data[y][x]
		}
	}
}

func (pbm *PBM) Flip() {
	for y := 0; y < pbm.height; y++ {
		for x := 0; x < pbm.width/2; x++ {
			pbm.data[y][x], pbm.data[y][pbm.width-x-1] = pbm.data[y][pbm.width-x-1], pbm.data[y][x]
		}
	}
}

func (pbm *PBM) Flop() {
	for y := 0; y < pbm.height/2; y++ {
		for x := 0; x < pbm.width; x++ {
			pbm.data[y][x], pbm.data[pbm.height-y-1][x] = pbm.data[pbm.height-y-1][x], pbm.data[y][x]
		}
	}
}

func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}
