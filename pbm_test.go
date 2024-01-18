package netpbm

import (
	"os"
	"reflect"
	"testing"
)

const imageWidth = 5
const imageHeight = 3

var testData = []bool{
	true, false, true, false, true,
	false, true, false, true, false,
	true, false, true, false, true,
}

var imageDataP4 = []bool{
	true, false, true, false, true,
	false, true, false, true, false,
	false, false, false, true, false,
}

func TestReadPBM(t *testing.T) {
	/* I create a temporary test file*/
	testFilename := "./imagesTest/pbmImages/imageP1.pbm"
	file, err := os.Create(testFilename)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testFilename)

	/* I write P1 data to the file*/
	file.WriteString("P1\n")
	file.WriteString("5 3\n")
	file.WriteString("1 0 1 0 1\n")
	file.WriteString("0 1 0 1 0\n")
	file.WriteString("1 0 1 0 1\n")
	file.Close()

	pbm, err := ReadPBM(testFilename)
	if err != nil {
		t.Fatal(err)
	}

	if pbm.magicNumber != "P1" {
		t.Error("Wrong magic number")
	}

	if pbm.width != imageWidth || pbm.height != imageHeight {
		t.Error("Wrong size")
	}

	if !reflect.DeepEqual(pbm.data, [][]bool{{true, false, true, false, true}, {false, true, false, true, false}, {true, false, true, false, true}}) {
		t.Error("Wrong data for P1")
	}

	file, err = os.Create(testFilename)
	if err != nil {
		t.Fatal(err)
	}
	file.WriteString("P4\n")
	file.WriteString("5 3\n")
	file.WriteString("10101010\n")
	file.WriteString("01010101\n")
	file.WriteString("01010100\n")
	file.Close()

	/*Here I prefer to read before check the magic number, the image with P4 magic number*/
	pbm, err = ReadPBM(testFilename)
	if err != nil {
		t.Fatal(err)
	}

	if pbm.magicNumber != "P4" {
		t.Error("Wrong magic number")
	}

	if pbm.width != imageWidth || pbm.height != imageHeight {
		t.Error("Wrong size")
	}

	if !reflect.DeepEqual(pbm.data, [][]bool{{true, false, true, false, true}, {false, true, false, true, false}, {false, false, false, true, false}}) {
		t.Error("Wrong data for P4")
	}
}

func TestSize(t *testing.T) {
	pbm := &PBM{data: [][]bool{{true, false, true, false, true}, {false, true, false, true, false}, {true, false, true, false, true}}, width: imageWidth, height: imageHeight, magicNumber: "P1"}
	w, h := pbm.Size()
	if w != imageWidth || h != imageHeight {
		t.Error("Wrong size")
	}
}

func TestAt(t *testing.T) {
	pbm := &PBM{data: [][]bool{{true, false, true, false, true}, {false, true, false, true, false}, {true, false, true, false, true}}, width: imageWidth, height: imageHeight, magicNumber: "P1"}
	if pbm.At(0, 2) != true {
		t.Error("Wrong value")
	}
}

func TestSet(t *testing.T) {
	pbm := &PBM{data: [][]bool{{true, false, true, false, true}, {false, true, false, true, false}, {true, false, true, false, true}}, width: imageWidth, height: imageHeight, magicNumber: "P1"}
	pbm.Set(1, 1, true)
	if pbm.At(1, 1) != true {
		t.Error("Wrong value")
	}
}

func TestSave(t *testing.T) {
	testFilename := "./imagesTest/pbmImages/imageP1.pbm"
	defer os.Remove(testFilename)

	pbm := &PBM{data: [][]bool{{true, false, true, false, true}, {false, true, false, true, false}, {true, false, true, false, true}}, width: imageWidth, height: imageHeight, magicNumber: "P1"}

	err := pbm.Save(testFilename)
	if err != nil {
		t.Fatal(err)
	}

	savedPBM, err := ReadPBM(testFilename)
	if err != nil {
		t.Fatal(err)
	}

	if savedPBM.magicNumber != "P1" {
		t.Error("Wrong magic number")
	}

	if savedPBM.width != imageWidth || savedPBM.height != imageHeight {
		t.Error("Wrong size")
	}

	if !reflect.DeepEqual(savedPBM.data, [][]bool{{true, false, true, false, true}, {false, true, false, true, false}, {true, false, true, false, true}}) {
		t.Error("Wrong data")
	}
}

func TestInvert(t *testing.T) {
	pbm := &PBM{data: [][]bool{{true, false, true, false, true}, {false, true, false, true, false}, {true, false, true, false, true}}, width: imageWidth, height: imageHeight, magicNumber: "P1"}
	pbm.Invert()

	if !reflect.DeepEqual(pbm.data, [][]bool{{false, true, false, true, false}, {true, false, true, false, true}, {false, true, false, true, false}}) {
		t.Error("Wrong data after inverting")
	}
}

func TestFlip(t *testing.T) {
	pbm := &PBM{data: [][]bool{{true, false, true, false, true}, {false, true, false, true, false}, {true, false, true, false, true}}, width: imageWidth, height: imageHeight, magicNumber: "P1"}
	pbm.Flip()

	if !reflect.DeepEqual(pbm.data, [][]bool{{true, false, true, false, true}, {false, true, false, true, false}, {true, false, true, false, true}}) {
		t.Error("Wrong data after flipping")
	}
}

func TestFlop(t *testing.T) {
	pbm := &PBM{
		data:        [][]bool{{true, false, true, false, true}, {false, true, false, true, false}, {true, false, true, false, true}},
		width:       imageWidth,
		height:      imageHeight,
		magicNumber: "P1",
	}

	pbm.Flop()

	expectedData := [][]bool{{true, false, true, false, true}, {false, true, false, true, false}, {true, false, true, false, true}}

	if !reflect.DeepEqual(pbm.data, expectedData) {
		t.Error("Wrong data after flop")
	}
}

func TestSetMagicNumber(t *testing.T) {
	pbm := &PBM{data: [][]bool{{true, false, true, false, true}, {false, true, false, true, false}, {true, false, true, false, true}}, width: imageWidth, height: imageHeight, magicNumber: "P1"}
	pbm.SetMagicNumber("P4")
	if pbm.magicNumber != "P4" {
		t.Error("Wrong magic number after setting")
	}
}
