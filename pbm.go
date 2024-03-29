package Netpbm

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

/* Here we have width, height, and pixel data.*/
func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	/*Là je crée un scanner pour lire dans le fichier puisque j'ai déjà fait du java*/
	scanner := bufio.NewScanner(file)

	var pbm PBM
	scanner.Scan()
	pbm.magicNumber = scanner.Text()

	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "%d %d", &pbm.width, &pbm.height)

	pbm.data = make([][]bool, pbm.height)
	for i := 0; i < pbm.height && scanner.Scan(); i++ {
		line := strings.Fields(scanner.Text())
		pbm.data[i] = make([]bool, pbm.width) /*Here it's initialize the slice for storing pixel values for the current line.*/
		for j := 0; j < pbm.width; j++ {
			if j < len(line) {
				value, err := strconv.Atoi(line[j])
				if err != nil {
					return nil, err
				}
				pbm.data[i][j] = value != 0
			} else {
				return nil, fmt.Errorf("Not enough data at line %d", i+1)
			}
		}
	}

	return &pbm, nil
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

	fmt.Fprintf(writer, "%s\n%d %d\n", pbm.magicNumber, pbm.width, pbm.height)

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
		pbm.data[y], pbm.data[pbm.height-y-1] = pbm.data[pbm.height-y-1], pbm.data[y]
	}
}

func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}
